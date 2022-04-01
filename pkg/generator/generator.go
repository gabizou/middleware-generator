package generator

import (
	"fmt"
	"log"
	"strings"

	"github.com/gabizou/middleware-generator/pkg/interpreter"

	"github.com/dave/jennifer/jen"

	"golang.org/x/tools/go/packages"
)

// Generator is an object to take an interpreted interpreter.ServiceModel and
// using a Customizer, generates the middleware output.
type Generator struct {
	f                    *jen.File // The generating file we're working on
	ourType              string
	ourPtr               rune
	svcPtr               string
	service              *ServiceModel
	interpretedFunctions []interpreter.DeclaredFunction
	customizer           Customizer
}

const (
	_packageParsing = packages.NeedTypes |
		packages.NeedTypesSizes |
		packages.NeedImports |
		packages.NeedName |
		packages.NeedFiles |
		packages.NeedCompiledGoFiles |
		packages.NeedTypesInfo |
		packages.NeedSyntax
)

func (g *Generator) parsePackage(file *File) {
	cfg := &packages.Config{
		Mode:       _packageParsing,
		Tests:      false,
		BuildFlags: []string{},
	}

	pkgs, err := packages.Load(cfg, file.Directory)
	if err != nil {
		log.Fatal(err)
	}
	if len(pkgs) != 1 {
		log.Fatalf("error: %d packages found", len(pkgs))
	}
	g.f = jen.NewFile(pkgs[0].Name)
}

func (g *Generator) AddFileHeader(header string) {
	g.f.HeaderComment(fmt.Sprintf("Code generated by \"middleware-generator %s\"; DO NOT EDIT.", header))
}

func (g *Generator) AddModel(model *ServiceModel) {
	pointerName := string(strings.ToLower(model.TypeName)[0])
	// Specify the default import name because we don't want to alias this
	for alias, v := range g.customizer.GetRequiredImportNames() {
		g.f.ImportName(alias, v)
	}
	g.customizer.ConfigureModel(model)

	middlewareTypeName := fmt.Sprintf(model.StructPrefix, string(model.TypeName[0]))
	g.ourType = middlewareTypeName
	g.ourPtr = []rune(strings.ToLower(middlewareTypeName))[0]
	model.StructPtr = string(g.ourPtr)
	model.ServicePtr = pointerName
	g.svcPtr = pointerName
	g.service = model
	g.interpretedFunctions = model.Interface
	g.genFactoryMethod(model)
	g.genStruct(model)
	g.genInterfaceMethods()
}

// genStruct creates the following:
// type logger${shortenedName} struct {
//   l log.Logger
//   ${shortenedName} ${ServiceModel.TypeName}
// }
func (g *Generator) genStruct(model *ServiceModel) *jen.Statement {
	fields := make([]jen.Code, len(model.InputParameters))
	for i, parameter := range model.InputParameters {
		field := jen.Id(parameter.FieldName)
		if parameter.TypePath != "" {
			field.Qual(parameter.TypePath, parameter.TypeName)
		} else {
			field.Id(parameter.TypeName)
		}
		fields[i] = field
	}
	fields = append(fields, jen.Id(g.svcPtr).Id(model.TypeName))

	return g.f.Type().
		Id(g.ourType).
		Struct(
			fields...,
		)
}

// genFactoryMethod creates the following:
//  func New${ServiceModel.TypeName}(${ServiceModel.InputParameters}.Name ${ServiceModel.InputParameters}) ${ServiceModel.Middleware} {
//    return func($shortenedName ${ServiceModel.TypeName}) ${ServiceModel.TypeName} {
//      return &${ServiceModel.${shortenedName}{
//          fields: variables,
//          ${shortenedName}: ${$shortenedName},
//    }
//  }
//}
func (g *Generator) genFactoryMethod(model *ServiceModel) *jen.Statement {
	genParams := make([]jen.Code, len(model.InputParameters))
	for i, parameter := range model.InputParameters {
		genParams[i] = jen.Id(parameter.VariableName).Qual(parameter.TypePath, parameter.TypeName)
	}
	return g.f.Func().
		Id(fmt.Sprintf("New%s%s", model.TypeName, g.customizer.FactorySuffix())).
		Params(genParams...).
		Id(model.Middleware).
		BlockFunc(func(gr *jen.Group) {
			gr.ReturnFunc(func(ig *jen.Group) {
				ig.Func().
					Params(jen.Id(g.svcPtr).Id(model.TypeName)).
					Id(model.TypeName).
					BlockFunc(func(ng *jen.Group) {
						fieldSetters := make(jen.Dict)
						fieldSetters[jen.Id(g.svcPtr)] = jen.Id(g.svcPtr)
						for _, parameter := range model.InputParameters {
							fieldSetters[jen.Id(parameter.FieldName)] = jen.Id(parameter.VariableName)
						}
						ng.Return(jen.Op("&").Id(g.ourType).Values(fieldSetters))
					})
			})
		})
}

func (g *Generator) genInterfaceMethods() {
	for _, method := range g.interpretedFunctions {
		// get the function for naming
		gennedFunc := g.genFunctionDeclaration(method)
		g.f.Add(gennedFunc)
	}
}

func (g *Generator) genFunctionDeclaration(method interpreter.DeclaredFunction) jen.Code {
	genedFunction := jen.Func().
		Params(jen.Id(g.service.StructPtr).Op("*").Id(g.ourType))
	genedFunction = genedFunction.
		Id(method.FunctionName())
	var genParams []jen.Code
	for _, variable := range method.Parameters() {
		parameter := variable.AsFunctionParam(variable.Name())
		genParams = append(genParams, parameter)
	}
	genedFunction.Params(genParams...)

	genedFunction.Add(method.ReturnDefinition())
	g.customizer.GenerateFunctionImplementation(genedFunction, g.service, method)
	return genedFunction
}

func (g *Generator) Print() {
	err := g.f.Save(fmt.Sprintf("%s_%s.go", g.customizer.FileNamePrefix(), strings.ToLower(g.service.TypeName)))
	if err != nil {
		panic(err)
	}
}