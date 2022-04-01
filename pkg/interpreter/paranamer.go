package interpreter

import (
	"github.com/dave/jennifer/jen"
)

func (p *primitive) AsFunctionParam(name string) jen.Code {
	return jen.Id(name).Id(p.goType.Name())
}

func (f *functionLiteral) AsFunctionParam(name string) jen.Code {
	sig := jen.Id(name).Func()
	return f.appendFunction(sig)
}

func (i *interfaceLiteral) AsFunctionParam(name string) jen.Code {
	fns := make([]jen.Code, len(i.functions))
	for i2, fn := range i.functions {
		fns[i2] = fn.AsReturnType()
	}
	return jen.Id(name).Interface(fns...)
}

func (d *declaredFunc) AsFunctionParam(name string) jen.Code {
	genedFunction := jen.Id(name).Func()
	params := d.sig.Params()
	genParams := make([]jen.Code, len(d.params))
	for i, param := range d.params {
		genParams[i] = param.AsFunctionParam(params.At(i).Name())
	}
	genedFunction.Params(genParams...)
	genReturns := make([]jen.Code, len(d.returns))
	for i, param := range d.returns {
		genReturns[i] = param.AsReturnType()
	}
	if len(genReturns) == 1 {
		genedFunction.Add(genReturns...)
	} else if len(genReturns) > 1 {
		genedFunction.Params(genReturns...)
	}
	return genedFunction
}

func (s *structLiteral) AsFunctionParam(name string) jen.Code {
	flds := make([]jen.Code, len(s.fields))
	for i, field := range s.fields {
		flds[i] = jen.Id(field.Name()).Add(field.AsReturnType())
	}
	return jen.Id(name).Struct(flds...)
}

func (s *sliceLiteral) AsFunctionParam(name string) jen.Code {
	return jen.Id(name).Index().Add(s.inner.AsReturnType())
}

func (p *pointerLiteral) AsFunctionParam(name string) jen.Code {
	return jen.Id(name).Op("*").Add(p.inner.AsReturnType())
}

func (m *mapLiteral) AsFunctionParam(name string) jen.Code {
	return jen.Id(name).Map(m.key.AsReturnType()).Add(m.val.AsReturnType())
}

func (n *named) AsFunctionParam(_ string) jen.Code {
	return jen.Id(n.name).Add(n.inner.AsReturnType())
}

func (n *namedLiteral) AsFunctionParam(name string) jen.Code {
	code := jen.Id(name)
	obj := n.named.Obj()
	pkg := obj.Pkg()
	if pkg == nil {
		return code.Id(obj.Name())
	}
	if !obj.Exported() {
		return code.Id(obj.Name())
	}

	return code.Qual(pkg.Path(), obj.Name())
}
