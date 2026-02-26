package ssh

import (
	"encoding/json"
	"os"
)

// SSHConfig representa la configuración del CLI
type SSHConfig struct {
	VMHost     string `json:"vm_host"`
	VMUser     string `json:"vm_user"`
	VMPort     string `json:"vm_port"`
	AgentName  string `json:"agent_name"`
	RemoteUser string `json:"remote_user"`
}

// Save guarda la configuración en un archivo
func (c *SSHConfig) Save(path string) error {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

// Load carga la configuración desde un archivo
func Load(path string) (*SSHConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg SSHConfig
	err = json.Unmarshal(data, &cfg)
	return &cfg, err
}

// DefaultConfigPath retorna la ruta por defecto de configuración
func DefaultConfigPath() string {
	home, _ := os.UserHomeDir()
	return home + "/.vm-cli/config.json"
}

// EnsureConfigDir asegura que el directorio de configuración exista
func EnsureConfigDir() error {
	home, _ := os.UserHomeDir()
	return os.MkdirAll(home+"/.vm-cli", 0755)
}
