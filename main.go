package main

import (
	"fmt"
	"os"

	"bursavich.dev/fs-shim/io/fs"
	"github.com/JosiahWitt/cli-ensure/internal/cmd"
	"github.com/JosiahWitt/cli-ensure/internal/ensurefile"
	"github.com/JosiahWitt/cli-ensure/internal/fswrite"
	"github.com/JosiahWitt/cli-ensure/internal/mockgen"
	"github.com/JosiahWitt/cli-ensure/internal/runcmd"
)

//nolint:gochecknoglobals // Allows injecting the version
// Version of the CLI.
// Should be tied to the release version.
var Version = "0.1.0"

func main() {
	app := cmd.App{
		Version: Version,

		Getwd:            os.Getwd,
		EnsureFileLoader: &ensurefile.Loader{FS: fs.DirFS("")},
		MockGenerator:    &mockgen.Generator{CmdRun: &runcmd.Runner{}, FSWrite: &fswrite.FSWrite{}},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Printf("ERROR: %v\n", err) //nolint:forbidigo // Allow printing error messages
		os.Exit(1)
	}
}
