package interpreter

import (
	"github.com/dave/jennifer/jen"
)

func (p *primitive) AsReturnType() jen.Code {
	return jen.Id(p.goType.Name())
}

func (n *namedLiteral) AsReturnType() jen.Code {
	obj := n.named.Obj()
	pkg := obj.Pkg()
	if pkg == nil {
		return jen.Id(obj.Name())
	}
	if !obj.Exported() {
		return jen.Id(obj.Name())
	}
	return jen.Qual(pkg.Path(), obj.Name())
}

func (f *functionLiteral) AsReturnType() jen.Code {
	return f.appendFunction(jen.Func())
}

func (i *interfaceLiteral) AsReturnType() jen.Code {
	fns := make([]jen.Code, len(i.functions))
	for i2, fn := range i.functions {
		fns[i2] = fn.AsReturnType()
	}
	return jen.Interface(fns...)
}

func (d *declaredFunc) AsReturnType() jen.Code {
	genedFunction := jen.Id(d.m.Name())
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

func (s *structLiteral) AsReturnType() jen.Code {
	flds := make([]jen.Code, len(s.fields))
	for i, field := range s.fields {
		flds[i] = jen.Id(field.Name()).Add(field.AsReturnType())
	}
	return jen.Struct(flds...)
}

func (s *sliceLiteral) AsReturnType() jen.Code {
	return jen.Index().Add(s.inner.AsReturnType())
}

func (p *pointerLiteral) AsReturnType() jen.Code {
	return jen.Op("*").Add(p.inner.AsReturnType())
}

func (n *named) AsReturnType() jen.Code {
	return n.inner.AsReturnType()
}

func (m *mapLiteral) AsReturnType() jen.Code {
	return jen.Map(m.key.AsReturnType()).Add(m.val.AsReturnType())
}
