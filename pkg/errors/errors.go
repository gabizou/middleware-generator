package errors

import (
	"fmt"
	"go/types"
)

// BadMethodSignatureTypeErr represents a type assertion failure that we want to know
// about where a given Func turns out to be not a types.Signature, which ends up being
// an unknown/unhandled function for us to generate a wrapper around.
type BadMethodSignatureTypeErr struct {
	Func *types.Func
}

func (b BadMethodSignatureTypeErr) Error() string {
	return fmt.Sprintf("method by %s has %T instead of *types.Signature, violation of types.Func", b.Func.Name(), b.Func.Type())
}

type UndeclaredTypeErr struct {
	Obj types.Object
}

func (u UndeclaredTypeErr) Error() string {
	return fmt.Sprintf("%v is not a named type", u.Obj)
}

type NotAnInterfaceErr struct {
	Obj types.Object
}

func (n NotAnInterfaceErr) Error() string {
	return fmt.Sprintf("type %v is not an interface!", n.Obj)
}
