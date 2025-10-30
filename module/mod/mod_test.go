package mod

import (
	"testing"

	"github.com/romberli/log"
	"github.com/stretchr/testify/assert"
)

const (
	testModDir     = "/Users/romber/source_code/go/src/github.com/romberli/go-util"
	testModName    = "github.com/tikv/client-go/v2"
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

	log.SetDisableEscape(true)
	log.SetDisableDoubleQuotes(true)
	log.SetLevel(log.WarnLevel)

	err := testController.PrintParentChain(testModName, testModVersion, false)
	asst.Nil(err, "test PrintParentChain() failed")
	if err != nil {
		log.Errorf("test PrintParentChain() failed: %+v", err)
	}
}
