// Package config provides utilities for loading and managing app config
// settings using the Viper library. It includes error types for handling
// config read errors and a function to load settings from a file path and ext.
package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// ReadConfigError wraps an error that occurs during configuration file
// reading, providing better context.
type ReadConfigError struct {
	Err error
}

// ReadConfigMessage is displayed when a configuration read error occurs,
// formatted to include the underlying error.
const ReadConfigMessage = "Failed to read config: %w"

// Error returns the formatted error message for ReadConfigError.
func (e *ReadConfigError) Error() string {
	return fmt.Errorf(ReadConfigMessage, e.Err).Error()
}

// NewReadConfigError creates a new ReadConfigError instance with the
// specified underlying error.
func NewReadConfigError(err error) error {
	return &ReadConfigError{Err: err}
}

// LoadConfig reads configuration from the specified file path and type,
// returning a map of settings. It panics on read errors.
func LoadConfig(configType string, filePath string) map[string]string {
	v := viper.New()
	v.SetConfigType(configType)
	v.SetConfigFile(filePath)

	if err := v.ReadInConfig(); err != nil {
		panic(NewReadConfigError(err))
	}

	config := v.AllSettings()
	configMap := map[string]string{}
	for k, v := range config {
		configMap[k] = v.(string)
	}

	return configMap
}
