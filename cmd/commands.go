package cmd

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"vm-cli/internal/ssh"

	"github.com/urfave/cli/v2"
)

// Profile represents a VM profile
type Profile struct {
	Name     string `json:"name"`
	Host     string `json:"host"`
	User     string `json:"user"`
	Password string `json:"password"`
	Port     string `json:"port"`
}

// ProfileStore manages VM profiles
type ProfileStore struct {
	profiles map[string]Profile
	configPath string
}

// NewProfileStore creates a new profile store
func NewProfileStore() *ProfileStore {
	configDir := os.Getenv("HOME") + "/.vm-cli"
	os.MkdirAll(configDir, 0755)
	
	store := &ProfileStore{
		profiles: make(map[string]Profile),
		configPath: configDir + "/profiles.json",
	}
	
	store.load()
	return store
}

// load loads profiles from disk
func (s *ProfileStore) load() {
	data, err := os.ReadFile(s.configPath)
	if err != nil {
		return
	}
	json.Unmarshal(data, &s.profiles)
}

// save saves profiles to disk
func (s *ProfileStore) save() {
	data, _ := json.MarshalIndent(s.profiles, "", "  ")
	os.WriteFile(s.configPath, data, 0600)
}

// Add saves a new profile
func (s *ProfileStore) Add(name, host, user, password, port string) {
	s.profiles[name] = Profile{name, host, user, password, port}
	s.save()
}

// Get retrieves a profile by name
func (s *ProfileStore) Get(name string) (Profile, bool) {
	p, ok := s.profiles[name]
	return p, ok
}

// Delete removes a profile
func (s *ProfileStore) Delete(name string) {
	delete(s.profiles, name)
	s.save()
}

// List returns all profile names
func (s *ProfileStore) List() []string {
	names := make([]string, 0, len(s.profiles))
	for name := range s.profiles {
		names = append(names, name)
	}
	return names
}

// getConnectionParams returns host, user, password, port from profile or flags
func getConnectionParams(cCtx *cli.Context) (host, user, password, port string, err error) {
	profileName := cCtx.String("profile")
	
	// Check profile first
	if profileName != "" {
		store := NewProfileStore()
		if p, ok := store.Get(profileName); ok {
			return p.Host, p.User, p.Password, p.Port, nil
		}
		return "", "", "", "", fmt.Errorf("profile '%s' not found", profileName)
	}
	
	// Fallback to flags
	host = cCtx.String("host")
	user = cCtx.String("user")
	password = getPassword(cCtx)
	port = cCtx.String("port")
	
	if password == "" {
		return "", "", "", "", fmt.Errorf("password required: --profile or --password or VM_CLI_PASSWORD")
	}
	
	return host, user, password, port, nil
}

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
			&cli.StringFlag{Name: "profile", Usage: "Use a saved profile"},
			&cli.StringFlag{Name: "host", Usage: "VM hostname or IP"},
			&cli.StringFlag{Name: "user", Aliases: []string{"u"}, Usage: "SSH username"},
			&cli.StringFlag{Name: "password", Aliases: []string{"p"}, Usage: "SSH password (or VM_CLI_PASSWORD)"},
			&cli.StringFlag{Name: "port", Value: "22", Aliases: []string{"P"}, Usage: "SSH port"},
		},
		Action: func(cCtx *cli.Context) error {
			host, user, password, port, err := getConnectionParams(cCtx)
			if err != nil {
				return err
			}

			client := ssh.NewClient(host, user, password, port)

			fmt.Printf("🔌 Conectando a %s...\n", host)
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
			&cli.StringFlag{Name: "profile", Usage: "Use a saved profile"},
			&cli.StringFlag{Name: "host", Usage: "VM hostname or IP"},
			&cli.StringFlag{Name: "user", Aliases: []string{"u"}, Usage: "SSH username"},
			&cli.StringFlag{Name: "password", Aliases: []string{"p"}, Usage: "SSH password"},
			&cli.StringFlag{Name: "port", Value: "22", Aliases: []string{"P"}, Usage: "SSH port"},
			&cli.StringFlag{Name: "command", Required: true, Aliases: []string{"c"}, Usage: "Command to execute"},
		},
		Action: func(cCtx *cli.Context) error {
			host, user, password, port, err := getConnectionParams(cCtx)
			if err != nil {
				return err
			}

			client := ssh.NewClient(host, user, password, port)

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

// NewDockerCommand returns the docker command
func NewDockerCommand() *cli.Command {
	return &cli.Command{
		Name:  "docker",
		Usage: "Docker management commands",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "profile", Usage: "Use a saved profile"},
			&cli.StringFlag{Name: "host", Usage: "VM hostname or IP"},
			&cli.StringFlag{Name: "user", Aliases: []string{"u"}, Usage: "SSH username"},
			&cli.StringFlag{Name: "password", Aliases: []string{"p"}, Usage: "SSH password"},
			&cli.StringFlag{Name: "port", Value: "22", Aliases: []string{"P"}, Usage: "SSH port"},
		},
		Subcommands: []*cli.Command{
			{
				Name:  "ps",
				Usage: "List containers",
				Action: func(cCtx *cli.Context) error {
					host, user, password, port, err := getConnectionParams(cCtx)
					if err != nil {
						return err
					}
					client := ssh.NewClient(host, user, password, port)
					if err := client.Connect(); err != nil {
						return err
					}
					defer client.Close()
					output, _ := client.Execute("docker ps")
					fmt.Print(output)
					return nil
				},
			},
			{
				Name:  "logs",
				Usage: "Get container logs",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "container", Required: true, Aliases: []string{"c"}, Usage: "Container name"},
					&cli.StringFlag{Name: "lines", Aliases: []string{"n"}, Value: "50", Usage: "Number of lines"},
				},
				Action: func(cCtx *cli.Context) error {
					host, user, password, port, err := getConnectionParams(cCtx)
					if err != nil {
						return err
					}
					client := ssh.NewClient(host, user, password, port)
					if err := client.Connect(); err != nil {
						return err
					}
					defer client.Close()
					lines := cCtx.String("lines")
					container := cCtx.String("container")
					output, _ := client.Execute(fmt.Sprintf("docker logs --tail %s %s", lines, container))
					fmt.Print(output)
					return nil
				},
			},
		},
	}
}

// NewUserCreateCommand returns the user-create command
func NewUserCreateCommand() *cli.Command {
	return &cli.Command{
		Name:  "user-create",
		Usage: "Create a new user on the remote VM",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "profile", Usage: "Use a saved profile"},
			&cli.StringFlag{Name: "host", Usage: "VM hostname or IP"},
			&cli.StringFlag{Name: "user", Aliases: []string{"u"}, Usage: "SSH username"},
			&cli.StringFlag{Name: "password", Aliases: []string{"p"}, Usage: "SSH password"},
			&cli.StringFlag{Name: "port", Value: "22", Aliases: []string{"P"}, Usage: "SSH port"},
			&cli.StringFlag{Name: "new-user", Aliases: []string{"nu"}, Required: true, Usage: "New username"},
			&cli.StringFlag{Name: "new-password", Aliases: []string{"np"}, Required: true, Usage: "New user password"},
		},
		Action: func(cCtx *cli.Context) error {
			host, user, password, port, err := getConnectionParams(cCtx)
			if err != nil {
				return err
			}

			client := ssh.NewClient(host, user, password, port)
			if err := client.Connect(); err != nil {
				return err
			}
			defer client.Close()

			newUser := cCtx.String("new-user")
			newPass := cCtx.String("new-password")

			cmd := fmt.Sprintf("sudo useradd -m -s /bin/bash %s && echo '%s:%s' | sudo chpasswd", newUser, newUser, newPass)
			output, err := client.Execute(cmd)
			if err != nil {
				fmt.Printf("Output: %s\n", output)
				return err
			}

			fmt.Printf("✅ Usuario '%s' creado!\n", newUser)
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
			&cli.StringFlag{Name: "profile", Usage: "Use a saved profile"},
			&cli.StringFlag{Name: "host", Usage: "VM hostname or IP"},
			&cli.StringFlag{Name: "user", Aliases: []string{"u"}, Usage: "SSH username"},
			&cli.StringFlag{Name: "password", Aliases: []string{"p"}, Usage: "SSH password"},
			&cli.StringFlag{Name: "port", Value: "22", Aliases: []string{"P"}, Usage: "SSH port"},
			&cli.StringFlag{Name: "check-user", Aliases: []string{"cu"}, Required: true, Usage: "User to check"},
		},
		Action: func(cCtx *cli.Context) error {
			host, user, password, port, err := getConnectionParams(cCtx)
			if err != nil {
				return err
			}

			client := ssh.NewClient(host, user, password, port)
			if err := client.Connect(); err != nil {
				return err
			}
			defer client.Close()

			checkUser := cCtx.String("check-user")
			cmd := fmt.Sprintf("id %s", checkUser)
			output, err := client.Execute(cmd)

			if err != nil {
				fmt.Printf("❌ Usuario '%s' NO existe\n", checkUser)
				return nil
			}

			fmt.Printf("✅ Usuario '%s' existe: %s\n", checkUser, output)
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
			&cli.StringFlag{Name: "profile", Usage: "Use a saved profile"},
			&cli.StringFlag{Name: "host", Usage: "VM hostname or IP"},
			&cli.StringFlag{Name: "user", Aliases: []string{"u"}, Usage: "SSH username"},
			&cli.StringFlag{Name: "password", Aliases: []string{"p"}, Usage: "SSH password"},
			&cli.StringFlag{Name: "port", Value: "22", Aliases: []string{"P"}, Usage: "SSH port"},
			&cli.StringFlag{Name: "delete-user", Aliases: []string{"du"}, Required: true, Usage: "User to delete"},
		},
		Action: func(cCtx *cli.Context) error {
			host, user, password, port, err := getConnectionParams(cCtx)
			if err != nil {
				return err
			}

			client := ssh.NewClient(host, user, password, port)
			if err := client.Connect(); err != nil {
				return err
			}
			defer client.Close()

			deleteUser := cCtx.String("delete-user")
			cmd := fmt.Sprintf("sudo userdel -r %s", deleteUser)
			_, err = client.Execute(cmd)

			if err != nil {
				fmt.Printf("❌ Error al eliminar usuario: %v\n", err)
				return err
			}

			fmt.Printf("✅ Usuario '%s' eliminado!\n", deleteUser)
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
			&cli.StringFlag{Name: "profile", Usage: "Use a saved profile"},
			&cli.StringFlag{Name: "host", Usage: "VM hostname or IP"},
			&cli.StringFlag{Name: "user", Aliases: []string{"u"}, Usage: "SSH username"},
			&cli.StringFlag{Name: "password", Aliases: []string{"p"}, Usage: "SSH password"},
			&cli.StringFlag{Name: "port", Value: "22", Aliases: []string{"P"}, Usage: "SSH port"},
		},
		Action: func(cCtx *cli.Context) error {
			host, user, password, port, err := getConnectionParams(cCtx)
			if err != nil {
				return err
			}

			client := ssh.NewClient(host, user, password, port)
			if err := client.Connect(); err != nil {
				return err
			}
			defer client.Close()

			cmds := []string{
				"echo '=== UPTIME ===' && uptime",
				"echo '=== MEMORY ===' && free -h",
				"echo '=== DISK ===' && df -h",
				"echo '=== CPU ===' && lscpu | grep 'Model name'",
			}

			for _, cmd := range cmds {
				output, _ := client.Execute(cmd)
				fmt.Print(output)
				fmt.Println()
			}

			return nil
		},
	}
}

// NewInitCommand returns the init command
func NewInitCommand() *cli.Command {
	return &cli.Command{
		Name:  "init",
		Usage: "Initialize configuration file",
		Action: func(cCtx *cli.Context) error {
			configDir := os.Getenv("HOME") + "/.vm-cli"
			os.MkdirAll(configDir, 0755)
			configFile := configDir + "/profiles.json"
			
			if _, err := os.ReadFile(configFile); os.IsNotExist(err) {
				os.WriteFile(configFile, []byte("{}"), 0644)
				fmt.Println("✅ Config initialized at:", configFile)
			} else {
				fmt.Println("✅ Config already exists at:", configFile)
			}
			return nil
		},
	}
}

// NewCopyCommand returns the copy command
func NewCopyCommand() *cli.Command {
	return &cli.Command{
		Name:  "copy",
		Usage: "Copy files between local and remote VM",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "profile", Usage: "Use a saved profile"},
			&cli.StringFlag{Name: "host", Usage: "VM hostname or IP"},
			&cli.StringFlag{Name: "user", Aliases: []string{"u"}, Usage: "SSH username"},
			&cli.StringFlag{Name: "password", Aliases: []string{"p"}, Usage: "SSH password"},
			&cli.StringFlag{Name: "port", Value: "22", Aliases: []string{"P"}, Usage: "SSH port"},
			&cli.StringFlag{Name: "source", Aliases: []string{"s"}, Required: true, Usage: "Source file (use user@host:path for remote)"},
			&cli.StringFlag{Name: "dest", Aliases: []string{"d"}, Required: true, Usage: "Destination path"},
			&cli.BoolFlag{Name: "to-remote", Aliases: []string{"to"}, Usage: "Copy to remote (default is from remote)"},
		},
		Action: func(cCtx *cli.Context) error {
			host, user, password, port, err := getConnectionParams(cCtx)
			if err != nil {
				return err
			}

			source := cCtx.String("source")
			dest := cCtx.String("dest")
			toRemote := cCtx.Bool("to-remote")

			client := ssh.NewClient(host, user, password, port)
			if err := client.Connect(); err != nil {
				return err
			}
			defer client.Close()

			if toRemote {
				// Local to Remote: read local file, encode base64, decode on remote
				data, err := os.ReadFile(source)
				if err != nil {
					return fmt.Errorf("failed to read local file: %v", err)
				}
				encoded := base64.StdEncoding.EncodeToString(data)
				
				// Write base64 to remote and decode
				cmd := "echo '" + encoded + "' | base64 -d > " + dest
				_, err = client.Execute(cmd)
				if err != nil {
					return err
				}
				fmt.Printf("✅ Copied %s -> %s@%s:%s\n", source, user, host, dest)
			} else {
				// Remote to Local: read remote file, encode base64, decode locally
				cmd := "base64 " + source
				output, err := client.Execute(cmd)
				if err != nil {
					return fmt.Errorf("failed to read remote file: %v", err)
				}
				
				encoded := strings.TrimSpace(output)
				data, err := base64.StdEncoding.DecodeString(encoded)
				if err != nil {
					return fmt.Errorf("failed to decode: %v", err)
				}
				
				err = os.WriteFile(dest, data, 0644)
				if err != nil {
					return fmt.Errorf("failed to write local file: %v", err)
				}
				fmt.Printf("✅ Copied %s@%s:%s -> %s\n", user, host, source, dest)
			}

			return nil
		},
	}
}
