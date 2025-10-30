package mod

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/buger/jsonparser"
	"github.com/pingcap/errors"
	"github.com/romberli/go-util/common"
	"github.com/romberli/go-util/constant"
	"github.com/romberli/go-util/linux"
	"github.com/romberli/log"
	"golang.org/x/mod/module"
)

const (
	AtString = "@"

	requireJSON = "Require"
	replaceJSON = "Replace"
	newJSON     = "New"

	findModFilesCommandTemplate = "find %s -type f -name go.mod"
	// getPackagesCommand          = `go list -m -f '{{if and (not .Indirect) (not .Replace)}}{{.Path}}@{{.Version}}{{end}}' all | sed '1d'`
	// getPackagesCommand           = `go mod edit -json | grep -v '"Indirect": true' | grep -A1 '"Path":' | grep -v '"Indirect":' | grep -E '(Path|Version)' | awk -F'"' 'NR%2==1{path=$4} NR%2==0{print path"@"$4}'`
	goModEditCommand             = `go mod edit -json`
	noSuchFileOrDirectoryMessage = "No such file or directory"
	missingGoModFile             = "go.mod: no such file or directory"
	missingGoSumFile             = "missing go.sum entry for go.mod file"
	gg
	goModDownloadCommand = "go mod download"
)

type PackageInfo struct {
	Path     string
	Version  string
	Indirect bool
}

func (pi *PackageInfo) String() string {
	s := pi.Path
	if pi.Version != constant.EmptyString {
		s += AtString + pi.Version
	}

	return s
}

type ReplaceInfo struct {
	Old *PackageInfo
	New *PackageInfo
}

type Node struct {
	RootPath string
	FullName string
	Name     string
	Version  string
	Finished bool

	ParentNodes []*Node
	ChildNodes  []*Node
}

func NewNode(rootPath, fullName string) *Node {
	var (
		name    string
		version string
	)
	if fullName != constant.EmptyString {
		nameList := strings.Split(fullName, AtString)
		if len(nameList) >= constant.OneInt {
			name = nameList[constant.ZeroInt]
			if len(nameList) == constant.TwoInt {
				version = nameList[constant.OneInt]
			}
		}
	}

	return &Node{
		RootPath: rootPath,
		FullName: fullName,
		Name:     name,
		Version:  version,
	}
}

func (n *Node) String() string {
	return n.FullName
}

func (n *Node) AddParentNode(parentNode *Node) {
	n.ParentNodes = append(n.ParentNodes, parentNode)
}

func (n *Node) AddChildNode(childNode *Node) {
	n.ChildNodes = append(n.ChildNodes, childNode)
}

func (n *Node) GetParentChain() [][]*Node {
	var result [][]*Node

	n.getParentChain(&result, []*Node{})

	return result
}

func (n *Node) getParentChain(result *[][]*Node, current []*Node) {
	current = append(current, n)

	if len(n.ParentNodes) == constant.ZeroInt {
		var tmp []*Node
		for i := len(current) - constant.OneInt; i >= constant.ZeroInt; i-- {
			tmp = append(tmp, current[i])
		}
		*result = append(*result, tmp)
		return
	}

	for _, parentNode := range n.ParentNodes {
		parentNode.getParentChain(result, current)
	}
}

func (n *Node) Resolve(m map[string]*Node) error {
	queue := []*Node{n}

	for len(queue) > 0 {
		currentNode := queue[0]
		queue = queue[1:]

		if currentNode.FullName != constant.EmptyString {
			m[currentNode.FullName] = currentNode
		}
		if currentNode.Finished {
			continue
		}

		packages, err := currentNode.getChildPackages()
		if err != nil {
			return err
		}

		for _, pkg := range packages {
			rootPath := currentNode.RootPath
			if currentNode.FullName == constant.EmptyString {
				// root node
				rootPath = defaultPackageRootPath
			}
			childNode, ok := m[pkg]
			if !ok {
				childNode = NewNode(rootPath, pkg)
			}
			childNode.AddParentNode(currentNode)
			currentNode.AddChildNode(childNode)

			queue = append(queue, childNode)
		}

		currentNode.Finished = true
	}

	return nil
}

func (n *Node) getChildPackages() ([]string, error) {
	modDirs, err := n.getModDirs()
	if err != nil {
		return nil, err
	}

	var (
		packages        []string
		requirePackages []*PackageInfo
		replacePackages []*ReplaceInfo
	)

	for _, dir := range modDirs {
		output, err := linux.ExecuteCommand(goModEditCommand, linux.WorkDirOption(dir), linux.UseSHCOption())
		if err != nil {
			return nil, err
		}
		// require
		data, _, _, err := jsonparser.Get(common.StringToBytes(output), requireJSON)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(data, &requirePackages)
		if err != nil {
			return nil, errors.Trace(err)
		}
		// replace
		data, _, _, err = jsonparser.Get(common.StringToBytes(output), replaceJSON)
		if err != nil {
			return nil, errors.Trace(err)
		}
		err = json.Unmarshal(data, &replacePackages)
		if err != nil {
			return nil, errors.Trace(err)
		}
	Loop:
		for _, requirePackage := range requirePackages {
			for _, replacePackage := range replacePackages {
				if requirePackage.Path == replacePackage.Old.Path {
					if strings.HasPrefix(replacePackage.New.Path, constant.DotString) ||
						replacePackage.New.Version == constant.EmptyString {
						log.Warnf("replace new path is a file path, will ignore this. packageName: %s, oldPath: %s, newPath: %s, newVersion: %s",
							n.FullName, replacePackage.Old.Path, replacePackage.New.Path, replacePackage.New.Version)
						continue Loop
					}
					requirePackage.Path = replacePackage.New.Path
					requirePackage.Version = replacePackage.New.Version
					break
				}
			}

			if !requirePackage.Indirect {
				packages = append(packages, requirePackage.String())
			}
		}

		// packagesList := strings.Split(strings.TrimSpace(output), constant.CRLFString)
		// for _, pkg := range packageList {
		// 	if pkg != constant.EmptyString && !common.ElementInSlice(packages, pkg) {
		// 		if strings.Contains(pkg, missingGoModFile) {
		// 			log.Warnf("package can not find appropriate go.mod, will ignore it. packageName: %s", pkg)
		// 			continue
		// 		}
		// 		if strings.Contains(pkg, missingGoSumFile) {
		// 			log.Warnf("package is missing go.sum file, will ignore it. packageName: %s", pkg)
		// 			continue
		// 		}
		// 		if strings.Contains(pkg, goModDownloadCommand) {
		// 			log.Warnf("packag is not downloaded, will ignore it. packageName: %s", pkg)
		// 			continue
		// 		}
		// 		packages = append(packages, pkg)
		// 	}
		// }
	}

	return packages, nil
}

func (n *Node) getModDirs() ([]string, error) {
	var err error
	path := n.RootPath
	if n.FullName != constant.EmptyString {
		path, err = module.EscapePath(strings.TrimSuffix(n.Name, constant.SlashString))
		if err != nil {
			return nil, errors.Trace(err)
		}
		path = filepath.Join(n.RootPath, path) + AtString + n.Version
	}
	cmd := fmt.Sprintf(findModFilesCommandTemplate, path)
	output, err := linux.ExecuteCommand(cmd)
	if err != nil {
		if strings.Contains(output, noSuchFileOrDirectoryMessage) {
			log.Warnf("path does not exist, maybe because the package is only dependent by certain build conditions, will ignore it. path: %s", path)
			return nil, nil
		}
		return nil, err
	}

	output = strings.TrimSpace(output)
	if output == constant.EmptyString {
		return nil, nil
	}

	var modDirs []string

	outputList := strings.Split(output, constant.CRLFString)
	for _, line := range outputList {
		dir := filepath.Dir(line)
		modDirs = append(modDirs, dir)
	}

	return modDirs, nil
}

type NodeList []*Node

func (nl NodeList) Reverse() NodeList {
	var result NodeList

	for i := len(nl) - constant.OneInt; i >= constant.ZeroInt; i-- {
		result = append(nl, nl[i])
	}

	return result
}
