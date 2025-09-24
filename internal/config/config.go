package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	CurrentUserName string `json:"current_user_name"`
	DatabaseURL     string `json:"db_url"`
}

const configFileName = ".gatorconfig.json"

func getConfigFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, configFileName), nil
}

func Read() (*Config, error) {
	configPath, err := getConfigFilePath()
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func (c *Config) save() error {
	configPath, err := getConfigFilePath()
	if err != nil {
		return err
	}
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(configPath, data, 0644)
}

func (c *Config) SetUser(name string) error {
	c.CurrentUserName = name
	if err := c.save(); err != nil {
		return err
	}
	return nil
}
