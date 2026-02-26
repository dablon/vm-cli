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

// ExecuteWithSudo ejecuta un comando que requiere sudo pidiendo TTY
func (c *Client) ExecuteWithSudo(command string) (string, error) {
	if c.client == nil {
		return "", fmt.Errorf("no hay conexión SSH activa")
	}

	session, err := c.client.NewSession()
	if err != nil {
		return "", fmt.Errorf("error al crear sesión: %w", err)
	}
	defer session.Close()

	// Request PTY para sudo interactivo
	modes := ssh.TerminalModes{
		ssh.ECHO:          0,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}
	
	if err := session.RequestPty("xterm", 80, 40, modes); err != nil {
		return "", fmt.Errorf("error al solicitar PTY: %w", err)
	}

	// Ejecutar comando con sudo
	cmd := fmt.Sprintf("sudo %s", command)
	output, err := session.CombinedOutput(cmd)
	if err != nil {
		return "", fmt.Errorf("error al ejecutar comando sudo: %w", err)
	}

	return string(output), nil
}

// RunWithSudo ejecuta un comando que requiere sudo usando -S y <<< (herestring)
func (c *Client) RunWithSudo(command string) (string, error) {
	sudoCmd := fmt.Sprintf("sudo -S -k %s <<< '%s'", command, c.Password)
	return c.Execute(sudoCmd)
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
			ssh.KeyboardInteractive(func(name, instruction string, questions []string, echos []bool) ([]string, error) {
				// Responder automáticamente con la contraseña
				answers := make([]string, len(questions))
				for i := range questions {
					answers[i] = c.Password
				}
				return answers, nil
			}),
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
	// Crear usuario con sudo usando PTY
	cmd := fmt.Sprintf("/usr/sbin/useradd -m -s /bin/bash %s", username)
	output, err := c.ExecuteWithSudo(cmd)
	if err != nil {
		return "", fmt.Errorf("error al crear usuario: %w", err)
	}
	
	// Establecer contraseña
	passCmd := fmt.Sprintf("echo '%s:%s' | /usr/sbin/chpasswd", username, password)
	_, err = c.ExecuteWithSudo(passCmd)
	if err != nil {
		fmt.Printf("⚠️  Warning: no se pudo establecer contraseña: %v\n", err)
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
