package main

import (
	"fmt"
	"os"

	"vm-cli/cmd"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:                 "vm-cli",
		Usage:                "CLI for remote VM management via SSH",
		EnableBashCompletion: true,
		Commands: []*cli.Command{
			{
				Name:  "profile-save",
				Usage: "Save a VM profile",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "name", Aliases: []string{"n"}, Required: true, Usage: "Profile name"},
					&cli.StringFlag{Name: "host", Required: true, Usage: "VM hostname or IP"},
					&cli.StringFlag{Name: "user", Aliases: []string{"u"}, Required: true, Usage: "SSH username"},
					&cli.StringFlag{Name: "password", Aliases: []string{"p"}, Required: true, Usage: "SSH password"},
					&cli.StringFlag{Name: "port", Value: "22", Usage: "SSH port"},
				},
				Action: func(c *cli.Context) error {
					store := cmd.NewProfileStore()
					store.Add(c.String("name"), c.String("host"), c.String("user"), c.String("password"), c.String("port"))
					fmt.Printf("✅ Profile '%s' saved!\n", c.String("name"))
					return nil
				},
			},
			{
				Name:  "profile-list",
				Usage: "List all saved profiles",
				Action: func(c *cli.Context) error {
					store := cmd.NewProfileStore()
					profiles := store.List()
					if len(profiles) == 0 { fmt.Println("No profiles saved."); return nil }
					fmt.Println("📋 Saved profiles:")
					for _, name := range profiles {
						p, _ := store.Get(name)
						fmt.Printf("  • %s → %s@%s:%s\n", name, p.User, p.Host, p.Port)
					}
					return nil
				},
			},
			{
				Name:  "profile-delete",
				Usage: "Delete a profile",
				Flags: []cli.Flag{&cli.StringFlag{Name: "name", Aliases: []string{"n"}, Required: true, Usage: "Profile name"}},
				Action: func(c *cli.Context) error {
					store := cmd.NewProfileStore()
					store.Delete(c.String("name"))
					fmt.Printf("✅ Profile '%s' deleted!\n", c.String("name"))
					return nil
				},
			},
			cmd.NewConnectCommand(),
			cmd.NewExecCommand(),
			cmd.NewUserCreateCommand(),
			cmd.NewUserExistsCommand(),
			cmd.NewUserDeleteCommand(),
			cmd.NewSystemInfoCommand(),
			cmd.NewDockerCommand(),
			cmd.NewInitCommand(),
			cmd.NewCopyCommand(),
		},
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "profile", Usage: "Use saved profile instead of host/user/password flags"},
			&cli.BoolFlag{Name: "verbose", Usage: "Verbose output"},
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}
