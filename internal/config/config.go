package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type ReadConfigError struct {
	Err error
}

const ReadConfigMessage = "Failed to read config: %w"

func (e *ReadConfigError) Error() string {
	return fmt.Errorf(ReadConfigMessage, e.Err).Error()
}

func NewReadConfigError(err error) error {
	return &ReadConfigError{Err: err}
}

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
