package mockgen

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"
	"sync"

	"github.com/JosiahWitt/ensure-cli/internal/ensurefile"
	"github.com/JosiahWitt/ensure-cli/internal/fswrite"
	"github.com/JosiahWitt/ensure-cli/internal/runcmd"
	"github.com/JosiahWitt/erk"
	"github.com/JosiahWitt/erk/erg"
)

const (
	defaultPrimaryDestination  = "internal/mocks"
	defaultInternalDestination = "mocks"
	gomockReflectDirPattern    = "gomock_reflect_*"
)

type (
	ErkInvalidConfig    struct{ erk.DefaultKind }
	ErkMultipleFailures struct{ erk.DefaultKind }
	ErkMockGenError     struct{ erk.DefaultKind }
	ErkFSWriteError     struct{ erk.DefaultKind }
)

var (
	ErrMissingMockConfig = erk.New(ErkInvalidConfig{}, "Missing `mocks` config in .ensure.yml file. For example:\n\n"+ensurefile.ExampleFile)
	ErrMissingPackages   = erk.New(ErkInvalidConfig{},
		"No mocks to generate. Please add some to `mocks.packages` in .ensure.yml file. For example:\n\n"+ensurefile.ExampleFile,
	)
	ErrDuplicatePackagePath = erk.New(ErkInvalidConfig{}, "Found duplicate package path: {{.packagePath}}. Package paths must be unique.")

	ErrMissingPackagePath       = erk.New(ErkInvalidConfig{}, "Missing `path` key for package.")
	ErrMissingPackageInterfaces = erk.New(ErkInvalidConfig{},
		"Package '{{.packagePath}}' has no interfaces to generate. Please add them using the `interfaces` key.",
	)

	ErrMultipleGenerationFailures = erk.New(ErkMultipleFailures{}, "Unable to generate at least one mock")
	ErrMockGenFailed              = erk.New(ErkMockGenError{}, "Could not run mockgen successfully for '{{.packageDescription}}': {{.err}}")

	ErrUnableToCreateDir  = erk.New(ErkFSWriteError{}, "Could not create directory '{{.path}}': {{.err}}")
	ErrUnableToCreateFile = erk.New(ErkFSWriteError{}, "Could not create file '{{.path}}': {{.err}}")
)

type GeneratorIface interface {
	GenerateMocks(config *ensurefile.Config, registerCleanup func(func() error)) error
}

type Generator struct {
	CmdRun  runcmd.RunnerIface
	FSWrite fswrite.FSWriteIface
	Logger  *log.Logger
}

var _ GeneratorIface = &Generator{}

// GenerateMocks for the provided configuration.
func (g *Generator) GenerateMocks(config *ensurefile.Config, registerCleanup func(func() error)) error {
	if err := validateConfig(config); err != nil {
		return err
	}

	mockDestinations, err := computeMockDestinations(config)
	if err != nil {
		return err
	}

	for _, pwd := range mockDestinations.uniquePWDs() {
		pwd := pwd // Pin range variable

		registerCleanup(func() error {
			pattern := filepath.Join(pwd, gomockReflectDirPattern)
			g.Logger.Printf("Cleaning up: %s", pattern)
			return g.FSWrite.GlobRemoveAll(pattern)
		})
	}

	asyncParams := &generateMockAsyncParams{
		errors: erg.NewAs(ErrMultipleGenerationFailures),
	}

	g.Logger.Println("Generating mocks:")
	for _, mockDestination := range mockDestinations {
		if !config.DisableParallelGeneration {
			asyncParams.wg.Add(1)
			go g.generateMockAsync(mockDestination, asyncParams)
		} else {
			g.Logger.Printf(" - Generating: %s\n", mockDestination.Package.String())

			if err := g.generateMock(mockDestination); err != nil {
				asyncParams.errors = erg.Append(asyncParams.errors, err)
			}
		}
	}

	asyncParams.wg.Wait()
	if erg.Any(asyncParams.errors) {
		return asyncParams.errors
	}

	return nil
}

type generateMockAsyncParams struct {
	wg       sync.WaitGroup
	errors   error
	errorsMu sync.Mutex
}

func (g *Generator) generateMockAsync(mockDestination *mockDestination, asyncParams *generateMockAsyncParams) {
	defer asyncParams.wg.Done()

	if err := g.generateMock(mockDestination); err != nil {
		asyncParams.errorsMu.Lock()
		defer asyncParams.errorsMu.Unlock()
		asyncParams.errors = erg.Append(asyncParams.errors, err)
		return
	}

	g.Logger.Printf(" - Generated: %s\n", mockDestination.Package.String())
}

func (g *Generator) generateMock(mockDestination *mockDestination) error {
	pkg := mockDestination.Package

	if pkg.Path == "" {
		return ErrMissingPackagePath
	}

	if len(pkg.Interfaces) < 1 {
		return erk.WithParams(ErrMissingPackageInterfaces, erk.Params{
			"packagePath": pkg.Path,
		})
	}

	result, err := g.CmdRun.Exec(&runcmd.ExecParams{
		PWD: mockDestination.PWD,
		CMD: "mockgen", // TODO: Allow overriding
		Args: []string{
			pkg.Path,
			strings.Join(pkg.Interfaces, ","),
		},
	})
	if err != nil {
		return erk.WrapWith(ErrMockGenFailed, err, erk.Params{
			"packageDescription": pkg.String(),
		})
	}

	result += createNEWMethods(pkg.Interfaces)

	mockFilePath := mockDestination.fullPath()
	mockDirPath := filepath.Dir(mockFilePath)

	if err := g.FSWrite.MkdirAll(mockDirPath, 0775); err != nil {
		return erk.WrapWith(ErrUnableToCreateDir, err, erk.Params{
			"path": mockDirPath,
		})
	}

	if err := g.FSWrite.WriteFile(mockFilePath, result, 0664); err != nil {
		return erk.WrapWith(ErrUnableToCreateFile, err, erk.Params{
			"path": mockFilePath,
		})
	}

	return nil
}

func validateConfig(config *ensurefile.Config) error {
	if config.Mocks == nil {
		return ErrMissingMockConfig
	}

	if config.Mocks.PrimaryDestination == "" {
		config.Mocks.PrimaryDestination = defaultPrimaryDestination
	}

	if config.Mocks.InternalDestination == "" {
		config.Mocks.InternalDestination = defaultInternalDestination
	}

	packages := config.Mocks.Packages
	if len(packages) < 1 {
		return ErrMissingPackages
	}

	// Ensure no duplicate package paths, since the last one would overwrite the first
	packagePaths := map[string]bool{}
	for _, pkg := range packages {
		if _, ok := packagePaths[pkg.Path]; ok {
			return erk.WithParams(ErrDuplicatePackagePath, erk.Params{
				"packagePath": pkg.Path,
			})
		}

		packagePaths[pkg.Path] = true
	}

	return nil
}

func createNEWMethods(interfaces []string) string {
	str := ""

	for _, iface := range interfaces {
		str += fmt.Sprintf(
			"\n// NEW creates a Mock%s.\n"+
				"func (*Mock%s) NEW(ctrl *gomock.Controller) *Mock%s {\n"+
				"\treturn NewMock%s(ctrl)\n"+
				"}\n",
			iface, iface, iface, iface,
		)
	}

	return str
}
