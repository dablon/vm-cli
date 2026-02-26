package ssh

import (
	"testing"
)

func TestNewClient(t *testing.T) {
	client := NewClient("192.168.1.1", "user", "password", "22")
	
	if client.config.Host != "192.168.1.1" {
		t.Errorf("Expected host 192.168.1.1, got %s", client.config.Host)
	}
	
	if client.config.User != "user" {
		t.Errorf("Expected user user, got %s", client.config.User)
	}
	
	if client.config.Password != "password" {
		t.Errorf("Expected password, got %s", client.config.Password)
	}
	
	if client.config.Port != "22" {
		t.Errorf("Expected port 22, got %s", client.config.Port)
	}
}

func TestNewClientDefaultPort(t *testing.T) {
	client := NewClient("192.168.1.1", "user", "password", "")
	
	if client.config.Port != "22" {
		t.Errorf("Expected default port 22, got %s", client.config.Port)
	}
}

func TestIsValidUsername(t *testing.T) {
	tests := []struct {
		username string
		valid    bool
	}{
		{"root", true},
		{"user123", true},
		{"test-user", true},
		{"test_user", true},
		{"a", true},
		{"user with spaces", false},
		{"user@special", false},
		{"", false},
		{"verylongusernamethatexceeds32characters", false},
	}

	for _, tt := range tests {
		result := isValidUsername(tt.username)
		if result != tt.valid {
			t.Errorf("isValidUsername(%s) = %v, expected %v", tt.username, result, tt.valid)
		}
	}
}

func TestClientNotConnected(t *testing.T) {
	client := NewClient("192.168.1.1", "user", "password", "22")
	
	// Test Execute without connection
	_, err := client.Execute("ls")
	if err == nil {
		t.Error("Expected error when executing without connection")
	}
	
	// Test Close without connection
	err = client.Close()
	if err != nil {
		t.Errorf("Unexpected error when closing disconnected client: %v", err)
	}
}
