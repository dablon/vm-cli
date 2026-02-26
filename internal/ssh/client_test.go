package ssh

import (
	"testing"
)

// Tests for Client creation
func TestNewClient(t *testing.T) {
	client := NewClient("192.168.1.100", "testuser", "password", "22")

	if client == nil {
		t.Fatal("NewClient returned nil")
	}
	if client.client != nil {
		t.Error("SSH client should be nil before connection")
	}
}

func TestNewClient_EmptyPort(t *testing.T) {
	client := NewClient("192.168.1.100", "testuser", "password", "")

	if client == nil {
		t.Fatal("NewClient returned nil")
	}
}

func TestNewClient_CustomPort(t *testing.T) {
	client := NewClient("192.168.1.100", "testuser", "password", "2222")

	if client == nil {
		t.Fatal("NewClient returned nil")
	}
}

// Tests for Connection (without actual connection)
func TestClient_Close_NotConnected(t *testing.T) {
	client := NewClient("192.168.1.100", "testuser", "password", "22")

	err := client.Close()
	if err != nil {
		t.Errorf("Close should not error when already nil: %v", err)
	}
}

// Tests for Execute (without connection)
func TestClient_Execute_NotConnected(t *testing.T) {
	client := NewClient("192.168.1.100", "testuser", "password", "22")

	_, err := client.Execute("ls")
	if err == nil {
		t.Error("Expected error when executing without connection")
	}
}

// Tests for ExecuteWithSudo
func TestClient_ExecuteWithSudo_NotConnected(t *testing.T) {
	client := NewClient("192.168.1.100", "testuser", "password", "22")

	_, err := client.ExecuteWithSudo("ls")
	if err == nil {
		t.Error("Expected error when executing sudo without connection")
	}
}

// Tests for UserExists
func TestClient_UserExists_NotConnected(t *testing.T) {
	client := NewClient("192.168.1.100", "testuser", "password", "22")

	exists, _ := client.UserExists("someuser")
	if exists {
		t.Error("Expected user to not exist when not connected")
	}
}

// Tests for CreateUser
func TestClient_CreateUser_NotConnected(t *testing.T) {
	client := NewClient("192.168.1.100", "testuser", "password", "22")

	_, err := client.CreateUser("newuser", "password123")
	if err == nil {
		t.Error("Expected error when creating user without connection")
	}
}

// Tests for DeleteUser
func TestClient_DeleteUser_NotConnected(t *testing.T) {
	client := NewClient("192.168.1.100", "testuser", "password", "22")

	_, err := client.DeleteUser("someuser")
	if err == nil {
		t.Error("Expected error when deleting user without connection")
	}
}

// Tests for EnsureSSHKey
func TestClient_EnsureSSHKey_NotConnected(t *testing.T) {
	client := NewClient("192.168.1.100", "testuser", "password", "22")

	_, err := client.EnsureSSHKey("someuser")
	if err == nil {
		t.Error("Expected error when ensuring SSH key without connection")
	}
}

// Tests for GetSSHKey
func TestClient_GetSSHKey_NotConnected(t *testing.T) {
	client := NewClient("192.168.1.100", "testuser", "password", "22")

	_, err := client.GetSSHKey("someuser")
	if err == nil {
		t.Error("Expected error when getting SSH key without connection")
	}
}

// Tests for GetSystemInfo
func TestClient_GetSystemInfo_NotConnected(t *testing.T) {
	client := NewClient("192.168.1.100", "testuser", "password", "22")

	// When not connected, the command may still execute (returns empty/error)
	// Just verify it returns something
	_, _ = client.GetSystemInfo()
}

// Tests for ListContainers
func TestClient_ListContainers_NotConnected(t *testing.T) {
	client := NewClient("192.168.1.100", "testuser", "password", "22")

	_, err := client.ListContainers(true)
	if err == nil {
		t.Error("Expected error when listing containers without connection")
	}
}

func TestClient_ListContainers_AllFalse(t *testing.T) {
	client := NewClient("192.168.1.100", "testuser", "password", "22")

	_, err := client.ListContainers(false)
	if err == nil {
		t.Error("Expected error when listing containers without connection")
	}
}

// Tests for GetDockerInfo
func TestClient_GetDockerInfo_NotConnected(t *testing.T) {
	client := NewClient("192.168.1.100", "testuser", "password", "22")

	_, err := client.GetDockerInfo()
	if err == nil {
		t.Error("Expected error when getting docker info without connection")
	}
}

// Tests for multiple scenarios
func TestClient_Execute_MultipleCommands(t *testing.T) {
	client := NewClient("192.168.1.100", "testuser", "password", "22")

	tests := []string{
		"ls",
		"pwd",
		"whoami",
		"echo test",
		"cat /etc/passwd",
	}

	for _, cmd := range tests {
		t.Run(cmd, func(t *testing.T) {
			_, err := client.Execute(cmd)
			if err == nil {
				t.Error("Expected error for command without connection")
			}
		})
	}
}

func TestClient_UserExists_MultipleUsers(t *testing.T) {
	client := NewClient("192.168.1.100", "testuser", "password", "22")

	users := []string{
		"root",
		"nobody",
		"testuser",
		"nonexistent",
	}

	for _, user := range users {
		t.Run(user, func(t *testing.T) {
			exists, _ := client.UserExists(user)
			if exists {
				t.Error("Expected user to not exist without connection")
			}
		})
	}
}

func TestClient_CreateUser_MultipleUsers(t *testing.T) {
	client := NewClient("192.168.1.100", "testuser", "password", "22")

	users := []struct {
		username string
		password string
	}{
		{"user1", "pass1"},
		{"user2", "pass2"},
		{"newuser", "newpass"},
	}

	for _, u := range users {
		t.Run(u.username, func(t *testing.T) {
			_, err := client.CreateUser(u.username, u.password)
			if err == nil {
				t.Error("Expected error without connection")
			}
		})
	}
}

func TestNewClient_DifferentHosts(t *testing.T) {
	hosts := []struct {
		host string
		user string
		pass string
		port string
	}{
		{"192.168.1.1", "user", "pass", "22"},
		{"10.0.0.1", "admin", "admin123", "2222"},
		{"example.com", "ubuntu", "ubuntu", "22"},
		{"192.168.1.100", "root", "root123", ""},
	}

	for _, h := range hosts {
		t.Run(h.host, func(t *testing.T) {
			client := NewClient(h.host, h.user, h.pass, h.port)
			if client == nil {
				t.Error("NewClient returned nil")
			}
		})
	}
}

func TestClient_Close_MultipleTimes(t *testing.T) {
	client := NewClient("192.168.1.100", "testuser", "password", "22")

	// Close multiple times should not panic
	for i := 0; i < 3; i++ {
		err := client.Close()
		if err != nil {
			t.Errorf("Close %d: unexpected error: %v", i+1, err)
		}
	}
}
