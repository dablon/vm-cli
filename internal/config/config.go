package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	VMHost     string `json:"vm_host"`
	VMUser     string `json:"vm_user"`
	VMPort     string `json:"vm_port"`
	AgentName  string `json:"agent_name"`
	RemoteUser string `json:"remote_user"`
}

func (c *Config) Save(path string) error {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	err = json.Unmarshal(data, &cfg)
	return &cfg, err
}

func DefaultConfigPath() string {
	home, _ := os.UserHomeDir()
	return home + "/.vm-cli/config.json"
}
