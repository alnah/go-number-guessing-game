package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// ReadConfigError represents an error that occurs when reading the configuration.
type ReadConfigError struct {
	Err error
}

// ReadConfigMessage is the error message format for configuration read failures.
const ReadConfigMessage = "Failed to read config: %w"

// Error returns the formatted error message for ReadConfigError.
func (e *ReadConfigError) Error() string {
	return fmt.Errorf(ReadConfigMessage, e.Err).Error()
}

// NewReadConfigError creates a new instance of ReadConfigError with the given error.
func NewReadConfigError(err error) error {
	return &ReadConfigError{Err: err}
}

// LoadConfig loads the configuration from the provided data based on the specified
// config type. It returns a map of configuration settings as key-value pairs.
func LoadConfig(configType string, filePath string) map[string]string {
	v := viper.New()
	v.SetConfigType(configType)
	v.SetConfigFile(filePath)

	// Attempt to read the configuration data. If it fails, panic with a custom error.
	if err := v.ReadInConfig(); err != nil {
		panic(NewReadConfigError(err))
	}

	// Retrieve all settings from the viper instance and convert them to a map.
	config := v.AllSettings()
	configMap := map[string]string{}
	for k, v := range config {
		configMap[k] = v.(string)
	}

	return configMap
}
