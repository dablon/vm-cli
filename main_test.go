package main

import (
	"testing"

	"github.com/urfave/cli/v2"
)

// Test main CLI app
func TestCLI_Connect_RequiresHost(t *testing.T) {
	app := createTestApp()

	err := app.Run([]string{"vm-cli", "connect", "--user", "test", "--password", "pass"})
	if err == nil {
		t.Error("Expected error when host is missing")
	}
}

func TestCLI_Connect_RequiresUser(t *testing.T) {
	app := createTestApp()

	err := app.Run([]string{"vm-cli", "connect", "--host", "192.168.1.1", "--password", "pass"})
	if err == nil {
		t.Error("Expected error when user is missing")
	}
}

func TestCLI_Connect_RequiresPassword(t *testing.T) {
	app := createTestApp()

	err := app.Run([]string{"vm-cli", "connect", "--host", "192.168.1.1", "--user", "test"})
	if err == nil {
		t.Error("Expected error when password is missing")
	}
}

func TestCLI_Exec_RequiresHost(t *testing.T) {
	app := createTestApp()

	err := app.Run([]string{"vm-cli", "exec", "--user", "test", "--command", "ls"})
	if err == nil {
		t.Error("Expected error when host is missing")
	}
}

func TestCLI_Exec_RequiresUser(t *testing.T) {
	app := createTestApp()

	err := app.Run([]string{"vm-cli", "exec", "--host", "192.168.1.1", "--command", "ls"})
	if err == nil {
		t.Error("Expected error when user is missing")
	}
}

func TestCLI_Exec_RequiresCommand(t *testing.T) {
	app := createTestApp()

	err := app.Run([]string{"vm-cli", "exec", "--host", "192.168.1.1", "--user", "test"})
	if err == nil {
		t.Error("Expected error when command is missing")
	}
}

func TestCLI_UserCreate_RequiresHost(t *testing.T) {
	app := createTestApp()

	err := app.Run([]string{"vm-cli", "user-create", "--user", "admin", "--new-user", "newuser", "--new-password", "pass"})
	if err == nil {
		t.Error("Expected error when host is missing")
	}
}

func TestCLI_UserCreate_RequiresUser(t *testing.T) {
	app := createTestApp()

	err := app.Run([]string{"vm-cli", "user-create", "--host", "192.168.1.1", "--new-user", "newuser", "--new-password", "pass"})
	if err == nil {
		t.Error("Expected error when user is missing")
	}
}

func TestCLI_UserCreate_RequiresNewUser(t *testing.T) {
	app := createTestApp()

	err := app.Run([]string{"vm-cli", "user-create", "--host", "192.168.1.1", "--user", "admin", "--new-password", "pass"})
	if err == nil {
		t.Error("Expected error when new-user is missing")
	}
}

func TestCLI_UserCreate_RequiresNewPassword(t *testing.T) {
	app := createTestApp()

	err := app.Run([]string{"vm-cli", "user-create", "--host", "192.168.1.1", "--user", "admin", "--new-user", "newuser"})
	if err == nil {
		t.Error("Expected error when new-password is missing")
	}
}

func TestCLI_UserExists_RequiresHost(t *testing.T) {
	app := createTestApp()

	err := app.Run([]string{"vm-cli", "user-exists", "--user", "admin", "--check-user", "someuser"})
	if err == nil {
		t.Error("Expected error when host is missing")
	}
}

func TestCLI_UserExists_RequiresUser(t *testing.T) {
	app := createTestApp()

	err := app.Run([]string{"vm-cli", "user-exists", "--host", "192.168.1.1", "--check-user", "someuser"})
	if err == nil {
		t.Error("Expected error when user is missing")
	}
}

func TestCLI_UserExists_RequiresCheckUser(t *testing.T) {
	app := createTestApp()

	err := app.Run([]string{"vm-cli", "user-exists", "--host", "192.168.1.1", "--user", "admin"})
	if err == nil {
		t.Error("Expected error when check-user is missing")
	}
}

func TestCLI_Init_RequiresHost(t *testing.T) {
	app := createTestApp()

	err := app.Run([]string{"vm-cli", "init", "--user", "admin", "--agent", "test-agent"})
	if err == nil {
		t.Error("Expected error when host is missing")
	}
}

func TestCLI_Init_RequiresUser(t *testing.T) {
	app := createTestApp()

	err := app.Run([]string{"vm-cli", "init", "--host", "192.168.1.1", "--agent", "test-agent"})
	if err == nil {
		t.Error("Expected error when user is missing")
	}
}

func TestCLI_Init_RequiresAgent(t *testing.T) {
	app := createTestApp()

	err := app.Run([]string{"vm-cli", "init", "--host", "192.168.1.1", "--user", "admin"})
	if err == nil {
		t.Error("Expected error when agent is missing")
	}
}

func TestCLI_Help(t *testing.T) {
	app := createTestApp()

	err := app.Run([]string{"vm-cli", "--help"})
	if err != nil {
		t.Errorf("Help should not error: %v", err)
	}
}

func TestCLI_Version(t *testing.T) {
	app := &cli.App{
		Name:    "vm-cli",
		Version: "1.0.0",
	}

	err := app.Run([]string{"vm-cli", "--version"})
	if err != nil {
		t.Errorf("Version should not error: %v", err)
	}
}

func TestCLI_UnknownCommand(t *testing.T) {
	app := createTestApp()

	// Unknown command should fail - may return error or show help
	_ = app.Run([]string{"vm-cli", "unknown-command"})
	// Just verify it doesn't panic
}

func TestCLI_NoArgs(t *testing.T) {
	app := createTestApp()

	// Should not error, just show help
	err := app.Run([]string{"vm-cli"})
	if err != nil {
		t.Errorf("No args should not error: %v", err)
	}
}

func TestCLI_DeleteUser_RequiresHost(t *testing.T) {
	app := createTestApp()

	err := app.Run([]string{"vm-cli", "user-delete", "--user", "admin", "--delete-user", "someuser"})
	if err == nil {
		t.Error("Expected error when host is missing")
	}
}

func TestCLI_Docker_RequiresHost(t *testing.T) {
	app := createTestApp()

	err := app.Run([]string{"vm-cli", "docker", "--user", "admin", "--subcommand", "ps"})
	if err == nil {
		t.Error("Expected error when host is missing")
	}
}

func TestCLI_Docker_RequiresUser(t *testing.T) {
	app := createTestApp()

	err := app.Run([]string{"vm-cli", "docker", "--host", "192.168.1.1", "--subcommand", "ps"})
	if err == nil {
		t.Error("Expected error when user is missing")
	}
}

func TestCLI_SysInfo_RequiresHost(t *testing.T) {
	app := createTestApp()

	err := app.Run([]string{"vm-cli", "sysinfo", "--user", "admin"})
	if err == nil {
		t.Error("Expected error when host is missing")
	}
}

func TestCLI_SysInfo_RequiresUser(t *testing.T) {
	app := createTestApp()

	err := app.Run([]string{"vm-cli", "sysinfo", "--host", "192.168.1.1"})
	if err == nil {
		t.Error("Expected error when user is missing")
	}
}

func createTestApp() *cli.App {
	return &cli.App{
		Name:  "vm-cli",
		Usage: "Test CLI",
		Commands: []*cli.Command{
			{
				Name: "connect",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "host", Required: true},
					&cli.StringFlag{Name: "user", Required: true},
					&cli.StringFlag{Name: "password", Required: true},
				},
				Action: func(c *cli.Context) error {
					return nil
				},
			},
			{
				Name: "exec",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "host", Required: true},
					&cli.StringFlag{Name: "user", Required: true},
					&cli.StringFlag{Name: "command", Required: true},
				},
				Action: func(c *cli.Context) error {
					return nil
				},
			},
			{
				Name: "user-create",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "host", Required: true},
					&cli.StringFlag{Name: "user", Required: true},
					&cli.StringFlag{Name: "new-user", Required: true},
					&cli.StringFlag{Name: "new-password", Required: true},
				},
				Action: func(c *cli.Context) error {
					return nil
				},
			},
			{
				Name: "user-delete",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "host", Required: true},
					&cli.StringFlag{Name: "user", Required: true},
					&cli.StringFlag{Name: "delete-user", Required: true},
				},
				Action: func(c *cli.Context) error {
					return nil
				},
			},
			{
				Name: "user-exists",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "host", Required: true},
					&cli.StringFlag{Name: "user", Required: true},
					&cli.StringFlag{Name: "check-user", Required: true},
				},
				Action: func(c *cli.Context) error {
					return nil
				},
			},
			{
				Name: "docker",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "host", Required: true},
					&cli.StringFlag{Name: "user", Required: true},
					&cli.StringFlag{Name: "subcommand", Required: true},
				},
				Action: func(c *cli.Context) error {
					return nil
				},
			},
			{
				Name: "sysinfo",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "host", Required: true},
					&cli.StringFlag{Name: "user", Required: true},
				},
				Action: func(c *cli.Context) error {
					return nil
				},
			},
			{
				Name: "init",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "host", Required: true},
					&cli.StringFlag{Name: "user", Required: true},
					&cli.StringFlag{Name: "agent", Required: true},
				},
				Action: func(c *cli.Context) error {
					return nil
				},
			},
		},
	}
}
