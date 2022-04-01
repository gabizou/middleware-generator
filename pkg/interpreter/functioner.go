package interpreter

import (
	"github.com/dave/jennifer/jen"
)

type DeclaredFunction interface {
	InterpretedVariable
	FunctionName() string
	Parameters() []NamedVariable
	Returns() []NamedVariable
	ReturnDefinition() jen.Code
}

func (d *declaredFunc) FunctionName() string {
	return d.m.Name()
}

func (d *declaredFunc) Parameters() []NamedVariable {
	return d.params
}

func (d *declaredFunc) Returns() []NamedVariable {
	return d.returns
}

func (d *declaredFunc) ReturnDefinition() jen.Code {
	builder := jen.Add()
	var genReturnTypes []jen.Code
	paramsOverride := false
	for _, variable := range d.Returns() {
		if variable.Name() != "" {
			paramsOverride = true
		}
		genReturnTypes = append(genReturnTypes, variable.AsReturnType())
	}
	if len(genReturnTypes) == 1 {
		if paramsOverride {
			builder.Params(genReturnTypes...)
		} else {
			builder.Add(genReturnTypes...)
		}
	} else if len(genReturnTypes) > 1 {
		builder.Params(genReturnTypes...)
	}
	return builder
}
