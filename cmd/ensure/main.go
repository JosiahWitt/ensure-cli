package main

import (
	"fmt"
	"log"
	"os"

	"bursavich.dev/fs-shim/io/fs"
	"github.com/JosiahWitt/ensure-cli/internal/cmd"
	"github.com/JosiahWitt/ensure-cli/internal/ensurefile"
	"github.com/JosiahWitt/ensure-cli/internal/fswrite"
	"github.com/JosiahWitt/ensure-cli/internal/mockgen"
	"github.com/JosiahWitt/ensure-cli/internal/runcmd"
)

//nolint:gochecknoglobals // Allows injecting the version
// Version of the CLI.
// Should be tied to the release version.
var Version = "0.1.2"

func main() {
	logger := log.New(os.Stdout, "", 0)

	app := cmd.App{
		Version: Version,

		Logger:           logger,
		Getwd:            os.Getwd,
		EnsureFileLoader: &ensurefile.Loader{FS: fs.DirFS("")},
		MockGenerator: &mockgen.Generator{
			CmdRun:  &runcmd.Runner{},
			FSWrite: &fswrite.FSWrite{},
			Logger:  logger,
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Printf("ERROR: %v\n", err) //nolint:forbidigo // Allow printing error messages
		os.Exit(1)
	}
}
