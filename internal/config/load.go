package config

import (
	"fmt"
	"os"
	"gopkg.in/yaml.v3"
)

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	var cfg Config
	if err:= yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse cfg: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("parse cfg: %w", err)
	}

	return &cfg, nil
}