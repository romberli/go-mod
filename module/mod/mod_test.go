package mod

import (
	"testing"

	"github.com/romberli/log"
	"github.com/stretchr/testify/assert"
)

const (
	testModDir     = "/Users/romber/source_code/go/src/github.com/romberli/go-util"
	testModName    = "gopkg.in/yaml.v2"
	testModVersion = ""
)

var (
	testController *Controller
)

func init() {
	testController = NewController(testModDir)
}

func TestModController_All(t *testing.T) {
	TestModController_PrintParentChain(t)
}

func TestModController_PrintParentChain(t *testing.T) {
	asst := assert.New(t)

	err := testController.PrintParentChain(testModName, testModVersion, true)
	asst.Nil(err, "test PrintParentChain() failed")
	log.SetDisableEscape(true)
	log.SetDisableDoubleQuotes(true)
	if err != nil {
		log.Errorf("test PrintParentChain() failed: %+v", err)
	}
}
