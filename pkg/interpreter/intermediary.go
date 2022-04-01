package interpreter

import (
	"fmt"
	"go/types"
	"strconv"
	"strings"

	"github.com/gabizou/middleware-generator/pkg/errors"

	"github.com/dave/jennifer/jen"
)

type InterpretedVariable interface {
	DebugString() string
	Stringer(name string) jen.Code
	Name() string
	UnderlyingType() types.Type
	AsFunctionParam(name string) jen.Code
	AsReturnType() jen.Code
	AssignImports(file *jen.File)
}

func DeriveInterface(iface *types.Interface) []DeclaredFunction {
	parameter := deriveParameter(0, iface)
	derivedInterface, ok := parameter.(*interfaceLiteral)
	if !ok {
		return nil
	}
	return derivedInterface.functions
}

const _recursiveTypeResolutionLimit = 10

func deriveParameter(attempts int32, variable types.Type) InterpretedVariable {
	if attempts > _recursiveTypeResolutionLimit {
		panic("got too complicated, don't make 5 nested types")
	}
	indent := strings.Repeat(" ", int(attempts))
	switch kind := variable.(type) {
	case *types.Pointer:
		p := &pointerLiteral{inner: deriveParameter(attempts+1, kind.Elem())}
		fmt.Printf("%s%s\n", indent, p.DebugString())
		return p
	case *types.Basic:
		p := &primitive{goType: kind}
		fmt.Printf("%s%s\n", indent, p.DebugString())
		return p
	case *types.Named:
		n := &namedLiteral{named: kind}
		fmt.Printf("%s%s\n", indent, n.DebugString())
		return n
	case *types.Slice:
		s := &sliceLiteral{inner: deriveParameter(attempts+1, kind.Elem())}
		fmt.Printf("%s%s\n", indent, s.DebugString())
		return s
	case *types.Signature:
		f := &functionLiteral{sig: kind}
		params := make([]NamedVariable, kind.Params().Len())
		paramNames := make([]string, kind.Params().Len())
		for p := 0; p < kind.Params().Len(); p++ {
			param := kind.Params().At(p)
			derivedParameter := deriveParameter(attempts+1, param.Type())
			paramNames[p] = param.Name()
			params[p] = &named{
				variable: param,
				name:     param.Name(),
				inner:    derivedParameter,
			}
		}
		returns := make([]NamedVariable, kind.Results().Len())
		for p := 0; p < kind.Results().Len(); p++ {
			ret := kind.Results().At(p)
			derivedReturn := deriveParameter(attempts+1, ret.Type())
			returns[p] = &named{
				variable: ret,
				name:     ret.Name(),
				inner:    derivedReturn,
			}
		}
		f.params = params
		f.paramNames = paramNames
		f.returns = returns
		fmt.Printf("%s%s\n", indent, f.DebugString())
		return f
	case *types.Map:
		k := deriveParameter(attempts+1, kind.Key())
		v := deriveParameter(attempts+1, kind.Elem())
		m := &mapLiteral{kind: kind, key: k, val: v}
		fmt.Printf("%s%s\n", indent, m.DebugString())
		return m
	case *types.Interface:
		fns := make([]DeclaredFunction, kind.NumExplicitMethods())
		i := &interfaceLiteral{
			iface:     kind,
			functions: fns,
		}
		for fn := 0; fn < kind.NumExplicitMethods(); fn++ {
			m := kind.Method(fn)
			sig, ok := m.Type().(*types.Signature)
			if !ok {
				panic(errors.BadMethodSignatureTypeErr{Func: m})
			}
			dc := &declaredFunc{m: m, sig: sig}
			params := sig.Params()
			derivedParams := make([]NamedVariable, params.Len())
			for p := 0; p < params.Len(); p++ {
				param := params.At(p)
				derivedParameter := deriveParameter(attempts+1, param.Type())

				name := param.Name()
				if name == "" {
					suffix := ""
					if p != 0 {
						suffix = strconv.Itoa(p)
					}
					name = fmt.Sprintf("%s%s", derivedParameter.Name(), suffix)
				}
				derivedParams[p] = &named{
					variable: param,
					name:     name,
					inner:    derivedParameter,
				}
			}
			dc.params = derivedParams

			res := sig.Results()
			derivedRes := make([]NamedVariable, res.Len())
			for r := 0; r < res.Len(); r++ {
				res := res.At(r)
				derivedResult := deriveParameter(attempts+1, res.Type())

				name := res.Name()
				wasNamed := true
				if name == "" {
					wasNamed = false
					suffix := ""
					if r != 0 {
						suffix = strconv.Itoa(r)
					}
					name = fmt.Sprintf("%s%s", derivedResult.Name(), suffix)
				}
				derivedRes[r] = &named{
					variable:   res,
					name:       name,
					trulyNamed: wasNamed,
					inner:      derivedResult,
				}
			}
			dc.returns = derivedRes
			fns[fn] = dc
		}
		fmt.Printf("%s%s\n", indent, i.DebugString())
		return i
	case *types.Struct:
		fields := make([]NamedVariable, kind.NumFields())
		s := &structLiteral{
			st:     kind,
			fields: fields,
		}
		for f := 0; f < kind.NumFields(); f++ {
			field := kind.Field(f)
			derivedField := deriveParameter(attempts+1, field.Type())
			fields[f] = &named{
				variable:   field,
				name:       field.Name(),
				inner:      derivedField,
				trulyNamed: true,
			}
		}
		fmt.Printf("%s%s\n", indent, s.DebugString())
		return s
	}
	return nil
}

type primitive struct {
	goType *types.Basic
}

type namedLiteral struct {
	named *types.Named
}

type functionLiteral struct {
	sig        *types.Signature
	params     []NamedVariable
	returns    []NamedVariable
	paramNames []string
}

func (f *functionLiteral) appendFunction(sig *jen.Statement) jen.Code {
	gennedParams := make([]jen.Code, len(f.params))
	for p := 0; p < len(f.params); p++ {
		param := f.params[p]
		starter := jen.Id(f.paramNames[p])
		funcParam := starter.Add(param.AsReturnType())
		gennedParams = append(gennedParams, funcParam)
	}
	sig = sig.Params(gennedParams...)
	gennedReturns := make([]jen.Code, len(f.returns))
	for i, returned := range f.returns {
		gennedReturns[i] = returned.AsReturnType()
	}
	if len(gennedReturns) == 0 {
		return sig
	} else if len(gennedReturns) == 1 {
		sig.Add(gennedReturns[0])
	} else {
		sig.Params(gennedReturns...)
	}
	return sig
}

type interfaceLiteral struct {
	iface     *types.Interface
	functions []DeclaredFunction
}

type declaredFunc struct {
	m       *types.Func
	sig     *types.Signature
	params  []NamedVariable
	returns []NamedVariable
}

type structLiteral struct {
	st     *types.Struct
	fields []NamedVariable
}

type sliceLiteral struct {
	inner InterpretedVariable
}

type pointerLiteral struct {
	inner InterpretedVariable
}

type mapLiteral struct {
	kind *types.Map
	key  InterpretedVariable
	val  InterpretedVariable
}

type named struct {
	variable   *types.Var
	name       string
	inner      InterpretedVariable
	trulyNamed bool
}
