package interpreter

import (
	"fmt"

	"github.com/dave/jennifer/jen"
)

type NamedVariable interface {
	InterpretedVariable
	NamedParameter() jen.Code
}

func (d *declaredFunc) Name() string {
	return d.m.Name()
}

func (f *functionLiteral) Name() string {
	return "func"
}

func (i *interfaceLiteral) Name() string {
	return "interface"
}

func (m *mapLiteral) Name() string {
	return "map"
}

func (n *namedLiteral) Name() string {
	return n.named.Obj().Name()
}

func (p *pointerLiteral) Name() string {
	return fmt.Sprintf("%s%s", "ptr", p.inner.Name())
}

func (p *primitive) Name() string {
	return p.goType.Name()
}

func (s *sliceLiteral) Name() string {
	return fmt.Sprintf("%s%s", "slice", s.inner.Name())
}

func (s *structLiteral) Name() string {
	return "strct"
}

func (n *named) Name() string {
	return n.name
}

func (n *named) NamedParameter() jen.Code {
	return jen.Id(n.name)
}
