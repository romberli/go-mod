package mod

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/romberli/go-util/constant"
	"github.com/romberli/go-util/linux"

	"github.com/romberli/go-mod/config"
)

const (
	getPackageRootPathCommand = "go env GOMODCACHE"
	goModGraphCommandTemplate = "cd %s && go mod graph"

	defaultSpaceNum  = 2
	outputPrefix     = "├ "
	outputPrefixLast = "└ "
)

var (
	defaultPackageRootPath string
)

type Controller struct {
	baseDir string

	RootNode *Node
	m        map[string]*Node
}

func NewController(baseDir string) *Controller {
	if baseDir == constant.EmptyString {
		baseDir = config.DefaultModDir
	}
	return &Controller{
		baseDir: baseDir,
		m:       make(map[string]*Node),
	}
}

func NewControllerWithDefault() *Controller {
	return NewController(config.DefaultModDir)
}

func (c *Controller) Init() error {
	if !filepath.IsAbs(c.baseDir) {
		absPath, err := filepath.Abs(c.baseDir)
		if err != nil {
			return err
		}
		c.baseDir = absPath
	}

	var err error
	defaultPackageRootPath, err = getPackageRootPath()
	if err != nil {
		return err
	}

	c.RootNode = NewNode(c.baseDir, constant.EmptyString)

	return c.RootNode.Resolve(c.m)
}

func (c *Controller) GetNodes(name, version string) []*Node {
	var result []*Node

	for _, node := range c.m {
		if name == node.Name && (version == constant.EmptyString || version == node.Version) {
			result = append(result, node)
		}
	}

	return result
}

func (c *Controller) GetParentChain(name, version string) [][]*Node {
	nodes := c.GetNodes(name, version)
	var result [][]*Node

	for _, node := range nodes {
		chain := node.GetParentChain()
		result = append(result, chain...)
	}

	return result
}

func (c *Controller) PrintParentChain(name, version string) error {
	err := c.Init()
	if err != nil {
		return err
	}
	nodesList := c.GetParentChain(name, version)
	c.PrintNodesList(nodesList)

	return nil
}

func (c *Controller) PrintNodesList(nodesList [][]*Node) {
	prefix := outputPrefixLast
	for _, nodes := range nodesList {
		for i, node := range nodes {
			if i == constant.ZeroInt {
				fmt.Println(node.RootPath)
				continue
			}
			fmt.Printf("%s%s\n", strings.Repeat(constant.SpaceString, defaultSpaceNum*i), prefix+node.String())
		}
	}
}

func (c *Controller) GetGoModGraph() (string, error) {
	cmd := fmt.Sprintf(goModGraphCommandTemplate, c.baseDir)

	return linux.ExecuteCommand(cmd)
}

func getPackageRootPath() (string, error) {
	output, err := linux.ExecuteCommand(getPackageRootPathCommand)
	if err != nil {
		return constant.EmptyString, err
	}

	return strings.TrimSpace(output), nil
}
