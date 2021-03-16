package cmd

import (
	"github.com/urfave/cli/v2"
)

func (a *App) generateCmd() *cli.Command {
	return &cli.Command{
		Name:   "generate",
		Usage:  "DEPRECATED: Please use `ensure mocks generate`. This command (`ensure generate mocks`) will be removed in the next minor release.",
		Hidden: true,

		Subcommands: []*cli.Command{
			a.generateMocksCmd(),
		},
	}
}

func (a *App) generateMocksCmd() *cli.Command {
	return &cli.Command{
		Name:   "mocks",
		Usage:  "DEPRECATED: Please use `ensure mocks generate`. This command (`ensure generate mocks`) will be removed in the next minor release.",
		Hidden: true,

		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "disable-parallel",
				Usage: "Disables generating the mocks in parallel",
			},
		},

		Action: func(c *cli.Context) error {
			a.Logger.Print(
				"WARNING: `ensure generate mocks` is deprecated. Please use `ensure mocks generate`." +
					"This command (`ensure generate mocks`) will be removed in the next minor release.\n\n",
			)

			pwd, err := a.Getwd()
			if err != nil {
				return err
			}

			config, err := a.EnsureFileLoader.LoadConfig(pwd)
			if err != nil {
				return err
			}

			config.DisableParallelGeneration = c.Bool("disable-parallel")
			return a.MockGenerator.GenerateMocks(a.Cleanup.ToContext(c.Context), config)
		},
	}
}
