package mod

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/pingcap/errors"
	"github.com/romberli/go-util/common"
	"github.com/romberli/go-util/constant"
	"github.com/romberli/go-util/linux"
	"github.com/romberli/log"
	"golang.org/x/mod/module"
)

const (
	AtString = "@"

	findModFilesCommandTemplate  = "find %s -type f -name go.mod"
	getPackagesCommand           = `go list -m -f '{{if not .Indirect}}{{.Path}}@{{.Version}}{{end}}' all | sed '1d'`
	noSuchFileOrDirectoryMessage = "No such file or directory"
	missingGoModFile             = "go.mod: no such file or directory"
	missingGoSumFile             = "missing go.sum entry for go.mod file"

	goModDownloadCommand = "go mod download"
)

type Node struct {
	RootPath string
	FullName string
	Name     string
	Version  string
	Finished bool

	ParentNodes []*Node
	ChildNodes  []*Node

	logger *log.Logger
}

func NewNode(rootPath, fullName string, logger *log.Logger) *Node {
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
		logger:   logger,
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
	if n.FullName != constant.EmptyString {
		m[n.FullName] = n
	}
	if n.Finished {
		return nil
	}
	packages, err := n.getChildPackages()
	if err != nil {
		return err
	}

	for _, pkg := range packages {
		rootPath := n.RootPath
		if n.FullName == constant.EmptyString {
			// root node
			rootPath = defaultPackageRootPath
		}
		childNode, ok := m[pkg]
		if !ok {
			childNode = NewNode(rootPath, pkg, n.logger)
		}
		childNode.AddParentNode(n)
		n.AddChildNode(childNode)

		err = childNode.Resolve(m)
		if err != nil {
			return err
		}
	}

	n.Finished = true

	return nil
}

func (n *Node) getChildPackages() ([]string, error) {
	modDirs, err := n.getModDirs()
	if err != nil {
		return nil, err
	}

	var packages []string

	for _, dir := range modDirs {
		output, err := linux.ExecuteCommand(getPackagesCommand, linux.WorkDirOption(dir), linux.UseSHCOption())
		if err != nil {
			return nil, err
		}

		packagesList := strings.Split(strings.TrimSpace(output), constant.CRLFString)
		for _, pkg := range packagesList {
			if pkg != constant.EmptyString && !common.ElementInSlice(packages, pkg) {
				if strings.Contains(pkg, missingGoModFile) {
					if n.logger != nil {
						n.logger.Warnf("package can not find appropriate go.mod, will ignore it. packageName: %s", pkg)
					}
					continue
				}
				if strings.Contains(pkg, missingGoSumFile) {
					if n.logger != nil {
						n.logger.Warnf("package is missing go.sum file, will ignore it. packageName: %s", pkg)
					}
					continue
				}
				if strings.Contains(pkg, goModDownloadCommand) {
					if n.logger != nil {
						n.logger.Warnf("packag is not downloaded, will ignore it. packageName: %s", pkg)
					}
					continue
				}
				packages = append(packages, pkg)
			}
		}
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
			if n.logger != nil {
				n.logger.Warnf("path does not exist, maybe because the package is only dependent by certain build conditions, will ignore it. path: %s", path)
			}
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
