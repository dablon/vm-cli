package ssh

import (
	"fmt"
	"os"

	"golang.org/x/crypto/ssh"
)

// Client representa una conexión SSH
type Client struct {
	Host     string
	User     string
	Password string
	Port     string
	client   *ssh.Client
}

// NewClient crea un nuevo cliente SSH
func NewClient(host, user, password, port string) *Client {
	return &Client{
		Host:     host,
		User:     user,
		Password: password,
		Port:     port,
	}
}

// Connect establece la conexión SSH
func (c *Client) Connect() error {
	config := &ssh.ClientConfig{
		User: c.User,
		Auth: []ssh.AuthMethod{
			ssh.Password(c.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	addr := fmt.Sprintf("%s:%s", c.Host, c.Port)
	client, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		return fmt.Errorf("error al conectar: %w", err)
	}

	c.client = client
	return nil
}

// Execute ejecuta un comando remoto
func (c *Client) Execute(command string) (string, error) {
	if c.client == nil {
		return "", fmt.Errorf("no hay conexión SSH activa")
	}

	session, err := c.client.NewSession()
	if err != nil {
		return "", fmt.Errorf("error al crear sesión: %w", err)
	}
	defer session.Close()

	output, err := session.CombinedOutput(command)
	if err != nil {
		return "", fmt.Errorf("error al ejecutar comando: %w", err)
	}

	return string(output), nil
}

// Close cierra la conexión SSH
func (c *Client) Close() error {
	if c.client != nil {
		return c.client.Close()
	}
	return nil
}

// CreateUser crea un nuevo usuario en la VM
func (c *Client) CreateUser(username, password string) (string, error) {
	// Crear usuario
	cmd := fmt.Sprintf("sudo useradd -m -s /bin/bash %s && echo '%s:%s' | sudo chpasswd", username, username, password)
	output, err := c.Execute(cmd)
	if err != nil {
		return "", fmt.Errorf("error al crear usuario: %w", err)
	}
	return output, nil
}

// UserExists verifica si un usuario existe
func (c *Client) UserExists(username string) (bool, error) {
	cmd := fmt.Sprintf("id %s", username)
	_, err := c.Execute(cmd)
	if err != nil {
		return false, nil
	}
	return true, nil
}

// EnsureSSHKey genera y configura clave SSH para el usuario
func (c *Client) EnsureSSHKey(username string) (string, error) {
	cmd := fmt.Sprintf(`
		sudo -u %s mkdir -p ~/.ssh && 
		sudo -u %s ssh-keygen -t ed25519 -f ~/.ssh/id_ed25519 -N "" -q &&
		cat ~/.ssh/id_ed25519.pub | sudo tee -a ~/.ssh/authorized_keys &&
		sudo chmod 700 ~/.ssh &&
		sudo chmod 600 ~/.ssh/authorized_keys &&
		sudo chown -R %s:%s ~/.ssh
	`, username, username, username, username)
	
	output, err := c.Execute(cmd)
	if err != nil {
		return "", fmt.Errorf("error al generar SSH key: %w", err)
	}
	return output, nil
}

// GetSSHKey obtiene la clave SSH pública del usuario
func (c *Client) GetSSHKey(username string) (string, error) {
	cmd := fmt.Sprintf("sudo cat /home/%s/.ssh/id_ed25519.pub", username)
	output, err := c.Execute(cmd)
	if err != nil {
		return "", fmt.Errorf("error al obtener SSH key: %w", err)
	}
	return output, nil
}

// SaveHostKey guarda la fingerprint del host para conexiones futuras
func SaveHostKey(host, key string) error {
	homeDir, _ := os.UserHomeDir()
	configPath := homeDir + "/.vm-cli/known_hosts"
	
	// Asegurar que existe el directorio
	os.MkdirAll(homeDir+"/.vm-cli", 0755)
	
	// Append al archivo
	f, err := os.OpenFile(configPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	
	_, err = f.WriteString(fmt.Sprintf("%s %s\n", host, key))
	return err
}
