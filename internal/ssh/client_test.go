package ssh

import (
	"testing"
)

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

func TestClient_Close_NotConnected(t *testing.T) {
	client := NewClient("192.168.1.100", "testuser", "password", "22")

	err := client.Close()
	if err != nil {
		t.Errorf("Close should not error when already nil: %v", err)
	}
}

func TestClient_Execute_NotConnected(t *testing.T) {
	client := NewClient("192.168.1.100", "testuser", "password", "22")

	_, err := client.Execute("ls")
	if err == nil {
		t.Error("Expected error when executing without connection")
	}
}

func TestClient_ExecuteWithSudo_NotConnected(t *testing.T) {
	client := NewClient("192.168.1.100", "testuser", "password", "22")

	_, err := client.ExecuteWithSudo("ls")
	if err == nil {
		t.Error("Expected error when executing sudo without connection")
	}
}

func TestClient_UserExists_NotConnected(t *testing.T) {
	client := NewClient("192.168.1.100", "testuser", "password", "22")

	// When not connected, id command will fail
	exists, _ := client.UserExists("someuser")
	if exists {
		t.Error("Expected user to not exist when not connected")
	}
}

func TestClient_CreateUser_NotConnected(t *testing.T) {
	client := NewClient("192.168.1.100", "testuser", "password", "22")

	_, err := client.CreateUser("newuser", "password123")
	if err == nil {
		t.Error("Expected error when creating user without connection")
	}
}

func TestClient_EnsureSSHKey_NotConnected(t *testing.T) {
	client := NewClient("192.168.1.100", "testuser", "password", "22")

	_, err := client.EnsureSSHKey("someuser")
	if err == nil {
		t.Error("Expected error when ensuring SSH key without connection")
	}
}

func TestClient_GetSSHKey_NotConnected(t *testing.T) {
	client := NewClient("192.168.1.100", "testuser", "password", "22")

	_, err := client.GetSSHKey("someuser")
	if err == nil {
		t.Error("Expected error when getting SSH key without connection")
	}
}

func TestSaveHostKey(t *testing.T) {
	err := SaveHostKey("192.168.1.100", "test-key-fingerprint")
	if err != nil {
		t.Errorf("SaveHostKey failed: %v", err)
	}
}
