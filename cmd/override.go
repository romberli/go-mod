package cmd

import (
	"strings"

	"github.com/romberli/go-util/constant"
	"github.com/spf13/viper"

	"github.com/romberli/go-mod/config"
)

// OverrideConfigByCLI read configuration from command line interface, it will override the config file configuration
func OverrideConfigByCLI() error {
	// override log
	err := overrideLogByCLI()
	if err != nil {
		return err
	}
	// override mod
	err = overrideModByCLI()
	if err != nil {
		return err
	}

	return nil
}

// overrideLogByCLI overrides the log section by command line interface
func overrideLogByCLI() error {
	// log.level
	if logLevel != constant.DefaultRandomString {
		logLevel = strings.ToLower(logLevel)
		viper.Set(config.LogLevelKey, logLevel)
	}
	// log.format
	if logFormat != constant.DefaultRandomString {
		logLevel = strings.ToLower(logFormat)
		viper.Set(config.LogFormatKey, logFormat)
	}

	return nil
}

// overrideModByCLI overrides the mod section by command line interface
func overrideModByCLI() error {
	// mod.dir
	if modDir != constant.DefaultRandomString {
		viper.Set(config.ModDirKey, modDir)
	}
	// mod.name
	if modName != constant.DefaultRandomString {
		viper.Set(config.ModNameKey, modName)
	}
	// mod.version
	if modVersion != constant.DefaultRandomString {
		viper.Set(config.ModVersionKey, modVersion)
	}
	// mod.useCompileVersion
	if modUseCompileVersionStr != constant.DefaultRandomString {
		viper.Set(config.ModUseCompileVersionKey, modUseCompileVersionStr)
	}

	return nil
}
