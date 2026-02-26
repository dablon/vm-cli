package cmd

import (
	"fmt"
	"os"

	"vm-cli/internal/ssh"

	"github.com/urfave/cli/v2"
)

func NewConnectCommand() *cli.Command {
	return &cli.Command{
		Name:  "connect",
		Usage: "Conectar a la VM remota",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "host", Required: true, Usage: "Host de la VM"},
			&cli.StringFlag{Name: "user", Required: true, Usage: "Usuario SSH"},
			&cli.StringFlag{Name: "password", Usage: "Contraseña SSH (o VM_CLI_PASSWORD)"},
			&cli.StringFlag{Name: "port", Value: "22", Usage: "Puerto SSH"},
		},
		Action: func(cCtx *cli.Context) error {
			password := getPassword(cCtx)
			if password == "" {
				return fmt.Errorf("se requiere contraseña: --password o VM_CLI_PASSWORD")
			}
			client := ssh.NewClient(cCtx.String("host"), cCtx.String("user"), password, cCtx.String("port"))
			fmt.Println("🔌 Conectando a", cCtx.String("host")+"...")
			if err := client.Connect(); err != nil {
				return fmt.Errorf("❌ Error de conexión: %w", err)
			}
			defer client.Close()
			fmt.Println("✅ ¡Conectado!")
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

func NewExecCommand() *cli.Command {
	return &cli.Command{
		Name:  "exec",
		Usage: "Ejecutar un comando en la VM",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "host", Required: true},
			&cli.StringFlag{Name: "user", Required: true},
			&cli.StringFlag{Name: "password", Usage: "o VM_CLI_PASSWORD"},
			&cli.StringFlag{Name: "port", Value: "22"},
			&cli.StringFlag{Name: "command", Required: true, Usage: "Comando a ejecutar"},
		},
		Action: func(cCtx *cli.Context) error {
			password := getPassword(cCtx)
			if password == "" {
				return fmt.Errorf("se requiere contraseña")
			}
			client := ssh.NewClient(cCtx.String("host"), cCtx.String("user"), password, cCtx.String("port"))
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

func NewUserCreateCommand() *cli.Command {
	return &cli.Command{
		Name:  "user-create",
		Usage: "Crear un nuevo usuario en la VM",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "host", Required: true},
			&cli.StringFlag{Name: "user", Required: true},
			&cli.StringFlag{Name: "password", Usage: "o VM_CLI_PASSWORD"},
			&cli.StringFlag{Name: "port", Value: "22"},
			&cli.StringFlag{Name: "new-user", Required: true},
			&cli.StringFlag{Name: "new-password", Required: true},
		},
		Action: func(cCtx *cli.Context) error {
			password := getPassword(cCtx)
			if password == "" {
				return fmt.Errorf("se requiere contraseña")
			}
			client := ssh.NewClient(cCtx.String("host"), cCtx.String("user"), password, cCtx.String("port"))
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

func NewUserExistsCommand() *cli.Command {
	return &cli.Command{
		Name:  "user-exists",
		Usage: "Verificar si un usuario existe",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "host", Required: true},
			&cli.StringFlag{Name: "user", Required: true},
			&cli.StringFlag{Name: "password", Usage: "o VM_CLI_PASSWORD"},
			&cli.StringFlag{Name: "port", Value: "22"},
			&cli.StringFlag{Name: "check-user", Required: true},
		},
		Action: func(cCtx *cli.Context) error {
			password := getPassword(cCtx)
			if password == "" {
				return fmt.Errorf("se requiere contraseña")
			}
			client := ssh.NewClient(cCtx.String("host"), cCtx.String("user"), password, cCtx.String("port"))
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

func NewUserDeleteCommand() *cli.Command {
	return &cli.Command{
		Name:  "user-delete",
		Usage: "Eliminar un usuario de la VM",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "host", Required: true},
			&cli.StringFlag{Name: "user", Required: true},
			&cli.StringFlag{Name: "password", Usage: "o VM_CLI_PASSWORD"},
			&cli.StringFlag{Name: "port", Value: "22"},
			&cli.StringFlag{Name: "username", Required: true},
		},
		Action: func(cCtx *cli.Context) error {
			password := getPassword(cCtx)
			if password == "" {
				return fmt.Errorf("se requiere contraseña")
			}
			client := ssh.NewClient(cCtx.String("host"), cCtx.String("user"), password, cCtx.String("port"))
			if err := client.Connect(); err != nil {
				return err
			}
			defer client.Close()
			username := cCtx.String("username")
			fmt.Printf("🗑️  Eliminando usuario %s...\n", username)
			_, err := client.DeleteUser(username)
			if err != nil {
				return fmt.Errorf("❌ Error al eliminar usuario: %w", err)
			}
			fmt.Printf("✅ Usuario %s eliminado\n", username)
			return nil
		},
	}
}

func NewSystemInfoCommand() *cli.Command {
	return &cli.Command{
		Name:  "sysinfo",
		Usage: "Obtener información del sistema",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "host", Required: true},
			&cli.StringFlag{Name: "user", Required: true},
			&cli.StringFlag{Name: "password", Usage: "o VM_CLI_PASSWORD"},
			&cli.StringFlag{Name: "port", Value: "22"},
		},
		Action: func(cCtx *cli.Context) error {
			password := getPassword(cCtx)
			if password == "" {
				return fmt.Errorf("se requiere contraseña")
			}
			client := ssh.NewClient(cCtx.String("host"), cCtx.String("user"), password, cCtx.String("port"))
			if err := client.Connect(); err != nil {
				return err
			}
			defer client.Close()
			fmt.Println("📊 Obteniendo información del sistema...")
			output, err := client.GetSystemInfo()
			if err != nil {
				return fmt.Errorf("❌ Error: %w", err)
			}
			fmt.Println(output)
			return nil
		},
	}
}

func NewDockerCommand() *cli.Command {
	return &cli.Command{
		Name:  "docker",
		Usage: "Ejecutar comandos Docker",
		Subcommands: []*cli.Command{
			{
				Name:  "ps",
				Usage: "				Flags: []cli.FlagListar contenedores",
{
					&cli.StringFlag{Name: "host", Required: true},
					&cli.StringFlag{Name: "user", Required: true},
					&cli.StringFlag{Name: "password", Usage: "o VM_CLI_PASSWORD"},
					&cli.BoolFlag{Name: "all", Usage: "Mostrar todos"},
				},
				Action: func(cCtx *cli.Context) error {
					password := getPassword(cCtx)
					if password == "" {
						return fmt.Errorf("se requiere contraseña")
					}
					client := ssh.NewClient(cCtx.String("host"), cCtx.String("user"), password, "22")
					if err := client.Connect(); err != nil {
						return err
					}
					defer client.Close()
					output, err := client.ListContainers(cCtx.Bool("all"))
					if err != nil {
						return err
					}
					fmt.Print(output)
					return nil
				},
			},
			{
				Name:  "info",
				Usage: "Información de Docker",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "host", Required: true},
					&cli.StringFlag{Name: "user", Required: true},
					&cli.StringFlag{Name: "password", Usage: "o VM_CLI_PASSWORD"},
				},
				Action: func(cCtx *cli.Context) error {
					password := getPassword(cCtx)
					if password == "" {
						return fmt.Errorf("se requiere contraseña")
					}
					client := ssh.NewClient(cCtx.String("host"), cCtx.String("user"), password, "22")
					if err := client.Connect(); err != nil {
						return err
					}
					defer client.Close()
					output, err := client.GetDockerInfo()
					if err != nil {
						return err
					}
					fmt.Print(output)
					return nil
				},
			},
		},
	}
}

func NewInitCommand() *cli.Command {
	return &cli.Command{
		Name:  "init",
		Usage: "Inicializar configuración",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "host", Required: true},
			&cli.StringFlag{Name: "user", Required: true},
			&cli.StringFlag{Name: "password", Usage: "o VM_CLI_PASSWORD"},
			&cli.StringFlag{Name: "port", Value: "22"},
			&cli.StringFlag{Name: "agent", Required: true},
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
				return fmt.Errorf("❌ Error al guardar config: %w", err)
			}
			fmt.Printf("✅ Configuración guardada en %s\n", path)
			return nil
		},
	}
}

func getPassword(cCtx *cli.Context) string {
	if val := cCtx.String("password"); val != "" {
		return val
	}
	return os.Getenv("VM_CLI_PASSWORD")
}
