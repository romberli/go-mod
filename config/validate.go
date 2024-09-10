package config

import (
	"github.com/pingcap/errors"
	"github.com/romberli/go-multierror"
	"github.com/romberli/go-util/common"
	"github.com/romberli/go-util/constant"
	"github.com/spf13/cast"
	"github.com/spf13/viper"

	"github.com/romberli/go-mod/pkg/message"
)

// ValidateConfig validates if the configuration is valid
func ValidateConfig() (err error) {
	merr := &multierror.Error{}

	// validate log section
	err = ValidateLog()
	if err != nil {
		merr = multierror.Append(merr, err)
	}
	// validate mod section
	err = ValidateMod()
	if err != nil {
		merr = multierror.Append(merr, err)
	}

	return errors.Trace(merr.ErrorOrNil())
}

// ValidateLog validates if log section is valid.
func ValidateLog() error {
	merr := &multierror.Error{}

	// validate log.level
	logLevel, err := cast.ToStringE(viper.Get(LogLevelKey))
	if err != nil {
		merr = multierror.Append(merr, errors.Trace(err))
	}
	if !common.ElementInSlice(ValidLogLevels, logLevel) {
		merr = multierror.Append(merr, message.NewMessage(message.ErrNotValidLogLevel, logLevel))
	}

	// validate log.format
	logFormat, err := cast.ToStringE(viper.Get(LogFormatKey))
	if err != nil {
		merr = multierror.Append(merr, errors.Trace(err))
	}
	if !common.ElementInSlice(ValidLogFormats, logFormat) {
		merr = multierror.Append(merr, message.NewMessage(message.ErrNotValidLogFormat, logFormat))
	}

	return merr.ErrorOrNil()
}

// ValidateMod validates if mod section is valid.
func ValidateMod() error {
	merr := &multierror.Error{}

	// validate mod.dir
	modDir, err := cast.ToStringE(viper.Get(ModDirKey))
	if err != nil {
		merr = multierror.Append(merr, errors.Trace(err))
	}
	if modDir == constant.EmptyString {
		merr = multierror.Append(merr, errors.New("mod directory should not be empty"))
	}

	// validate mod.name
	_, err = cast.ToStringE(viper.Get(ModNameKey))
	if err != nil {
		merr = multierror.Append(merr, errors.Trace(err))
	}

	// validate mod.version
	_, err = cast.ToStringE(viper.Get(ModVersionKey))
	if err != nil {
		merr = multierror.Append(merr, errors.Trace(err))
	}

	return merr.ErrorOrNil()
}
