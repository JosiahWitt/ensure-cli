package cmd

import (
	"log"

	"github.com/JosiahWitt/ensure-cli/internal/ensurefile"
	"github.com/JosiahWitt/ensure-cli/internal/exitcleanup"
	"github.com/JosiahWitt/ensure-cli/internal/mockgen"
	"github.com/urfave/cli/v2"
)

// App is the CLI application for ensure.
type App struct {
	Version string

	Logger           *log.Logger
	Getwd            func() (string, error)
	EnsureFileLoader ensurefile.LoaderIface
	MockGenerator    mockgen.MockGenerator
	Cleanup          exitcleanup.ExitCleaner
}

// Run the application given the os.Args array.
func (a *App) Run(args []string) error {
	cliApp := &cli.App{
		Name:    "ensure",
		Usage:   "A balanced test framework for Go 1.14+.",
		Version: a.Version,

		ExitErrHandler: func(context *cli.Context, err error) {}, // Bubble up error

		Commands: []*cli.Command{
			a.generateCmd(),
			a.mocksCmd(),
		},
	}

	return cliApp.Run(args)
}
