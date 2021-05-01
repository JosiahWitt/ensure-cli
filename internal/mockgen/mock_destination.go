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
	mockPackageName := dest.mockPackageName()
	destPkgFile := filepath.Join(filepath.Dir(dest.rawPackagePath), mockPackageName, mockPackageName+".go")

	return filepath.Join(dest.PWD, dest.MockDir, destPkgFile)
}

func (dest *mockDestination) mockPackageName() string {
	originalPackageName := filepath.Base(dest.rawPackagePath)
	mockPackageName := "mock_" + originalPackageName
	return mockPackageName
}

func (dests mockDestinations) uniquePWDs() []string {
	uniquePWDMap := map[string]bool{}
	for _, dest := range dests {
		uniquePWDMap[dest.PWD] = true
	}

	uniquePWDs := []string{}
	for uniquePWD := range uniquePWDMap {
		uniquePWDs = append(uniquePWDs, uniquePWD)
	}

	return uniquePWDs
}

func (dests mockDestinations) byFullMockDir() map[string]mockDestinations {
	byMockDir := map[string]mockDestinations{}
	for _, dest := range dests {
		key := filepath.Join(dest.PWD, dest.MockDir)
		byMockDir[key] = append(byMockDir[key], dest)
	}

	return byMockDir
}

func (dests mockDestinations) byPackagePath() map[string]*mockDestination {
	byPackagePath := map[string]*mockDestination{}
	for _, dest := range dests {
		byPackagePath[dest.Package.Path] = dest
	}

	return byPackagePath
}

func (dests mockDestinations) hasFullPathPrefix(prefix string) bool {
	for _, dest := range dests {
		if strings.HasPrefix(dest.fullPath(), prefix) {
			return true
		}
	}

	return false
}
