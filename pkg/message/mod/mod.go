package mod

import (
	"github.com/romberli/go-util/config"

	"github.com/romberli/go-mod/pkg/message"
)

func init() {
	initModDebugMessage()
	initModInfoMessage()
	initModErrorMessage()
}

const (
	// debug
	DebugModParentPrintParentChain = 100001

	// info
	InfoModParentPrintParentChain = 200001

	// error
	ErrModParentPrintParentChain = 400001
)

func initModDebugMessage() {

}

func initModInfoMessage() {

}

func initModErrorMessage() {
	message.Messages[ErrModParentPrintParentChain] = config.NewErrMessage(message.DefaultMessageHeader, ErrModParentPrintParentChain,
		"mod: print parent chain failed. mod directory: %s, mod name: %s, mod version: %s")
}
