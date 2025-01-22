package config_test

import (
	_ "embed"
	"testing"

	"github.com/go-number-guessing-game/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestIntegrationLoadConfig(t *testing.T) {
	t.Run("return config map", func(t *testing.T) {
		configMap := config.LoadConfig("yaml", "../../configs/mock.yaml")

		assert.Len(t, configMap, 1)
		assert.Equal(t, "test", configMap["key"])
	})

	t.Run("panic when error", func(t *testing.T) {
		assert.Panics(t, func() {
			config.LoadConfig("yaml", "../../configs/bad.yaml")
		}, "want panic when loading invalid JSON data structure")
	})
}
