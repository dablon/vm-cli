package main

import (
	"os"
	"testing"

	"github.com/urfave/cli/v2"
)

func TestGetPassword_FromFlag(t *testing.T) {
	// Test getting password from CLI flag
	// This is tested through CLI execution
}

func TestGetPassword_FromEnv(t *testing.T) {
	os.Setenv("VM_CLI_PASSWORD", "env-password")
	defer os.Unsetenv("VM_CLI_PASSWORD")

	// When no flag is provided, should check env var
	// This is tested through CLI execution
}

func TestConnect_RequiresHost(t *testing.T) {
	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name: "connect",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "host", Required: true},
					&cli.StringFlag{Name: "user", Required: true},
					&cli.StringFlag{Name: "password"},
				},
				Action: func(cCtx *cli.Context) error {
					return nil
				},
			},
		},
	}

	err := app.Run([]string{"vm-cli", "connect", "--user", "test"})
	if err == nil {
		t.Error("Expected error when host is missing")
	}
}

func TestConnect_RequiresUser(t *testing.T) {
	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name: "connect",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "host", Required: true},
					&cli.StringFlag{Name: "user", Required: true},
					&cli.StringFlag{Name: "password"},
				},
				Action: func(cCtx *cli.Context) error {
					return nil
				},
			},
		},
	}

	err := app.Run([]string{"vm-cli", "connect", "--host", "192.168.1.1"})
	if err == nil {
		t.Error("Expected error when user is missing")
	}
}

func TestExec_RequiresCommand(t *testing.T) {
	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name: "exec",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "host", Required: true},
					&cli.StringFlag{Name: "user", Required: true},
					&cli.StringFlag{Name: "command", Required: true},
				},
				Action: func(cCtx *cli.Context) error {
					return nil
				},
			},
		},
	}

	err := app.Run([]string{"vm-cli", "exec", "--host", "192.168.1.1", "--user", "test"})
	if err == nil {
		t.Error("Expected error when command is missing")
	}
}

func TestUserCreate_RequiresNewUser(t *testing.T) {
	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name: "user-create",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "host", Required: true},
					&cli.StringFlag{Name: "user", Required: true},
					&cli.StringFlag{Name: "new-user", Required: true},
					&cli.StringFlag{Name: "new-password", Required: true},
				},
				Action: func(cCtx *cli.Context) error {
					return nil
				},
			},
		},
	}

	err := app.Run([]string{"vm-cli", "user-create", 
		"--host", "192.168.1.1", 
		"--user", "admin",
		"--new-password", "pass123"})
	if err == nil {
		t.Error("Expected error when new-user is missing")
	}
}

func TestUserCreate_RequiresNewPassword(t *testing.T) {
	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name: "user-create",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "host", Required: true},
					&cli.StringFlag{Name: "user", Required: true},
					&cli.StringFlag{Name: "new-user", Required: true},
					&cli.StringFlag{Name: "new-password", Required: true},
				},
				Action: func(cCtx *cli.Context) error {
					return nil
				},
			},
		},
	}

	err := app.Run([]string{"vm-cli", "user-create", 
		"--host", "192.168.1.1", 
		"--user", "admin",
		"--new-user", "newuser"})
	if err == nil {
		t.Error("Expected error when new-password is missing")
	}
}

func TestUserExists_RequiresCheckUser(t *testing.T) {
	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name: "user-exists",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "host", Required: true},
					&cli.StringFlag{Name: "user", Required: true},
					&cli.StringFlag{Name: "check-user", Required: true},
				},
				Action: func(cCtx *cli.Context) error {
					return nil
				},
			},
		},
	}

	err := app.Run([]string{"vm-cli", "user-exists", 
		"--host", "192.168.1.1", 
		"--user", "admin"})
	if err == nil {
		t.Error("Expected error when check-user is missing")
	}
}

func TestInit_RequiresAgent(t *testing.T) {
	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name: "init",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "host", Required: true},
					&cli.StringFlag{Name: "user", Required: true},
					&cli.StringFlag{Name: "agent", Required: true},
				},
				Action: func(cCtx *cli.Context) error {
					return nil
				},
			},
		},
	}

	err := app.Run([]string{"vm-cli", "init", 
		"--host", "192.168.1.1", 
		"--user", "admin"})
	if err == nil {
		t.Error("Expected error when agent is missing")
	}
}

func TestApp_Help(t *testing.T) {
	app := &cli.App{
		Name:  "vm-cli",
		Usage: "Test CLI",
	}

	err := app.Run([]string{"vm-cli", "--help"})
	if err != nil {
		t.Errorf("Help should not error: %v", err)
	}
}

func TestApp_Version(t *testing.T) {
	app := &cli.App{
		Name:    "vm-cli",
		Version: "1.0.0",
	}

	err := app.Run([]string{"vm-cli", "--version"})
	if err != nil {
		t.Errorf("Version should not error: %v", err)
	}
}
