package main

import (
	"fmt"
	"os"

	"vm-cli/internal/config"
	"vm-cli/internal/ssh"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "vm-cli",
		Usage: "CLI para ejecutar comandos en VM remota via SSH",
		Commands: []*cli.Command{
			{
				Name:  "connect",
				Usage: "Conectar a la VM remota",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "host", Required: true, Usage: "Host de la VM"},
					&cli.StringFlag{Name: "user", Required: true, Usage: "Usuario SSH"},
					&cli.StringFlag{Name: "password", Required: true, Usage: "Contraseña SSH"},
					&cli.StringFlag{Name: "port", Value: "22", Usage: "Puerto SSH"},
				},
				Action: func(cCtx *cli.Context) error {
					client := ssh.NewClient(
						cCtx.String("host"),
						cCtx.String("user"),
						cCtx.String("password"),
						cCtx.String("port"),
					)
					
					fmt.Println("🔌 Conectando a", cCtx.String("host")+"...")
					if err := client.Connect(); err != nil {
						return fmt.Errorf("❌ Error de conexión: %w", err)
					}
					defer client.Close()
					
					fmt.Println("✅ ¡Conectado!")
					
					// Ejecutar comando de prueba
					fmt.Println("📟 Ejecutando uname -a...")
					output, err := client.Execute("uname -a")
					if err != nil {
						return fmt.Errorf("❌ Error al ejecutar: %w", err)
					}
					
					fmt.Println("📺 Output:")
					fmt.Println(output)
					return nil
				},
			},
			{
				Name:  "exec",
				Usage: "Ejecutar un comando en la VM",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "host", Required: true},
					&cli.StringFlag{Name: "user", Required: true},
					&cli.StringFlag{Name: "password", Required: true},
					&cli.StringFlag{Name: "port", Value: "22"},
					&cli.StringFlag{Name: "command", Required: true, Usage: "Comando a ejecutar"},
				},
				Action: func(cCtx *cli.Context) error {
					client := ssh.NewClient(
						cCtx.String("host"),
						cCtx.String("user"),
						cCtx.String("password"),
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
			},
			{
				Name:  "user-create",
				Usage: "Crear un nuevo usuario en la VM",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "host", Required: true},
					&cli.StringFlag{Name: "user", Required: true},
					&cli.StringFlag{Name: "password", Required: true},
					&cli.StringFlag{Name: "port", Value: "22"},
					&cli.StringFlag{Name: "new-user", Required: true, Usage: "Nombre del nuevo usuario"},
					&cli.StringFlag{Name: "new-password", Required: true, Usage: "Contraseña del nuevo usuario"},
				},
				Action: func(cCtx *cli.Context) error {
					client := ssh.NewClient(
						cCtx.String("host"),
						cCtx.String("user"),
						cCtx.String("password"),
						cCtx.String("port"),
					)
					
					if err := client.Connect(); err != nil {
						return err
					}
					defer client.Close()
					
					newUser := cCtx.String("new-user")
					newPass := cCtx.String("new-password")
					
					fmt.Printf("👤 Creando usuario %s...\n", newUser)
					
					output, err := client.CreateUser(newUser, newPass)
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
					if err == nil && pubKey != "" {
						fmt.Println("📋 SSH Public Key:")
						fmt.Println(pubKey)
					}
					
					fmt.Println(output)
					return nil
				},
			},
			{
				Name:  "user-exists",
				Usage: "Verificar si un usuario existe",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "host", Required: true},
					&cli.StringFlag{Name: "user", Required: true},
					&cli.StringFlag{Name: "password", Required: true},
					&cli.StringFlag{Name: "port", Value: "22"},
					&cli.StringFlag{Name: "check-user", Required: true},
				},
				Action: func(cCtx *cli.Context) error {
					client := ssh.NewClient(
						cCtx.String("host"),
						cCtx.String("user"),
						cCtx.String("password"),
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
			},
			{
				Name:  "init",
				Usage: "Inicializar configuración",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "host", Required: true},
					&cli.StringFlag{Name: "user", Required: true},
					&cli.StringFlag{Name: "password", Required: true},
					&cli.StringFlag{Name: "port", Value: "22"},
					&cli.StringFlag{Name: "agent", Required: true, Usage: "Nombre del agente"},
				},
				Action: func(cCtx *cli.Context) error {
					cfg := &config.Config{
						VMHost:    cCtx.String("host"),
						VMUser:    cCtx.String("user"),
						VMPort:    cCtx.String("port"),
						AgentName: cCtx.String("agent"),
					}
					
					path := config.DefaultConfigPath()
					if err := cfg.Save(path); err != nil {
						return fmt.Errorf("❌ Error al guardar config: %w", err)
					}
					
					fmt.Printf("✅ Configuración guardada en %s\n", path)
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}
