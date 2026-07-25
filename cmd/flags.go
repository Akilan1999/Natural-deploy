package cmd

import (
	"github.com/urfave/cli/v2"
)

// Variables declared for CLI
var (
	Server    bool
	StartTask string
)

var AppConfigFlags = []cli.Flag{
	&cli.StringFlag{
		Name:        "StartTask",
		Aliases:     []string{"st"},
		Usage:       "Starts the task",
		EnvVars:     []string{"ST"},
		Destination: &StartTask,
	},
	&cli.BoolFlag{
		Name:        "Server",
		Aliases:     []string{"s"},
		Usage:       "Starts the server mode",
		EnvVars:     []string{"S"},
		Destination: &Server,
	},
}
