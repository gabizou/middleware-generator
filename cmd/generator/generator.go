package main

import (
	"fmt"
	"os"

	"github.com/gabizou/middleware-generator/pkg/generator"
	_ "github.com/gabizou/middleware-generator/pkg/plugins/tracing"
)

const (
	argsLengthRequirement = 4
)

func main() {
	if len(os.Args) != argsLengthRequirement {
		panic(fmt.Errorf("expected exactly one argument: <source type>"))
	}
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	g := generator.Interpret(dir, os.Args)
	g.Print()
}
