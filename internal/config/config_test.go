package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestConfig_SaveAndLoad(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "config.json")

	cfg := &Config{
		VMHost:     "192.168.1.100",
		VMUser:     "testuser",
		VMPort:     "22",
		AgentName:  "test-agent",
		RemoteUser: "remote-test",
	}

	// Save
	err := cfg.Save(tmpFile)
	if err != nil {
		t.Fatalf("Failed to save config: %v", err)
	}

	// Load
	loaded, err := Load(tmpFile)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Verify
	if loaded.VMHost != cfg.VMHost {
		t.Errorf("VMHost mismatch: got %s, want %s", loaded.VMHost, cfg.VMHost)
	}
	if loaded.VMUser != cfg.VMUser {
		t.Errorf("VMUser mismatch: got %s, want %s", loaded.VMUser, cfg.VMUser)
	}
	if loaded.VMPort != cfg.VMPort {
		t.Errorf("VMPort mismatch: got %s, want %s", loaded.VMPort, cfg.VMPort)
	}
	if loaded.AgentName != cfg.AgentName {
		t.Errorf("AgentName mismatch: got %s, want %s", loaded.AgentName, cfg.AgentName)
	}
	if loaded.RemoteUser != cfg.RemoteUser {
		t.Errorf("RemoteUser mismatch: got %s, want %s", loaded.RemoteUser, cfg.RemoteUser)
	}
}

func TestLoad_FileNotFound(t *testing.T) {
	_, err := Load("/nonexistent/path/config.json")
	if err == nil {
		t.Error("Expected error for nonexistent file")
	}
}

func TestLoad_InvalidJSON(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "invalid.json")

	os.WriteFile(tmpFile, []byte("invalid json"), 0644)

	_, err := Load(tmpFile)
	if err == nil {
		t.Error("Expected error for invalid JSON")
	}
}

func TestDefaultConfigPath(t *testing.T) {
	path := DefaultConfigPath()
	if path == "" {
		t.Error("DefaultConfigPath returned empty string")
	}
	// Should contain .vm-cli/config.json
	if filepath.Base(path) != "config.json" {
		t.Errorf("Expected config.json as basename, got %s", filepath.Base(path))
	}
}

func TestConfig_Save_InvalidPath(t *testing.T) {
	cfg := &Config{
		VMHost: "test",
	}
	// Try to save to invalid path (directory doesn't exist)
	err := cfg.Save("/nonexistent directory/config.json")
	if err == nil {
		t.Error("Expected error for invalid path")
	}
}

func TestConfig_Empty(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "empty.json")

	cfg := &Config{}
	err := cfg.Save(tmpFile)
	if err != nil {
		t.Fatalf("Failed to save empty config: %v", err)
	}

	loaded, err := Load(tmpFile)
	if err != nil {
		t.Fatalf("Failed to load empty config: %v", err)
	}

	if loaded.VMHost != "" {
		t.Errorf("Expected empty VMHost, got %s", loaded.VMHost)
	}
}
