package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

// Load reads the configuration from a YAML file.
func Load(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		// If the file doesn't exist, return an empty config (or nil) and no error
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	defer f.Close()

	var cfg Config
	if err := yaml.NewDecoder(f).Decode(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
