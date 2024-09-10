package mod

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/pingcap/errors"
	"github.com/romberli/go-util/constant"
	"github.com/romberli/go-util/linux"
	"github.com/romberli/log"

	"github.com/romberli/go-mod/config"
)

const (
	getPackageRootPathCommand = "go env GOMODCACHE"
	goModGraphCommand         = "go mod graph"
	goModListCommandTemplate  = "go list -m %s"

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

	logger *log.Logger
}

func NewController(baseDir string, logger *log.Logger) *Controller {
	if baseDir == constant.EmptyString {
		baseDir = config.DefaultModDir
	}
	return &Controller{
		baseDir: baseDir,
		m:       make(map[string]*Node),
		logger:  logger,
	}
}

func NewControllerWithDefault() *Controller {
	return NewController(config.DefaultModDir, nil)
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

	c.RootNode = NewNode(c.baseDir, constant.EmptyString, c.logger)

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

func (c *Controller) PrintParentChain(name, version string, modUseCompileVersion bool) error {
	err := c.Init()
	if err != nil {
		return err
	}

	v := version
	if modUseCompileVersion {
		v, err = c.GetCompileVersion(name)
		if err != nil {
			return err
		}
	}

	nodesList := c.GetParentChain(name, v)
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
	return linux.ExecuteCommand(goModGraphCommand, linux.WorkDirOption(c.baseDir))
}

func (c *Controller) GetCompileVersion(name string) (string, error) {
	command := fmt.Sprintf(goModListCommandTemplate, name)
	output, err := linux.ExecuteCommand(command, linux.WorkDirOption(c.baseDir))
	if err != nil {
		return output, err
	}

	outputList := strings.Split(strings.TrimSpace(output), constant.SpaceString)
	if len(outputList) < constant.TwoInt {
		return output, errors.Errorf("Controller.GetCompileVersion(): output format is not valid. output: %s", output)
	}

	return outputList[constant.OneInt], nil
}

func getPackageRootPath() (string, error) {
	output, err := linux.ExecuteCommand(getPackageRootPathCommand)
	if err != nil {
		return output, err
	}

	return strings.TrimSpace(output), nil
}
