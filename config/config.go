package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Ethereum EthereumConfig `yaml:"ethereum"`
}

type ServerConfig struct {
	Port string `yaml:"port"`
	Host string `yaml:"host"`
}

type EthereumConfig struct {
	RPCURL  string `yaml:"rpc_url"`
	Timeout string `yaml:"timeout"`
}

func Load() *Config {
	config, err := loadFromYAML("config.yaml")
	if err == nil {
		return config
	}

	return &Config{
		Server: ServerConfig{
			Port: config.Server.Port,
			Host: config.Server.Host,
		},
		Ethereum: EthereumConfig{
			RPCURL:  config.Ethereum.RPCURL,
			Timeout: config.Ethereum.Timeout,
		},
	}
}

func loadFromYAML(filename string) (*Config, error) {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return nil, fmt.Errorf("config file %s not found", filename)
	}

	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse YAML config: %w", err)
	}

	return &config, nil
}

func LoadFromPath(configPath string) (*Config, error) {
	if !filepath.IsAbs(configPath) {
		wd, err := os.Getwd()
		if err != nil {
			return nil, fmt.Errorf("failed to get working directory: %w", err)
		}
		configPath = filepath.Join(wd, configPath)
	}

	return loadFromYAML(configPath)
}
