package generator

import (
	"fmt"
	"sync"

	"github.com/dave/jennifer/jen"
	"github.com/gabizou/middleware-generator/pkg/interpreter"
)

// Customizer is utilized by Generator to specify the generated output
// of a desired middleware. Each instance should be registered with
// generator.Register.
type Customizer interface {
	// FileNamePrefix gives a prefix for the generated file name.
	FileNamePrefix() string
	// FactorySuffix provides a variant typed suffix for a Factory method
	// generated for the middleware.
	FactorySuffix() string
	// ConfigureModel takes the given ServiceModel and can apply any specifications
	// such as desired inputs to the Factory method,
	ConfigureModel(model *ServiceModel)
	GetRequiredImportNames() map[string]string
	// GenerateFunctionImplementation will be passed in a builder
	// of pre-computed code statements as the function declaration
	// as described by the passed-in DeclaredFunction. It is important
	// to note that no code blocks have been created at this point, see
	// jen.Block for more details.
	GenerateFunctionImplementation(builder *jen.Statement, service *ServiceModel, method interpreter.DeclaredFunction) jen.Code
}

type MiddlewareParameter struct {
	VariableName string
	TypeName     string
	TypePath     string
	FieldName    string
}

var (
	customizersMu sync.RWMutex
	customizers   = make(map[string]Customizer)
)

// Register makes a Customizer available by name. If a Customizer by
// the same name is registered, an error is thrown.
func Register(name string, customizer Customizer) {
	customizersMu.Lock()
	defer customizersMu.Unlock()
	if customizer == nil {
		panic("generator: Register customizer is nil")
	}
	if _, dup := customizers[name]; dup {
		panic("generator: Register called twice for customizer " + name)
	}
	customizers[name] = customizer
}

func (g *Generator) SetupCustomizer(customizer string) {
	customizersMu.RLock()
	defer customizersMu.RUnlock()
	g.customizer = customizers[customizer]
	if g.customizer == nil {
		panic(fmt.Errorf("generator: No customizer found by name %s", customizer))
	}
}
