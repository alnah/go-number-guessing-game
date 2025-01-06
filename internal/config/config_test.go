package config_test

import (
	_ "embed"
	"testing"

	"github.com/go-number-guessing-game/internal/config"
	"github.com/stretchr/testify/assert"
)

const (
	goodDataTestYAML = "data_test.yaml"
	badDataTestJSON  = "data_test.json"
)

func TestUnitLoadConfig(t *testing.T) {
	t.Run("return config map", func(t *testing.T) {
		configMap := config.LoadConfig("yaml", goodDataTestYAML)

		assert.Len(t, configMap, 1)
		assert.Equal(t, "test", configMap["key"])
	})

	t.Run("panic when error", func(t *testing.T) {
		assert.Panics(t, func() {
			config.LoadConfig("json", badDataTestJSON)
		}, "want panic when loading invalid JSON data structure")
	})
}
