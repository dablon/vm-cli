package cmd

import (
	"fmt"
	"os"

	"vm-cli/internal/ssh"

	"github.com/urfave/cli/v2"
)

// NewConnectCommand retorna el comando de conexión
func NewConnectCommand() *cli.Command {
	return &cli.Command{
		Name:  "connect",
		Usage: "Conectar a la VM remota y ejecutar comando de prueba",
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

			client := ssh.NewClient(
				cCtx.String("host"),
				cCtx.String("user"),
				password,
				cCtx.String("port"),
			)

			fmt.Println("🔌 Conectando a", cCtx.String("host")+"...")
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

// NewExecCommand retorna el comando de ejecución
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
				return fmt.Errorf("se requiere contraseña: --password o VM_CLI_PASSWORD")
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

// NewUserCreateCommand retorna el comando de creación de usuario
func NewUserCreateCommand() *cli.Command {
	return &cli.Command{
		Name:  "user-create",
		Usage: "Crear un nuevo usuario en la VM",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "host", Required: true},
			&cli.StringFlag{Name: "user", Required: true},
			&cli.StringFlag{Name: "password", Usage: "o VM_CLI_PASSWORD"},
			&cli.StringFlag{Name: "port", Value: "22"},
			&cli.StringFlag{Name: "new-user", Required: true, Usage: "Nombre del nuevo usuario"},
			&cli.StringFlag{Name: "new-password", Required: true, Usage: "Contraseña del nuevo usuario"},
		},
		Action: func(cCtx *cli.Context) error {
			password := getPassword(cCtx)
			if password == "" {
				return fmt.Errorf("se requiere contraseña: --password o VM_CLI_PASSWORD")
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

			// Generar SSH key
			fmt.Printf("🔑 Generando SSH key para %s...\n", newUser)
			_, err = client.EnsureSSHKey(newUser)
			if err != nil {
				fmt.Printf("⚠️  Advertencia: no se pudo generar SSH key: %v\n", err)
			} else {
				fmt.Printf("✅ SSH key generada\n")
			}

			// Mostrar SSH key pública
			pubKey, err := client.GetSSHKey(newUser)
			if err == nil && pubKey != "" && pubKey != "no-key\n" {
				fmt.Println("📋 SSH Public Key:")
				fmt.Println(pubKey)
			}

			return nil
		},
	}
}

// NewUserExistsCommand retorna el comando de verificación de usuario
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
				return fmt.Errorf("se requiere contraseña: --password o VM_CLI_PASSWORD")
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

// NewUserDeleteCommand retorna el comando de eliminación de usuario
func NewUserDeleteCommand() *cli.Command {
	return &cli.Command{
		Name:  "user-delete",
		Usage: "Eliminar un usuario de la VM",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "host", Required: true},
			&cli.StringFlag{Name: "user", Required: true},
			&cli.StringFlag{Name: "password", Usage: "o VM_CLI_PASSWORD"},
			&cli.StringFlag{Name: "port", Value: "22"},
			&cli.StringFlag{Name: "username", Required: true, Usage: "Usuario a eliminar"},
		},
		Action: func(cCtx *cli.Context) error {
			password := getPassword(cCtx)
			if password == "" {
				return fmt.Errorf("se requiere contraseña: --password o VM_CLI_PASSWORD")
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

// NewSystemInfoCommand retorna el comando de información del sistema
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
				return fmt.Errorf("se requiere contraseña: --password o VM_CLI_PASSWORD")
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

// NewDockerCommand retorna el comando de Docker
func NewDockerCommand() *cli.Command {
	return &cli.Command{
		Name:  "docker",
		Usage: "Ejecutar comandos Docker",
		Subcommands: []*cli.Command{
			{
				Name:  "ps",
				Usage: "Listar contenedores",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "host", Required: true},
					&cli.StringFlag{Name: "user", Required: true},
					&cli.StringFlag{Name: "password", Usage: "o VM_CLI_PASSWORD"},
					&cli.BoolFlag{Name: "all", Usage: "Mostrar todos los contenedores"},
				},
				Action: func(cCtx *cli.Context) error {
					password := getPassword(cCtx)
					if password == "" {
						return fmt.Errorf("se requiere contraseña")
					}

					client := ssh.NewClient(
						cCtx.String("host"),
						cCtx.String("user"),
						password,
						"22",
					)

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

					client := ssh.NewClient(
						cCtx.String("host"),
						cCtx.String("user"),
						password,
						"22",
					)

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

// NewInitCommand retorna el comando de inicialización
func NewInitCommand() *cli.Command {
	return &cli.Command{
		Name:  "init",
		Usage: "Inicializar configuración",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "host", Required: true},
			&cli.StringFlag{Name: "user", Required: true},
			&cli.StringFlag{Name: "password", Usage: "o VM_CLI_PASSWORD"},
			&cli.StringFlag{Name: "port", Value: "22"},
			&cli.StringFlag{Name: "agent", Required: true, Usage: "Nombre del agente"},
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

// getPassword retorna la contraseña desde flag o variable de entorno
func getPassword(cCtx *cli.Context) string {
	if val := cCtx.String("password"); val != "" {
		return val
	}
	return os.Getenv("VM_CLI_PASSWORD")
}
