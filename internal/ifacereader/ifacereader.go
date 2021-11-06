package ifacereader

import (
	"fmt"
	"go/types"
	"sort"

	"golang.org/x/tools/go/packages"
)

type InterfaceReaderIface interface {
	ReadPackages(pkgDetails []*PackageDetails) ([]*Package, error)
}

type InterfaceReader struct{}

var _ InterfaceReaderIface = &InterfaceReader{}

type PackageDetails struct {
	Path       string
	Interfaces []string
}

type Package struct {
	Path       string
	Imports    []string
	Interfaces []*Interface
}

type Interface struct {
	Name    string
	Methods []*Method
}

type Method struct {
	Name    string
	Inputs  []*Tuple
	Outputs []*Tuple
}

type Tuple struct {
	VariableName string
	PackagePaths []string
	Type         string
}

func (r *InterfaceReader) ReadPackages(pkgDetails []*PackageDetails) ([]*Package, error) {
	pkgDetailsByPath := make(map[string]*PackageDetails, len(pkgDetails))
	pkgPaths := make([]string, 0, len(pkgDetails))
	for _, pkgDetail := range pkgDetails {
		if _, ok := pkgDetailsByPath[pkgDetail.Path]; ok {
			return nil, fmt.Errorf("Duplicate entry for path: %s", pkgDetail)
		}

		pkgDetailsByPath[pkgDetail.Path] = pkgDetail
		pkgPaths = append(pkgPaths, pkgDetail.Path)
	}

	cfg := &packages.Config{
		Mode: packages.NeedName | packages.NeedFiles | packages.NeedImports | packages.NeedTypes,
	}

	rawPkgs, err := packages.Load(cfg, pkgPaths...)
	if err != nil {
		return nil, err
	}

	pkgs := make([]*Package, 0, len(rawPkgs))
	for _, pkg := range rawPkgs {
		pkgDetail, ok := pkgDetailsByPath[pkg.PkgPath]
		if !ok {
			return nil, fmt.Errorf("Could not find package details with path: %s", pkg.PkgPath)
		}

		builtPkg, err := buildPackage(pkgDetail, pkg)
		if err != nil {
			return nil, err
		}

		pkgs = append(pkgs, builtPkg)
		delete(pkgDetailsByPath, pkg.PkgPath)
	}

	if len(pkgDetailsByPath) != 0 {
		return nil, fmt.Errorf("Leftover package details: %v", pkgDetailsByPath)
	}

	return pkgs, nil
}

func buildPackage(pkgDetail *PackageDetails, pkg *packages.Package) (*Package, error) {
	ifaces := make([]*Interface, 0, len(pkgDetail.Interfaces))
	importMap := make(map[string]bool)

	for _, ifaceName := range pkgDetail.Interfaces {
		rawIface := pkg.Types.Scope().Lookup(ifaceName)
		if rawIface == nil {
			return nil, fmt.Errorf("Interface %s not found in %s", ifaceName, pkgDetail.Path)
		}

		iface, ok := rawIface.Type().Underlying().(*types.Interface)
		if !ok {
			return nil, fmt.Errorf("Not an interface: %s", rawIface.String())
		}

		builtIface, err := buildIface(ifaceName, iface)
		if err != nil {
			return nil, err
		}

		for _, rawImport := range builtIface.rawImports() {
			importMap[rawImport] = true
		}

		ifaces = append(ifaces, builtIface)
	}

	imports := make([]string, 0, len(importMap))
	for importPath := range importMap {
		imports = append(imports, importPath)
	}
	sort.Strings(imports)

	return &Package{
		Path:       pkgDetail.Path,
		Imports:    imports,
		Interfaces: ifaces,
	}, nil
}

func buildIface(ifaceName string, iface *types.Interface) (*Interface, error) {
	methods := make([]*Method, 0, iface.NumMethods())

	for i := 0; i < iface.NumMethods(); i++ {
		builtMethod, err := buildMethod(iface.Method(i))
		if err != nil {
			return nil, err
		}

		methods = append(methods, builtMethod)
	}

	return &Interface{
		Name:    ifaceName,
		Methods: methods,
	}, nil
}

func buildMethod(method *types.Func) (*Method, error) {
	signature := method.Type().Underlying().(*types.Signature)

	inputs := make([]*Tuple, 0, signature.Params().Len())
	for i := 0; i < signature.Params().Len(); i++ {
		param := signature.Params().At(i)

		builtInput := buildTuple(param.Name(), param.Type())
		inputs = append(inputs, builtInput)
	}

	outputs := make([]*Tuple, 0, signature.Results().Len())
	for i := 0; i < signature.Results().Len(); i++ {
		result := signature.Results().At(i)

		builtOutput := buildTuple(result.Name(), result.Type())
		outputs = append(outputs, builtOutput)
	}

	return &Method{
		Name:    method.Name(),
		Inputs:  inputs,
		Outputs: outputs,
	}, nil
}

func buildTuple(variableName string, rawType types.Type) *Tuple {
	pkgPaths, err := extractPackagePaths(rawType)
	if err != nil {
		panic(err)
	}

	tuple := &Tuple{
		VariableName: variableName,
		Type: types.TypeString(rawType, func(p *types.Package) string {
			return p.Name()
		}),

		PackagePaths: pkgPaths,
	}

	// fmt.Printf("variableName '%s' :: type '%s' :: rawType: %T :: %T\n", tuple.VariableName, tuple.Type, rawType, rawType.Underlying())
	// namedType, ok := rawType.(*types.Named)
	// if ok {
	// 	obj := namedType.Obj()

	// 	if obj.Pkg() != nil {
	// 		tuple.PackagePath = obj.Pkg().Path()
	// 	}
	// } else {
	// 	// namedType, ok := rawType.(*types.Interface)
	// 	// if ok {
	// 	// 	fmt.Println("inface", namedType.Underlying())
	// 	// }
	// }

	return tuple
}

func extractPackagePaths(rawType types.Type) ([]string, error) {
	// fmt.Printf("extractPackagePaths %T: %s\n", rawType, rawType.String())

	switch t := rawType.(type) {
	case *types.Named:
		obj := t.Obj()
		if obj.Pkg() != nil {
			return []string{obj.Pkg().Path()}, nil
		}
		return nil, nil

	case *types.Basic:
		return nil, nil

	case *types.Interface:
		// TODO: Recurse since this is inline
		return nil, nil

	case *types.Struct:
		// TODO: Recurse since this is inline
		return nil, nil

	case *types.Slice:
		return extractPackagePaths(t.Elem())

	case *types.Array:
		return extractPackagePaths(t.Elem())

	case *types.Pointer:
		return extractPackagePaths(t.Elem())

	case *types.Chan:
		return extractPackagePaths(t.Elem())

	case *types.Map:
		keyPaths, err := extractPackagePaths(t.Key())
		if err != nil {
			return nil, err
		}

		elemPaths, err := extractPackagePaths(t.Elem())
		if err != nil {
			return nil, err
		}

		return append(keyPaths, elemPaths...), nil

	case *types.Signature:
		paths := make([]string, 0, t.Params().Len()+t.Results().Len())
		for i := 0; i < t.Params().Len(); i++ {
			paramPaths, err := extractPackagePaths(t.Params().At(i).Type())
			if err != nil {
				return nil, err
			}

			paths = append(paths, paramPaths...)
		}

		for i := 0; i < t.Results().Len(); i++ {
			paramPaths, err := extractPackagePaths(t.Results().At(i).Type())
			if err != nil {
				return nil, err
			}

			paths = append(paths, paramPaths...)
		}

		return paths, nil

	default:
		return nil, fmt.Errorf("not matched %T: %s", rawType, rawType.String())
	}
}

func (iface *Interface) rawImports() []string {
	imports := []string{}
	for _, method := range iface.Methods {
		for _, input := range method.Inputs {
			imports = append(imports, input.PackagePaths...)
		}

		for _, output := range method.Outputs {
			imports = append(imports, output.PackagePaths...)
		}
	}

	return imports
}
