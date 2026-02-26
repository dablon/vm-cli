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
		Usage:                "CLI para ejecutar comandos en VM remota via SSH",
		EnableBashCompletion: true,
		Commands: []*cli.Command{
			cmd.NewConnectCommand(),
			cmd.NewExecCommand(),
			cmd.NewUserCreateCommand(),
			cmd.NewUserExistsCommand(),
			cmd.NewUserDeleteCommand(),
			cmd.NewSystemInfoCommand(),
			cmd.NewDockerCommand(),
			cmd.NewInitCommand(),
		},
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "verbose",
				Usage: "Salida verbosa",
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}
