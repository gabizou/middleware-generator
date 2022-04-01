package interpreter

import (
	"fmt"
	"go/types"

	"github.com/dave/jennifer/jen"
)

func (p *primitive) DebugString() string {
	return p.goType.String()
}

func (p *primitive) Stringer(name string) jen.Code {
	basic := p.goType.Kind()
	switch basic {
	case types.Invalid:
		return jen.Lit(`"undefined"`)
	case types.Bool, types.UntypedBool:
		return jen.Qual("fmt", "Sprintf").Params(jen.Lit(`"%v"`), jen.Lit(name))
	case types.UntypedRune:
		return jen.Qual("fmt", "Sprintf").Params(jen.Lit(`"%c"`), jen.Lit(name))
	case types.Int, types.Int8,
		types.Int16, types.Int32,
		types.Int64, types.Uint,
		types.Uint8, types.Uint16,
		types.Uint32, types.Uint64,
		types.Float32, types.Float64,
		types.Complex64, types.Complex128,
		types.UntypedInt, types.UntypedFloat,
		types.UntypedComplex, types.Uintptr, types.UnsafePointer:
		return jen.Qual("fmt", "Sprintf").Params(jen.Lit(`"%v"`), jen.Lit(name))
	case types.String, types.UntypedString, types.UntypedNil:
		return jen.Lit(name)
	}
	return jen.Lit(`"undefined"`)
}

func (f *functionLiteral) DebugString() string {
	return f.sig.String()
}

func (f *functionLiteral) Stringer(name string) jen.Code {
	return jen.Qual("fmt", "Sprintf").Params(jen.Lit(`"%v"`), jen.Lit(name))
}

func (i *interfaceLiteral) DebugString() string {
	return fmt.Sprintf("Interface type: %s", i.iface.String())
}

func (i *interfaceLiteral) Stringer(name string) jen.Code {
	return jen.Qual("fmt", "Sprintf").Params(jen.Lit(`"%v"`), jen.Lit(name))
}

func (n *namedLiteral) DebugString() string {
	return n.named.String()
}

func (n *namedLiteral) Stringer(name string) jen.Code {
	return jen.Qual("fmt", "Sprintf").Params(jen.Lit(`"%v"`), jen.Lit(name))
}

func (d *declaredFunc) DebugString() string {
	return fmt.Sprintf("Declared Func: %s", d.sig.String())
}

func (d *declaredFunc) Stringer(name string) jen.Code {
	return jen.Qual("fmt", "Sprintf").Params(jen.Lit(`"%v"`), jen.Lit(name))
}

func (s *structLiteral) DebugString() string {
	return fmt.Sprintf("Struct type: %s", s.st.String())
}

func (s *structLiteral) Stringer(name string) jen.Code {
	return jen.Qual("fmt", "Sprintf").Params(jen.Lit(`"%v"`), jen.Lit(name))
}

func (s *sliceLiteral) DebugString() string {
	return fmt.Sprintf("Slice type: %s\n", s.inner.DebugString())
}

func (s *sliceLiteral) Stringer(name string) jen.Code {
	return jen.Qual("fmt", "Sprintf").Params(jen.Lit(`"%v"`), jen.Lit(name))
}

func (p *pointerLiteral) DebugString() string {
	return fmt.Sprintf("Pointer type: %#v\n", p.inner.DebugString())
}

func (p *pointerLiteral) Stringer(name string) jen.Code {
	return jen.Qual("fmt", "Sprintf").Params(jen.Lit(`"%v"`), jen.Lit(name))
}

func (m *mapLiteral) DebugString() string {
	return fmt.Sprintf("Map type: %s(%s[%s])", m.kind.String(), m.key.DebugString(), m.val.DebugString())
}

func (m *mapLiteral) Stringer(name string) jen.Code {
	return jen.Qual("fmt", "Sprintf").Params(jen.Lit(`"%v"`), jen.Lit(name))
}

func (n *named) DebugString() string {
	return fmt.Sprintf("Named %s(%s)", n.name, n.inner.DebugString())
}

func (n *named) Stringer(name string) jen.Code {
	return jen.Qual("fmt", "Sprintf").Params(jen.Lit(`"%v"`), jen.Lit(name))
}
