package main

import (
	"fmt"
	"log"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/checker/decls"
	"github.com/google/cel-go/common/types"
)

func main() {
	env, err := cel.NewEnv(
		cel.Declarations(
			decls.NewVar("name", decls.String),
			decls.NewVar("age", decls.Int),
		),
	)
	if err != nil {
		log.Fatalf("env create failed: %v", err)
	}

	expression := "name == 'Alice' && age > 20"
	ast, iss := env.Compile(expression)
	if iss.Err() != nil {
		log.Fatalf("compile expression failed: %v", iss.Err())
	}

	prg, err := env.Program(ast)
	if err != nil {
		log.Fatalf("create programe failed: %v", err)
	}

	input := map[string]interface{}{
		"name": "Alice",
		"age":  10,
	}

	out, _, err := prg.Eval(input)
	if err != nil {
		log.Fatalf("execute failed: %v", err)
	}

	if out == types.True {
		fmt.Println("expression result: true")
	} else {
		fmt.Println("exepression result: false")
	}
}
