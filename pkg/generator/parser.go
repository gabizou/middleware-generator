package generator

import (
	"fmt"
	"go/types"
	"os"
	"strings"

	"github.com/gabizou/middleware-generator/pkg/errors"
	"github.com/gabizou/middleware-generator/pkg/interpreter"

	"golang.org/x/tools/go/packages"
)

type ServiceModel struct {
	TypeName        string
	Middleware      string
	StructPrefix    string
	Interface       []interpreter.DeclaredFunction
	InputParameters []MiddlewareParameter
	StructPtr       string
	ServicePtr      string
}

type File struct {
	Directory  string
	TypeName   string // Name of the constant type.
	Package    *packages.Package
	Middleware string
	Customizer string
}

const (
	packageLoadingMode = packages.NeedName |
		packages.NeedTypes |
		packages.NeedDeps |
		packages.NeedSyntax |
		packages.NeedSyntax |
		packages.NeedTypesInfo
)

func Interpret(dir string, args []string) *Generator {
	targetFile := &File{
		Directory:  dir,
		TypeName:   args[1],
		Middleware: args[2],
		Customizer: args[3],
	}
	interpretedService, err := parseForService(targetFile)
	if err != nil {
		panic(err)
	}
	g := Generator{}
	g.parsePackage(targetFile)
	g.AddFileHeader(strings.Join(args[1:], " "))
	g.SetupCustomizer(targetFile.Customizer)

	// Run generate for each type.
	g.AddModel(interpretedService)
	return &g
}

func parseForService(file *File) (*ServiceModel, error) {
	// 2. Inspect package and use type checker to infer imported types
	file.Package = loadPackage(file.Directory)

	// 3. Lookup the given source type name in the package declarations
	obj := file.Package.Types.Scope().Lookup(file.TypeName)
	if obj == nil {
		failErr(fmt.Errorf("%s not found in declared types of %s",
			file.TypeName, file.Package))
	}

	// 4. We check if it is a declared type
	if _, ok := obj.(*types.TypeName); !ok {
		return nil, errors.UndeclaredTypeErr{Obj: obj}
	}
	// 5. We expect the underlying type to be a struct
	structType, ok := obj.Type().Underlying().(*types.Interface)
	if !ok {
		return nil, errors.NotAnInterfaceErr{Obj: obj}
	}

	// 6. Now we can iterate through fields and access tags
	iface := interpreter.DeriveInterface(structType)
	sm := &ServiceModel{Interface: iface, TypeName: file.TypeName, Middleware: file.Middleware}
	return sm, nil
}

func loadPackage(path string) *packages.Package {
	cfg := &packages.Config{
		Mode: packageLoadingMode,
	}
	pkgs, err := packages.Load(cfg, path)
	if err != nil {
		failErr(err)
	}
	if len(pkgs) != 1 {
		failErr(fmt.Errorf("error: %d packages found", len(pkgs)))
	}

	return pkgs[0]
}

func failErr(err error) {
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
