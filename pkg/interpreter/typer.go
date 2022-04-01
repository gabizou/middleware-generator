package interpreter

import (
	"go/types"
)

func (p *primitive) UnderlyingType() types.Type {
	return p.goType
}

func (n *namedLiteral) UnderlyingType() types.Type {
	return n.named
}

func (f *functionLiteral) UnderlyingType() types.Type {
	return f.sig
}

func (i *interfaceLiteral) UnderlyingType() types.Type {
	return i.iface
}

func (d *declaredFunc) UnderlyingType() types.Type {
	return d.sig
}

func (s *structLiteral) UnderlyingType() types.Type {
	return s.st
}

func (s *sliceLiteral) UnderlyingType() types.Type {
	return s.inner.UnderlyingType()
}

func (p *pointerLiteral) UnderlyingType() types.Type {
	return p.inner.UnderlyingType()
}

func (n *named) UnderlyingType() types.Type {
	return n.variable.Type()
}

func (m *mapLiteral) UnderlyingType() types.Type {
	return m.kind
}
