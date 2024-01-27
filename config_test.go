package sync

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestConfig(t *testing.T) {
	assert := assert.New(t)

	f, err := os.Open("config.example.yaml")
	if err != nil {
		assert.Fail(err.Error())
		return
	}

	var cfg *Config
	if err := yaml.NewDecoder(f).Decode(&cfg); err != nil {
		assert.Fail(err.Error())
		return
	}

	assert.Equal("ecapis.qdm.cloud", cfg.QDM.BaseURL)
	assert.Equal("mongodb://localhost:27017", cfg.Persistence.Address)
	assert.Equal("qdm", cfg.Persistence.Database)
}
