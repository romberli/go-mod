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
	"github.com/romberli/go-util/constant"
)

// global constant
const (
	DefaultCommandName = "go-mod"
	DefaultBaseDir     = constant.CurrentDir

	DefaultModParent  = "parent"
	DefaultModChild   = "child"
	DefaultModAll     = "all"
	DefaultModDir     = "./"
	DefaultModName    = constant.EmptyString
	DefaultModVersion = constant.EmptyString
)

// configuration constant
const (
	LogLevelKey  = "log.level"
	LogFormatKey = "log.format"

	ModDirKey     = "mod.dir"
	ModNameKey    = "mod.name"
	ModVersionKey = "mod.version"
)
