package config

import (
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// wrapper for .forge.yaml structure
type forgeConfig struct {
	Concat Config `yaml:"concat"`
}

// LoadDefaults attempts to load configuration from default files:
// 1. .forge.yaml (nested under 'concat:')
// 2. .concat.yaml (flat)
func LoadDefaults() (*Config, error) {
	// 1. Try .forge.yaml
	if _, err := os.Stat(".forge.yaml"); err == nil {
		f, err := os.Open(".forge.yaml")
		if err != nil {
			return nil, err
		}
		defer f.Close()

		var fc forgeConfig
		if err := yaml.NewDecoder(f).Decode(&fc); err != nil {
			return nil, err
		}
		return &fc.Concat, nil
	}

	// 2. Try .concat.yaml
	if _, err := os.Stat(".concat.yaml"); err == nil {
		f, err := os.Open(".concat.yaml")
		if err != nil {
			return nil, err
		}
		defer f.Close()

		var cfg Config
		if err := yaml.NewDecoder(f).Decode(&cfg); err != nil {
			return nil, err
		}
		return &cfg, nil
	}

	// No config found, return nil (use defaults)
	return nil, nil
}
