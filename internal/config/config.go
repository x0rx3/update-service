package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	DSN                  string `yaml:"dsn"`
	Delay                int    `yaml:"delay"`
	Cache                string `yaml:"cache"`
	LogLevel             string `yaml:"log_level"`
	WorkerLimit          int    `yaml:"worker_limit"`
	ServerAddress        string `yaml:"server_address"`
	UpdateServerUrl      string `yaml:"update_server_url"`
	UpdateServerLogin    string `yaml:"update_server_login"`
	UpdateServerPassword string `yaml:"update_server_password"`
}

func NewConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	cfg := &Config{}

	err = yaml.NewDecoder(file).Decode(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
