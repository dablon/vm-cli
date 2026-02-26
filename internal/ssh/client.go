package ssh

import (
	"encoding/json"
	"fmt"
	"os"

	"golang.org/x/crypto/ssh"
)

// Config represents CLI configuration
type Config struct {
	VMHost    string `json:"vm_host"`
	VMUser    string `json:"vm_user"`
	VMPort    string `json:"vm_port"`
	AgentName string `json:"agent_name"`
}

// Save saves config to file
func (c *Config) Save(path string) error {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	dir := path[:len(path)-len("/config.json")]
	os.MkdirAll(dir, 0755)
	return os.WriteFile(path, data, 0644)
}

// Load loads config from file
func (c *Config) Load(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, c)
}

// DefaultConfigPath returns default config path
func DefaultConfigPath() string {
	home, _ := os.UserHomeDir()
	return home + "/.vm-cli/config.json"
}

// Client represents an SSH connection
type Client struct {
	Host     string
	User     string
	Password string
	Port     string
	client   *ssh.Client
}

// NewClient creates a new SSH client
func NewClient(host, user, password, port string) *Client {
	if port == "" {
		port = "22"
	}
	return &Client{
		Host:     host,
		User:     user,
		Password: password,
		Port:     port,
	}
}

// Connect establishes SSH connection
func (c *Client) Connect() error {
	if c.client != nil {
		return nil
	}

	authMethods := []ssh.AuthMethod{
		ssh.Password(c.Password),
	}

	config := &ssh.ClientConfig{
		User:            c.User,
		Auth:            authMethods,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	addr := fmt.Sprintf("%s:%s", c.Host, c.Port)
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
		return c.client.Close()
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

// ExecuteWithSudo executes a command with sudo
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
		fmt.Printf("Warning: PTY not available\n")
	}

	sudoCmd := fmt.Sprintf("sudo %s", command)
	output, err := session.CombinedOutput(sudoCmd)
	if err != nil {
		return "", fmt.Errorf("sudo failed: %w", err)
	}

	return string(output), nil
}

// CreateUser creates a new user on the VM
func (c *Client) CreateUser(username, password string) (string, error) {
	if username == "" || password == "" {
		return "", fmt.Errorf("username and password required")
	}

	cmd := fmt.Sprintf("/usr/sbin/useradd -m -s /bin/bash %s && echo '%s:%s' | /usr/sbin/chpasswd",
		username, username, password)

	output, err := c.ExecuteWithSudo(cmd)
	if err != nil {
		return "", fmt.Errorf("failed to create user: %w", err)
	}
	return output, nil
}

// UserExists checks if a user exists
func (c *Client) UserExists(username string) (bool, error) {
	if username == "" {
		return false, fmt.Errorf("username required")
	}

	_, err := c.Execute(fmt.Sprintf("id %s", username))
	if err != nil {
		return false, nil
	}
	return true, nil
}

// DeleteUser deletes a user from the VM
func (c *Client) DeleteUser(username string) (string, error) {
	if username == "" {
		return "", fmt.Errorf("username required")
	}

	cmd := fmt.Sprintf("/usr/sbin/userdel -r %s", username)
	output, err := c.ExecuteWithSudo(cmd)
	if err != nil {
		return "", fmt.Errorf("failed to delete user: %w", err)
	}
	return output, nil
}

// EnsureSSHKey generates SSH key for user
func (c *Client) EnsureSSHKey(username string) (string, error) {
	if username == "" {
		return "", fmt.Errorf("username required")
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

// GetSSHKey gets public SSH key for user
func (c *Client) GetSSHKey(username string) (string, error) {
	if username == "" {
		return "", fmt.Errorf("username required")
	}

	cmd := fmt.Sprintf("sudo cat /home/%s/.ssh/id_ed25519.pub", username)
	output, err := c.Execute(cmd)
	if err != nil {
		return "", fmt.Errorf("failed to get SSH key: %w", err)
	}
	return output, nil
}

// GetSystemInfo gets system information
func (c *Client) GetSystemInfo() (string, error) {
	commands := []string{"uname -a", "df -h", "free -h", "uptime"}
	var output string

	for _, cmd := range commands {
		result, _ := c.Execute(cmd)
		output += result + "\n"
	}

	return output, nil
}

// ListContainers lists Docker containers
func (c *Client) ListContainers(all bool) (string, error) {
	flag := ""
	if all {
		flag = "-a"
	}
	cmd := fmt.Sprintf("docker ps %s", flag)
	output, err := c.Execute(cmd)
	if err != nil {
		return "", fmt.Errorf("failed to list containers: %w", err)
	}
	return output, nil
}

// GetDockerInfo gets Docker information
func (c *Client) GetDockerInfo() (string, error) {
	output, err := c.Execute("docker info 2>/dev/null | head -20")
	if err != nil {
		return "", fmt.Errorf("docker not available: %w", err)
	}
	return output, nil
}
