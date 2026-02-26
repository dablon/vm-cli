package cmd

import (
	"fmt"
	"os"

	"vm-cli/internal/ssh"

	"github.com/urfave/cli/v2"
)

// getPassword returns password from flag or environment variable
func getPassword(cCtx *cli.Context) string {
	if val := cCtx.String("password"); val != "" {
		return val
	}
	if val := os.Getenv("VM_CLI_PASSWORD"); val != "" {
		return val
	}
	return ""
}

// NewConnectCommand returns the connect command
func NewConnectCommand() *cli.Command {
	return &cli.Command{
		Name:  "connect",
		Usage: "Connect to remote VM and run test command",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "host", Required: true, Aliases: []string{"h"}, Usage: "VM hostname or IP"},
			&cli.StringFlag{Name: "user", Required: true, Aliases: []string{"u"}, Usage: "SSH username"},
			&cli.StringFlag{Name: "password", Aliases: []string{"p"}, Usage: "SSH password (or VM_CLI_PASSWORD)"},
			&cli.StringFlag{Name: "port", Value: "22", Aliases: []string{"P"}, Usage: "SSH port"},
		},
		Action: func(cCtx *cli.Context) error {
			password := getPassword(cCtx)
			if password == "" {
				return fmt.Errorf("password required: --password or VM_CLI_PASSWORD")
			}

			client := ssh.NewClient(
				cCtx.String("host"),
				cCtx.String("user"),
				password,
				cCtx.String("port"),
			)

			fmt.Printf("🔌 Conectando a %s...\n", cCtx.String("host"))
			if err := client.Connect(); err != nil {
				return fmt.Errorf("❌ Error de conexión: %w", err)
			}
			defer client.Close()

			fmt.Println("✅ ¡Conectado!")

			fmt.Println("📟 Ejecutando uname -a...")
			output, err := client.Execute("uname -a")
			if err != nil {
				return fmt.Errorf("❌ Error al ejecutar: %w", err)
			}

			fmt.Println("📺 Output:")
			fmt.Println(output)
			return nil
		},
	}
}

// NewExecCommand returns the exec command
func NewExecCommand() *cli.Command {
	return &cli.Command{
		Name:  "exec",
		Usage: "Execute a command on the remote VM",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "host", Required: true, Aliases: []string{"h"}},
			&cli.StringFlag{Name: "user", Required: true, Aliases: []string{"u"}},
			&cli.StringFlag{Name: "password", Aliases: []string{"p"}},
			&cli.StringFlag{Name: "port", Value: "22", Aliases: []string{"P"}},
			&cli.StringFlag{Name: "command", Required: true, Aliases: []string{"c"}, Usage: "Command to execute"},
		},
		Action: func(cCtx *cli.Context) error {
			password := getPassword(cCtx)
			if password == "" {
				return fmt.Errorf("password required: --password or VM_CLI_PASSWORD")
			}

			client := ssh.NewClient(
				cCtx.String("host"),
				cCtx.String("user"),
				password,
				cCtx.String("port"),
			)

			if err := client.Connect(); err != nil {
				return err
			}
			defer client.Close()

			output, err := client.Execute(cCtx.String("command"))
			if err != nil {
				return err
			}

			fmt.Print(output)
			return nil
		},
	}
}

// NewUserCreateCommand returns the user-create command
func NewUserCreateCommand() *cli.Command {
	return &cli.Command{
		Name:  "user-create",
		Usage: "Create a new user on the remote VM",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "host", Required: true, Aliases: []string{"h"}},
			&cli.StringFlag{Name: "user", Required: true, Aliases: []string{"u"}},
			&cli.StringFlag{Name: "password", Aliases: []string{"p"}},
			&cli.StringFlag{Name: "port", Value: "22", Aliases: []string{"P"}},
			&cli.StringFlag{Name: "new-user", Required: true, Aliases: []string{"n"}, Usage: "New username"},
			&cli.StringFlag{Name: "new-password", Required: true, Aliases: []string{"w"}, Usage: "Password for new user"},
		},
		Action: func(cCtx *cli.Context) error {
			password := getPassword(cCtx)
			if password == "" {
				return fmt.Errorf("password required: --password or VM_CLI_PASSWORD")
			}

			client := ssh.NewClient(
				cCtx.String("host"),
				cCtx.String("user"),
				password,
				cCtx.String("port"),
			)

			if err := client.Connect(); err != nil {
				return err
			}
			defer client.Close()

			newUser := cCtx.String("new-user")
			newPass := cCtx.String("new-password")

			fmt.Printf("👤 Creando usuario %s...\n", newUser)
			_, err := client.CreateUser(newUser, newPass)
			if err != nil {
				return fmt.Errorf("❌ Error al crear usuario: %w", err)
			}

			fmt.Printf("✅ Usuario %s creado\n", newUser)
			return nil
		},
	}
}

// NewUserExistsCommand returns the user-exists command
func NewUserExistsCommand() *cli.Command {
	return &cli.Command{
		Name:  "user-exists",
		Usage: "Check if a user exists on the remote VM",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "host", Required: true, Aliases: []string{"h"}},
			&cli.StringFlag{Name: "user", Required: true, Aliases: []string{"u"}},
			&cli.StringFlag{Name: "password", Aliases: []string{"p"}},
			&cli.StringFlag{Name: "port", Value: "22", Aliases: []string{"P"}},
			&cli.StringFlag{Name: "check-user", Required: true, Aliases: []string{"c"}, Usage: "Username to check"},
		},
		Action: func(cCtx *cli.Context) error {
			password := getPassword(cCtx)
			if password == "" {
				return fmt.Errorf("password required: --password or VM_CLI_PASSWORD")
			}

			client := ssh.NewClient(
				cCtx.String("host"),
				cCtx.String("user"),
				password,
				cCtx.String("port"),
			)

			if err := client.Connect(); err != nil {
				return err
			}
			defer client.Close()

			exists, err := client.UserExists(cCtx.String("check-user"))
			if err != nil {
				return err
			}

			if exists {
				fmt.Printf("✅ El usuario %s existe\n", cCtx.String("check-user"))
			} else {
				fmt.Printf("❌ El usuario %s no existe\n", cCtx.String("check-user"))
			}
			return nil
		},
	}
}

// NewUserDeleteCommand returns the user-delete command
func NewUserDeleteCommand() *cli.Command {
	return &cli.Command{
		Name:  "user-delete",
		Usage: "Delete a user from the remote VM",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "host", Required: true, Aliases: []string{"h"}},
			&cli.StringFlag{Name: "user", Required: true, Aliases: []string{"u"}},
			&cli.StringFlag{Name: "password", Aliases: []string{"p"}},
			&cli.StringFlag{Name: "port", Value: "22", Aliases: []string{"P"}},
			&cli.StringFlag{Name: "username", Required: true, Aliases: []string{"n"}, Usage: "Username to delete"},
		},
		Action: func(cCtx *cli.Context) error {
			password := getPassword(cCtx)
			if password == "" {
				return fmt.Errorf("password required: --password or VM_CLI_PASSWORD")
			}

			client := ssh.NewClient(
				cCtx.String("host"),
				cCtx.String("user"),
				password,
				cCtx.String("port"),
			)

			if err := client.Connect(); err != nil {
				return err
			}
			defer client.Close()

			username := cCtx.String("username")
			fmt.Printf("🗑️ Eliminando usuario %s...\n", username)
			_, err := client.DeleteUser(username)
			if err != nil {
				return fmt.Errorf("❌ Error al eliminar usuario: %w", err)
			}

			fmt.Printf("✅ Usuario %s eliminado\n", username)
			return nil
		},
	}
}

// NewSystemInfoCommand returns the sysinfo command
func NewSystemInfoCommand() *cli.Command {
	return &cli.Command{
		Name:  "sysinfo",
		Usage: "Get system information from the remote VM",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "host", Required: true, Aliases: []string{"h"}},
			&cli.StringFlag{Name: "user", Required: true, Aliases: []string{"u"}},
			&cli.StringFlag{Name: "password", Aliases: []string{"p"}},
			&cli.StringFlag{Name: "port", Value: "22", Aliases: []string{"P"}},
		},
		Action: func(cCtx *cli.Context) error {
			password := getPassword(cCtx)
			if password == "" {
				return fmt.Errorf("password required: --password or VM_CLI_PASSWORD")
			}

			client := ssh.NewClient(
				cCtx.String("host"),
				cCtx.String("user"),
				password,
				cCtx.String("port"),
			)

			if err := client.Connect(); err != nil {
				return err
			}
			defer client.Close()

			output, err := client.GetSystemInfo()
			if err != nil {
				return fmt.Errorf("failed to get system info: %w", err)
			}

			fmt.Println(output)
			return nil
		},
	}
}

// NewDockerCommand returns the docker command
func NewDockerCommand() *cli.Command {
	return &cli.Command{
		Name:  "docker",
		Usage: "Docker management commands",
		Subcommands: []*cli.Command{
			{
				Name:  "ps",
				Usage: "List containers",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "host", Required: true, Aliases: []string{"h"}},
					&cli.StringFlag{Name: "user", Required: true, Aliases: []string{"u"}},
					&cli.StringFlag{Name: "password", Aliases: []string{"p"}},
					&cli.StringFlag{Name: "port", Value: "22", Aliases: []string{"P"}},
					&cli.BoolFlag{Name: "all", Aliases: []string{"a"}, Usage: "Show all containers"},
				},
				Action: func(cCtx *cli.Context) error {
					password := getPassword(cCtx)
					if password == "" {
						return fmt.Errorf("password required: --password or VM_CLI_PASSWORD")
					}

					client := ssh.NewClient(
						cCtx.String("host"),
						cCtx.String("user"),
						password,
						cCtx.String("port"),
					)

					if err := client.Connect(); err != nil {
						return err
					}
					defer client.Close()

					output, err := client.ListContainers(cCtx.Bool("all"))
					if err != nil {
						return fmt.Errorf("failed to list containers: %w", err)
					}

					fmt.Print(output)
					return nil
				},
			},
			{
				Name:  "info",
				Usage: "Show Docker info",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "host", Required: true, Aliases: []string{"h"}},
					&cli.StringFlag{Name: "user", Required: true, Aliases: []string{"u"}},
					&cli.StringFlag{Name: "password", Aliases: []string{"p"}},
					&cli.StringFlag{Name: "port", Value: "22", Aliases: []string{"P"}},
				},
				Action: func(cCtx *cli.Context) error {
					password := getPassword(cCtx)
					if password == "" {
						return fmt.Errorf("password required: --password or VM_CLI_PASSWORD")
					}

					client := ssh.NewClient(
						cCtx.String("host"),
						cCtx.String("user"),
						password,
						cCtx.String("port"),
					)

					if err := client.Connect(); err != nil {
						return err
					}
					defer client.Close()

					output, err := client.GetDockerInfo()
					if err != nil {
						return fmt.Errorf("failed to get docker info: %w", err)
					}

					fmt.Print(output)
					return nil
				},
			},
		},
	}
}

// NewInitCommand returns the init command
func NewInitCommand() *cli.Command {
	return &cli.Command{
		Name:  "init",
		Usage: "Initialize configuration file",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "host", Required: true, Aliases: []string{"h"}},
			&cli.StringFlag{Name: "user", Required: true, Aliases: []string{"u"}},
			&cli.StringFlag{Name: "password", Aliases: []string{"p"}},
			&cli.StringFlag{Name: "port", Value: "22", Aliases: []string{"P"}},
			&cli.StringFlag{Name: "agent", Required: true, Aliases: []string{"a"}, Usage: "Agent name"},
		},
		Action: func(cCtx *cli.Context) error {
			cfg := &ssh.Config{
				VMHost:    cCtx.String("host"),
				VMUser:    cCtx.String("user"),
				VMPort:    cCtx.String("port"),
				AgentName: cCtx.String("agent"),
			}

			path := ssh.DefaultConfigPath()
			if err := cfg.Save(path); err != nil {
				return fmt.Errorf("failed to save config: %w", err)
			}

			fmt.Printf("✅ Configuración guardada en %s\n", path)
			return nil
		},
	}
}
