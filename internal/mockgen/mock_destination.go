package mockgen

import (
	"path/filepath"
	"strings"

	"github.com/JosiahWitt/ensure-cli/internal/ensurefile"
	"github.com/JosiahWitt/erk"
	"github.com/JosiahWitt/erk/erg"
)

type ErkMockDestination struct{ erk.DefaultKind }

var ErrInternalPackageOutsideModule = erk.New(ErkMockDestination{},
	"Cannot generate mock of internal package, since package '{{.packagePath}}' is not in the current module '{{.modulePath}}'",
)

type mockDestinations []*mockDestination

type mockDestination struct {
	Package        *ensurefile.Package
	PWD            string
	MockDir        string
	rawPackagePath string
}

func computeMockDestinations(config *ensurefile.Config) (mockDestinations, error) {
	errGroup := erg.NewAs(ErrMultipleGenerationFailures)

	destinations := mockDestinations{}
	for _, pkg := range config.Mocks.Packages {
		dest, err := computeMockDestination(config, pkg)
		if err != nil {
			errGroup = erg.Append(errGroup, err)
			continue
		}

		destinations = append(destinations, dest)
	}

	if erg.Any(errGroup) {
		return nil, errGroup
	}

	return destinations, nil
}

func computeMockDestination(config *ensurefile.Config, pkg *ensurefile.Package) (*mockDestination, error) {
	const internalPart = "internal/"

	// Check if package is internal
	idx := strings.LastIndex(pkg.Path, internalPart)
	if idx < 0 {
		return &mockDestination{
			Package:        pkg,
			PWD:            config.RootPath,
			MockDir:        config.Mocks.PrimaryDestination,
			rawPackagePath: pkg.Path,
		}, nil
	}

	if !strings.HasPrefix(pkg.Path, config.ModulePath) {
		return nil, erk.WithParams(ErrInternalPackageOutsideModule, erk.Params{
			"packagePath": pkg.Path,
			"modulePath":  config.ModulePath,
		})
	}

	// Remove both the module path prefix, and the last internal/... suffix
	pkgPathPrefix := strings.TrimPrefix(pkg.Path[:idx], config.ModulePath)

	// Everything after the last internal/...
	pkgPathSuffix := pkg.Path[idx+len(internalPart):]

	return &mockDestination{
		Package:        pkg,
		PWD:            filepath.Join(config.RootPath, pkgPathPrefix),
		MockDir:        filepath.Join(internalPart, config.Mocks.InternalDestination),
		rawPackagePath: pkgPathSuffix,
	}, nil
}

func (dest *mockDestination) fullPath() string {
	originalPackageName := filepath.Base(dest.rawPackagePath)
	mockPackageName := "mock_" + originalPackageName
	destPkgFile := filepath.Join(filepath.Dir(dest.rawPackagePath), mockPackageName, mockPackageName+".go")

	return filepath.Join(dest.PWD, dest.MockDir, destPkgFile)
}

func (dests mockDestinations) uniquePWDs() []string {
	return dests.uniqueString(func(dest *mockDestination) string {
		return dest.PWD
	})
}

func (dests mockDestinations) uniqueString(fn func(dest *mockDestination) string) []string {
	uniqueMap := map[string]bool{}
	for _, dest := range dests {
		str := fn(dest)
		uniqueMap[str] = true
	}

	items := []string{}
	for item := range uniqueMap {
		items = append(items, item)
	}

	return items
}
