package tracing

import (
	"go/types"

	"github.com/gabizou/middleware-generator/pkg/generator"
	"github.com/gabizou/middleware-generator/pkg/interpreter"

	"github.com/dave/jennifer/jen"
)

func init() { //nolint:gochecknoinits
	generator.Register("tracer", tracer{})
}

type tracer struct {
}

func (t tracer) FileNamePrefix() string {
	return "tracer"
}

func (t tracer) FactorySuffix() string {
	return "Tracer"
}

func (t tracer) ConfigureModel(model *generator.ServiceModel) {
	model.StructPrefix = "tracer%s"
	model.InputParameters = []generator.MiddlewareParameter{
		{
			VariableName: "tracer",
			TypeName:     "Tracer",
			TypePath:     "github.com/openzipkin/zipkin-go",
			FieldName:    "tr",
		},
	}
}

// GenerateFunctionImplementation generates a function
// that depending on the method having a context.Context variable,
// will either forward along, or create a new span from context
func (t tracer) GenerateFunctionImplementation(
	builder *jen.Statement,
	service *generator.ServiceModel,
	method interpreter.DeclaredFunction,
) jen.Code {
	ctxName := ""
	hasContext := false
	for _, variable := range method.Parameters() {
		if typed, ok := variable.UnderlyingType().(*types.Named); ok {
			obj := typed.Obj()
			if obj.Name() == "Context" && obj.Pkg().Path() == "context" {
				hasContext = true
				ctxName = variable.Name()
				break
			}
		}
	}
	lines := make([]jen.Code, 0)
	if hasContext {
		startSpan := jen.List(
			jen.Id("span"),
			jen.Id(ctxName),
		).Op(":=").Id(service.StructPtr).Dot("tr").
			Dot("StartSpanFromContext").
			Call(
				jen.Id(ctxName),
				jen.Lit(method.Name()),
			)
		/* code to generate
		span, $ctxName := ${service.StructPtr}.tr.StartSpanFromContext($ctxName, ${DeclaredFunction.Name})
		*/
		lines = append(lines, startSpan)

		finisher := jen.Defer().Func().Call().Block(
			jen.Id("span").Dot("Finish").Call(),
		).Call()
		/* code to generate
		defer func(){
		  span.Finish()
		}()
		*/
		lines = append(lines, jen.Line(), finisher, jen.Line())
	}
	var methodParams []jen.Code
	for _, p := range method.Parameters() {
		methodParams = append(methodParams, p.NamedParameter())
	}
	returns := jen.Return(jen.Id(service.StructPtr).Dot(service.ServicePtr).Dot(method.Name()).Params(methodParams...))
	/* code to generate
	return ${service.StructPtr}.${service.ServicePtr}.${DeclaredFunction.Name}(${DeclaredFunction.Parameters})
	*/
	lines = append(lines, returns)

	return builder.Block(lines...)
}

func (t tracer) GetRequiredImportNames() map[string]string {
	return map[string]string{
		"zipkin": "github.com/openzipkin/zipkin-go",
	}
}
