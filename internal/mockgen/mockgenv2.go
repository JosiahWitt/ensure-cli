package mockgen

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"path/filepath"

	"github.com/JosiahWitt/ensure-cli/internal/ensurefile"
	"github.com/JosiahWitt/ensure-cli/internal/fswrite"
	"github.com/JosiahWitt/ensure-cli/internal/ifacereader"
	"github.com/JosiahWitt/erk"
	"github.com/dave/jennifer/jen"
)

type MockGenV2 struct {
	FSWrite fswrite.FSWriteIface
	Logger  *log.Logger

	IfaceReader ifacereader.InterfaceReaderIface
}

var _ MockGenerator = &MockGenV2{}

// GenerateMocks for the provided configuration.
func (g *MockGenV2) GenerateMocks(ctx context.Context, config *ensurefile.Config) error {
	if err := validateConfig(config); err != nil {
		return err
	}

	g.Logger.Println("Loading packages...")

	mockDestinations, err := computeMockDestinations(config)
	if err != nil {
		return err
	}

	mockDestsByPkg := mockDestinations.byPackagePath()

	pkgDetails := make([]*ifacereader.PackageDetails, 0, len(config.Mocks.Packages))
	for _, pkgDetail := range config.Mocks.Packages {
		pkgDetails = append(pkgDetails, &ifacereader.PackageDetails{
			Path:       pkgDetail.Path,
			Interfaces: pkgDetail.Interfaces,
		})
	}

	pkgs, err := g.IfaceReader.ReadPackages(pkgDetails)
	if err != nil {
		return err
	}

	g.Logger.Println("Generating mocks:")
	for _, pkg := range pkgs {
		mockDest, ok := mockDestsByPkg[pkg.Path]
		if !ok {
			fmt.Println(mockDestsByPkg)
			return errors.New("oops " + pkg.Path)
		}

		g.Logger.Printf(" - Generating: %s\n", mockDest.Package.String())

		mockFile, err := generatePackage(pkg, mockDest)
		if err != nil {
			return err
		}

		mockFilePath := mockDest.fullPath()
		mockDirPath := filepath.Dir(mockFilePath)

		if err := g.FSWrite.MkdirAll(mockDirPath, 0775); err != nil {
			return erk.WrapWith(ErrUnableToCreateDir, err, erk.Params{
				"path": mockDirPath,
			})
		}

		if err := g.FSWrite.WriteFile(mockFilePath, mockFile, 0664); err != nil {
			return err
		}
	}

	return nil
}

func generatePackage(pkg *ifacereader.Package, mockDest *mockDestination) (string, error) {
	f := jen.NewFile(mockDest.mockPackageName())

	importDecl := jen.Line()
	for _, importPath := range pkg.Imports {
		importDecl.Lit(importPath).Line()
	}

	f.Id("import").Parens(importDecl)

	for _, iface := range pkg.Interfaces {
		generateInterface(f, pkg.Path, iface)
	}

	buf := &bytes.Buffer{}
	if err := f.Render(buf); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func generateInterface(f *jen.File, packagePath string, iface *ifacereader.Interface) {
	mockStruct := "Mock" + iface.Name
	recorderStruct := mockStruct + "MockRecorder"

	// Mock struct
	f.Type().Id(mockStruct).Struct(
		jen.Id("ctrl").Op("*").Qual("github.com/golang/mock/gomock", "Controller"),
		jen.Id("recorder").Op("*").Id(recorderStruct),
	).Line()

	// Recorder struct
	f.Type().Id(recorderStruct).Struct(
		jen.Id("mock").Op("*").Id(mockStruct),
	).Line()

	// New mock function
	f.Func().Id("New"+mockStruct).Params(
		jen.Id("ctrl").Op("*").Qual("github.com/golang/mock/gomock", "Controller"),
	).Params(jen.Op("*").Id(mockStruct)).Block(
		jen.Id("mock").Op(":=").Op("&").Id(mockStruct).Values(jen.Dict{
			jen.Id("ctrl"): jen.Id("ctrl"),
		}),

		jen.Id("mock.recorder").Op("=").Op("&").Id(recorderStruct).Values(jen.Dict{
			jen.Id("mock"): jen.Id("mock"),
		}),

		jen.Return(jen.Id("mock")),
	).Line()

	// EXPECT method
	f.Func().Params(
		jen.Id("m").Op("*").Id(mockStruct),
	).Id("EXPECT").Params().Params(jen.Op("*").Id(recorderStruct)).Block(
		jen.Return(jen.Id("m.recorder")),
	).Line()

	// NEW method
	f.Func().Params(
		jen.Id("m").Op("*").Id(mockStruct),
	).Id("NEW").Params(
		jen.Id("ctrl").Op("*").Qual("github.com/golang/mock/gomock", "Controller"),
	).Params(jen.Op("*").Id(mockStruct)).Block(
		jen.Return(jen.Id("New" + mockStruct).Call(jen.Id("ctrl"))),
	).Line()

	// methodDecls := make([]jen.Code, 0, len(iface.Methods))
	for _, meth := range iface.Methods {
		realInputs := buildTuples(packagePath, meth.Inputs, "")
		realOutputs := buildTuples(packagePath, meth.Outputs, "")

		implCallInputs := []jen.Code{jen.Id("m"), jen.Lit(meth.Name)}
		for _, input := range meth.Inputs {
			implCallInputs = append(implCallInputs, jen.Id("_"+input.VariableName))
		}

		returnCastings := make([]jen.Code, 0, len(meth.Outputs))
		for i, output := range meth.Outputs {
			returnCastings = append(returnCastings,
				jen.List(
					jen.Id(fmt.Sprintf("ret%d", i)),
					jen.Id("_"),
				).
					Op(":=").
					Id("ret").Index(jen.Lit(i)).
					Assert(jen.Id(output.Type)),
			)
		}

		if len(meth.Outputs) == 0 {
			returnCastings = append(returnCastings,
				jen.Var().Id("_").Op("=").Id("ret"),
			)
		}

		returnVars := make([]jen.Code, 0, len(meth.Outputs))
		for i := range meth.Outputs {
			returnVars = append(returnVars,
				jen.Id(fmt.Sprintf("ret%d", i)),
			)
		}

		// Implementation method
		f.Func().Params(
			jen.Id("m").Op("*").Id(mockStruct),
		).Id(meth.Name).Params(realInputs...).Params(realOutputs...).Block(
			append(
				append(
					[]jen.Code{
						jen.Id("m.ctrl.T.Helper").Call(),
						jen.Id("ret").Op(":=").Id("m.ctrl.Call").Call(implCallInputs...),
					},

					returnCastings...,
				),

				jen.Return(returnVars...),
			)...,
		).Line()

		recorderCallInputs := []jen.Code{
			jen.Id("mr.mock"),
			jen.Lit(meth.Name),
			// reflect.TypeOf((*MockMockGenerator)(nil).TidyMocks),
			jen.Qual("reflect", "TypeOf").Call(jen.Parens(jen.Op("*").Id(mockStruct)).Parens(jen.Nil()).Dot(meth.Name)),
		}
		for _, input := range meth.Inputs {
			recorderCallInputs = append(recorderCallInputs, jen.Id("_"+input.VariableName))
		}

		recorderInputs := buildTuples(packagePath, meth.Inputs, "interface{}")

		// Expectation method
		f.Func().Params(
			jen.Id("mr").Op("*").Id(recorderStruct),
		).Id(meth.Name).Params(recorderInputs...).Params(
			jen.Op("*").Qual("github.com/golang/mock/gomock", "Call"),
		).Block(
			jen.Id("mr.mock.ctrl.T.Helper").Call(),
			jen.Return(jen.Id("mr.mock.ctrl.RecordCallWithMethodType").Call(recorderCallInputs...)),
		).Line()

		// methodDecl := jen.Id(meth.Name).Params(inputs...).Params(outputs...)

		// methodDecls = append(methodDecls, methodDecl)
	}

	// f.Type().Id(iface.Name).Interface(methodDecls...)
	// f.Line()
}

func buildTuples(packagePath string, tuple []*ifacereader.Tuple, overrideType string) []jen.Code {
	params := make([]jen.Code, 0, len(tuple))

	for _, input := range tuple {
		param := jen.Null()
		if input.VariableName != "" {
			param = param.Id("_" + input.VariableName)
		}
		// if len(input.PackagePaths) == 0 {

		if overrideType != "" {
			param = param.Id(overrideType)
		} else {
			param = param.Id(input.Type)
		}
		// } else {
		// 	param = param.Qual(input.PackagePaths[0], input.Type) // TODO: Multiple paths
		// }

		params = append(params, param)
	}

	return params
}
