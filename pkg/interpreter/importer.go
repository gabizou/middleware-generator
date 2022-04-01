package interpreter

import (
	"github.com/dave/jennifer/jen"
)

func (p *primitive) AssignImports(_ *jen.File) {}

func (n *namedLiteral) AssignImports(f *jen.File) {
	pkg := n.named.Obj().Pkg()
	if pkg != nil {
		f.ImportName(pkg.Path(), pkg.Name())
	}
}

func (f *functionLiteral) AssignImports(file *jen.File) {
	for _, param := range f.params {
		param.AssignImports(file)
	}
	for _, ret := range f.returns {
		ret.AssignImports(file)
	}
}

func (i *interfaceLiteral) AssignImports(file *jen.File) {
	for _, fn := range i.functions {
		fn.AssignImports(file)
	}
}

func (d *declaredFunc) AssignImports(file *jen.File) {
	for _, p := range d.params {
		p.AssignImports(file)
	}
	for _, r := range d.returns {
		r.AssignImports(file)
	}
}

func (s *structLiteral) AssignImports(file *jen.File) {
	for _, field := range s.fields {
		field.AssignImports(file)
	}
}

func (s *sliceLiteral) AssignImports(f *jen.File) {
	s.inner.AssignImports(f)
}

func (p *pointerLiteral) AssignImports(_ *jen.File) {}

func (m *mapLiteral) AssignImports(file *jen.File) {
	m.key.AssignImports(file)
	m.val.AssignImports(file)
}

func (n *named) AssignImports(file *jen.File) {
	n.inner.AssignImports(file)
}
