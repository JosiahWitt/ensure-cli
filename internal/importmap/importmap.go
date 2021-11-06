package importmap

import "regexp"

var removeSpecialCharacters = regexp.MustCompile("[^a-zA-Z0-9]+")

type ImportMap struct {
	byFullPackagePath map[string]string
	byPackageAlias    map[string]string
}

func (m ImportMap) Lookup(fullPackagePath string) string {
	if packageAlias, ok := m.byFullPackagePath[fullPackagePath]; ok {
		return packageAlias
	}

	packageAlias := m.generatePackageAlias(fullPackagePath)
	m.byFullPackagePath[fullPackagePath] = packageAlias
	m.byPackageAlias[packageAlias] = fullPackagePath
	return packageAlias
}

func (m ImportMap) generatePackageAlias(fullPackagePath string) string {
	return fullPackagePath
}
