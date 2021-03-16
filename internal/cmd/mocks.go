package cmd

import (
	"github.com/urfave/cli/v2"
)

func (a *App) mocksCmd() *cli.Command {
	return &cli.Command{
		Name:  "mocks",
		Usage: "commands related to mocks",

		Subcommands: []*cli.Command{
			a.mocksGenerateCmd(),
			a.mocksTidyCmd(),
		},
	}
}

func (a *App) mocksGenerateCmd() *cli.Command {
	return &cli.Command{
		Name:  "generate",
		Usage: "generates GoMocks (https://github.com/golang/mock) for the packages and interfaces listed in .ensure.yml",

		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "disable-parallel",
				Usage: "Disables generating the mocks in parallel",
			},
		},

		Action: func(c *cli.Context) error {
			pwd, err := a.Getwd()
			if err != nil {
				return err
			}

			config, err := a.EnsureFileLoader.LoadConfig(pwd)
			if err != nil {
				return err
			}

			config.DisableParallelGeneration = c.Bool("disable-parallel")
			if err := a.MockGenerator.GenerateMocks(a.Cleanup.ToContext(c.Context), config); err != nil {
				return err
			}

			if config.Mocks.TidyAfterGenerate {
				if err := a.MockGenerator.TidyMocks(config); err != nil {
					return err
				}
			}

			return nil
		},
	}
}

func (a *App) mocksTidyCmd() *cli.Command {
	return &cli.Command{
		Name:  "tidy",
		Usage: "removes any files and directories that would not be generated for the packages and interfaces listed in .ensure.yml",

		Action: func(c *cli.Context) error {
			pwd, err := a.Getwd()
			if err != nil {
				return err
			}

			config, err := a.EnsureFileLoader.LoadConfig(pwd)
			if err != nil {
				return err
			}

			return a.MockGenerator.TidyMocks(config)
		},
	}
}
