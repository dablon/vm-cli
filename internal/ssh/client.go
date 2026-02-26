package ssh

import (
	"fmt"
	"strings"

	"golang.org/x/crypto/ssh"
)

// Config holds SSH connection configuration
type Config struct {
	Host     string
	User     string
	Password string
	Port     string
}

// Client represents an SSH connection to a remote VM
type Client struct {
	config *Config
	client *ssh.Client
}

// NewClient creates a new SSH client
func NewClient(host, user, password, port string) *Client {
	if port == "" {
		port = "22"
	}
	return &Client{
		config: &Config{
			Host:     host,
			User:     user,
			Password: password,
			Port:     port,
		},
	}
}

// Connect establishes an SSH connection
func (c *Client) Connect() error {
	if c.client != nil {
		return nil // Already connected
	}

	config := &ssh.ClientConfig{
		User: c.config.User,
		Auth: []ssh.AuthMethod{
			ssh.Password(c.config.Password),
			ssh.KeyboardInteractive(func(name, instruction string, questions []string, echos []bool) ([]string, error) {
				answers := make([]string, len(questions))
				for i := range questions {
					answers[i] = c.config.Password
				}
				return answers, nil
			}),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	addr := fmt.Sprintf("%s:%s", c.config.Host, c.config.Port)
	client, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}

	c.client = client
	return nil
}

// Close closes the SSH connection
func (c *Client) Close() error {
	if c.client != nil {
		err := c.client.Close()
		c.client = nil
		return err
	}
	return nil
}

// Execute runs a command on the remote VM
func (c *Client) Execute(command string) (string, error) {
	if c.client == nil {
		return "", fmt.Errorf("not connected")
	}

	session, err := c.client.NewSession()
	if err != nil {
		return "", fmt.Errorf("failed to create session: %w", err)
	}
	defer session.Close()

	output, err := session.CombinedOutput(command)
	if err != nil {
		return "", fmt.Errorf("command failed: %w", err)
	}

	return string(output), nil
}

// ExecuteWithSudo executes a command with sudo privileges
func (c *Client) ExecuteWithSudo(command string) (string, error) {
	if c.client == nil {
		return "", fmt.Errorf("not connected")
	}

	session, err := c.client.NewSession()
	if err != nil {
		return "", fmt.Errorf("failed to create session: %w", err)
	}
	defer session.Close()

	// Request PTY for interactive sudo
	modes := ssh.TerminalModes{
		ssh.ECHO:          0,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}

	if err := session.RequestPty("xterm", 80, 40, modes); err != nil {
		// Continue without PTY
		fmt.Printf("⚠️  PTY not available, continuing...\n")
	}

	sudoCmd := fmt.Sprintf("sudo %s", command)
	output, err := session.CombinedOutput(sudoCmd)
	if err != nil {
		return "", fmt.Errorf("sudo command failed: %w", err)
	}

	return string(output), nil
}

// CreateUser creates a new user on the remote VM
func (c *Client) CreateUser(username, password string) error {
	if c.client == nil {
		return fmt.Errorf("not connected")
	}

	// Validate username
	if !isValidUsername(username) {
		return fmt.Errorf("invalid username: %s", username)
	}

	// Create user
	cmd := fmt.Sprintf("/usr/sbin/useradd -m -s /bin/bash %s", username)
	_, err := c.ExecuteWithSudo(cmd)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	// Set password
	passCmd := fmt.Sprintf("echo '%s:%s' | /usr/sbin/chpasswd", username, password)
	_, err = c.ExecuteWithSudo(passCmd)
	if err != nil {
		return fmt.Errorf("failed to set password: %w", err)
	}

	return nil
}

// UserExists checks if a user exists on the remote VM
func (c *Client) UserExists(username string) (bool, error) {
	if c.client == nil {
		return false, fmt.Errorf("not connected")
	}

	cmd := fmt.Sprintf("id %s", username)
	_, err := c.Execute(cmd)
	if err != nil {
		return false, nil
	}
	return true, nil
}

// EnsureSSHKey generates and configures SSH key for user
func (c *Client) EnsureSSHKey(username string) (string, error) {
	if c.client == nil {
		return "", fmt.Errorf("not connected")
	}

	cmd := fmt.Sprintf(`
		sudo -u %s mkdir -p ~/.ssh && 
		sudo -u %s ssh-keygen -t ed25519 -f ~/.ssh/id_ed25519 -N "" -q &&
		cat ~/.ssh/id_ed25519.pub | sudo tee -a ~/.ssh/authorized_keys &&
		sudo chmod 700 ~/.ssh &&
		sudo chmod 600 ~/.ssh/authorized_keys &&
		sudo chown -R %s:%s ~/.ssh
	`, username, username, username, username)

	_, err := c.Execute(cmd)
	if err != nil {
		return "", fmt.Errorf("failed to generate SSH key: %w", err)
	}

	return c.GetSSHKey(username)
}

// GetSSHKey gets the public SSH key for a user
func (c *Client) GetSSHKey(username string) (string, error) {
	if c.client == nil {
		return "", fmt.Errorf("not connected")
	}

	cmd := fmt.Sprintf("sudo cat /home/%s/.ssh/id_ed25519.pub", username)
	output, err := c.Execute(cmd)
	if err != nil {
		return "", fmt.Errorf("failed to get SSH key: %w", err)
	}

	return strings.TrimSpace(output), nil
}

// isValidUsername validates a username
func isValidUsername(username string) bool {
	if len(username) < 1 || len(username) > 32 {
		return false
	}
	for _, c := range username {
		if !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '_' || c == '-') {
			return false
		}
	}
	return true
}
