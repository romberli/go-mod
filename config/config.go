/*
Copyright © 2020 Romber Li <romber2001@gmail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package config

import (
	"fmt"
	"strings"

	"github.com/romberli/go-util/constant"
	"github.com/romberli/log"
	"github.com/spf13/viper"
)

var (
	ValidLogLevels  = []string{"debug", "info", "warn", "warning", "error", "fatal"}
	ValidLogFormats = []string{"text", "json"}
)

// SetDefaultConfig set default configuration, it is the lowest priority
func SetDefaultConfig(baseDir string) {
	// log
	SetDefaultLog(baseDir)
	// mod
	SetDefaultMod()
}

// SetDefaultLog sets the default value of log
func SetDefaultLog(baseDir string) {
	viper.SetDefault(LogLevelKey, log.DefaultLogLevel)
	viper.SetDefault(LogFormatKey, log.DefaultLogFormat)
}

// SetDefaultMod sets the default value of mod
func SetDefaultMod() {
	viper.SetDefault(ModDirKey, DefaultModDir)
	viper.SetDefault(ModNameKey, DefaultModName)
	viper.SetDefault(ModVersionKey, DefaultModVersion)
	viper.SetDefault(ModUseCompileVersionKey, DefaultModUseCompileVersion)
}

// TrimSpaceOfArg trims spaces of given argument
func TrimSpaceOfArg(arg string) string {
	args := strings.SplitN(arg, constant.EqualString, 2)

	switch len(args) {
	case 1:
		return strings.TrimSpace(args[0])
	case 2:
		argName := strings.TrimSpace(args[0])
		argValue := strings.TrimSpace(args[1])
		return fmt.Sprintf("%s=%s", argName, argValue)
	default:
		return arg
	}
}
