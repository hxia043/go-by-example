package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/google/cel-go/checker/decls"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"gopkg.in/yaml.v3"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"
	v1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apiextensions-apiserver/pkg/apiserver/schema"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/validation/field"
	celconfig "k8s.io/apiserver/pkg/apis/cel"

	"github.com/google/cel-go/cel"
	apiextensionscel "k8s.io/apiextensions-apiserver/pkg/apiserver/schema/cel"
)

const (
	// ScopedVarName is the variable name assigned to the locally scoped data element of a CEL validation
	// expression.
	ScopedVarName = "self"

	// OldScopedVarName is the variable name assigned to the existing value of the locally scoped data element of a
	// CEL validation expression.
	OldScopedVarName = "oldSelf"
)

func primitiveType(typ, format string) schema.Structural {
	result := schema.Structural{
		Generic: schema.Generic{
			Type: typ,
		},
	}
	if len(format) != 0 {
		result.ValueValidation = &schema.ValueValidation{
			Format: format,
		}
	}
	return result
}

var (
	stringType = primitiveType("string", "")
)

type Foo struct {
	Spec FooSpec `json:"spec"`
}

type FooSpec struct {
	Replicas int    `json:"replicas"`
	Version  string `json:"version"`
}

func objs(val ...interface{}) map[string]interface{} {
	result := make(map[string]interface{}, len(val))
	for i, v := range val {
		result[fmt.Sprintf("val%d", i+1)] = v
	}
	return result
}

func schemas(valSchema ...schema.Structural) *schema.Structural {
	result := make(map[string]schema.Structural, len(valSchema))
	for i, v := range valSchema {
		result[fmt.Sprintf("val%d", i+1)] = v
	}
	return objectTypePtr(result)
}

func withRule(s schema.Structural, rule string) schema.Structural {
	s.Extensions.XValidations = v1.ValidationRules{
		{
			Rule: rule,
		},
	}
	return s
}

func objectType(props map[string]schema.Structural) schema.Structural {
	return schema.Structural{
		Generic: schema.Generic{
			Type: "object",
		},
		Properties: props,
	}
}

func objectTypePtr(props map[string]schema.Structural) *schema.Structural {
	o := objectType(props)
	return &o
}

func mapType(valSchema *schema.Structural) schema.Structural {
	result := schema.Structural{
		Generic: schema.Generic{
			Type:                 "object",
			AdditionalProperties: &schema.StructuralOrBool{Bool: true, Structural: valSchema},
		},
	}
	return result
}

func kubernetesCEL() {
	crdYaml := `
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: foos.example.com
spec:
  group: example.com
  versions:
  - name: v1
    served: true
    storage: true
    schema:
      openAPIV3Schema:
        type: object
        properties:
          spec:
            type: object
            properties:
              replicas:
                type: integer
              version:
                type: string
            x-kubernetes-validations:
            - rule: "self.replicas >= 1 && self.replicas <= 10"
              message: "replicas must be between 1 and 10"
            - rule: "self.version.matches('^[a-zA-Z0-9]+$')"
              message: "version must match the regex pattern ^[a-zA-Z0-9]+$"
            - rule: "addNumbers(self.replicas, self.replicas)  == 10"
              message: "add numbers of replicas"
scope: Namespaced
names:
  plural: foos
  singular: foo
  kind: Foo
  shortNames:
  - f
`

	crYaml := `
apiVersion: apiextensions.k8s.io/v1
kind: Foo
metadata:
  name: foo
spec:
  replicas: 8
  version: v1
`

	var crdUnstructured unstructured.Unstructured
	if err := yaml.Unmarshal([]byte(crdYaml), &crdUnstructured.Object); err != nil {
		log.Fatalf("Error unmarshaling YAML to unstructured: %v", err)
	}

	crdJSON, err := crdUnstructured.MarshalJSON()
	if err != nil {
		log.Fatalf("Error marshaling unstructured to JSON: %v", err)
	}

	var crd v1.CustomResourceDefinition
	if err = json.Unmarshal(crdJSON, &crd); err != nil {
		log.Fatalf("Error marshaling unstructured to CustomResourceDefinition: %v", err)
	}

	convertedProps := &apiextensions.JSONSchemaProps{}
	if err := v1.Convert_v1_JSONSchemaProps_To_apiextensions_JSONSchemaProps(crd.Spec.Versions[0].Schema.OpenAPIV3Schema, convertedProps, nil); err != nil {
		log.Fatalf("Error convert json schema props: %v", err)
	}

	structural, err := schema.NewStructural(convertedProps)
	if err != nil {
		log.Fatalf("Error creating structural schema: %v", err)
	}

	celValidator := apiextensionscel.NewValidator(structural, true, celconfig.PerCallLimit)
	if celValidator == nil {
		log.Fatalf("Error: expected non nil validator")
	}

	var crUnstructured unstructured.Unstructured
	if err := yaml.Unmarshal([]byte(crYaml), &crUnstructured.Object); err != nil {
		log.Fatalf("Error unmarshaling YAML to unstructured: %v", err)
	}

	var oldObject interface{}
	errs, _ := celValidator.Validate(context.TODO(), field.NewPath("foo"), structural, crUnstructured.Object, oldObject, celconfig.RuntimeCELCostBudget)
	if len(errs) > 0 {
		for _, e := range errs {
			fmt.Printf("Validation error: %v\n", e)
		}
	} else {
		fmt.Println("Validation succeeded")
	}
}

func commonCEL() {
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
		"age":  30,
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

func addNumbers_double_double(args ...ref.Val) ref.Val {
	if len(args) != 2 {
		return types.NewErr("addNumbers requires 2 arguments")
	}

	return types.Double(args[0].(types.Double) + args[1].(types.Double))
}

func customCel() {
	env, err := cel.NewEnv(
		cel.Declarations(
			decls.NewVar("name", decls.String),
			decls.NewVar("age", decls.Int),
		),
		cel.Function("addNumbers",
			cel.Overload(
				"addNumbers_double_double",
				[]*cel.Type{cel.DoubleType, cel.DoubleType},
				cel.DoubleType,
				cel.FunctionBinding(addNumbers_double_double)),
		),
	)
	if err != nil {
		log.Fatalf("env create failed: %v", err)
	}

	//expression := "name == 'Alice' && age > 20 && addNumbers(number, number)"
	expression := "addNumbers(3.0, 4.0)"
	ast, iss := env.Compile(expression)
	if iss.Err() != nil {
		log.Fatalf("compile expression failed: %v", iss.Err())
	}

	prg, err := env.Program(ast)
	if err != nil {
		log.Fatalf("create programe failed: %v", err)
	}

	input := map[string]interface{}{
		"name":   "Alice",
		"age":    30,
		"number": 4.3,
	}

	out, _, err := prg.Eval(input)
	if err != nil {
		log.Fatalf("execute failed: %v", err)
	}

	fmt.Println(out)

	if out == types.True {
		fmt.Println("expression result: true")
	} else {
		fmt.Println("exepression result: false")
	}
}

func defaultCEL() {
	tests := []struct {
		name          string
		schema        *schema.Structural
		oldSchema     *schema.Structural
		obj           interface{}
		oldObj        interface{}
		valid         []string
		errors        map[string]string // rule -> string that error message must contain
		costBudget    int64
		isRoot        bool
		expectSkipped bool
	}{
		{name: "enums",
			obj: map[string]interface{}{"enumStr": "Pending"},
			schema: objectTypePtr(map[string]schema.Structural{"enumStr": {
				Generic: schema.Generic{
					Type: "string",
				},
				ValueValidation: &schema.ValueValidation{
					Enum: []schema.JSON{
						{Object: "Pending"},
						{Object: "Available"},
						{Object: "Bound"},
						{Object: "Released"},
						{Object: "Failed"},
					},
				},
			}}),
			valid: []string{
				"self.enumStr == 'Pending'",
				"self.enumStr in ['Pending', 'Available']",
			},
		},
		{name: "maps",
			obj:    objs(map[string]interface{}{"k1": "a", "k2": "b"}, map[string]interface{}{"k2": "b", "k1": "a"}),
			schema: schemas(mapType(&stringType), mapType(&stringType)),
			valid: []string{
				"self.val1 == self.val2", // equal even though order is different
				"'k1' in self.val1",
				"!('k3' in self.val1)",
				"self.val1 == {'k1': 'a', 'k2': 'b'}",
			},
			errors: map[string]string{
				// Mixed type maps are not allowed since we have HomogeneousAggregateLiterals enabled
				"{'k1': 'a', 'k2': 1, 'k2': 'b'}":     "expected type 'string' but found 'int'",
				"{'k1': 'a', 'k2': 'b', 'k2': false}": "expected type 'string' but found 'bool'",
			},
		},
	}

	for i := range tests {
		tt := tests[i]
		for j := range tt.valid {
			validRule := tt.valid[j]
			s := withRule(*tt.schema, validRule)
			celValidator := apiextensionscel.NewValidator(&s, tt.isRoot, celconfig.PerCallLimit)
			if celValidator == nil {
				fmt.Println("expected non nil validator")
			}

			errs, remainingBudget := celValidator.Validate(context.TODO(), field.NewPath("root"), &s, tt.obj, tt.oldObj, celconfig.RuntimeCELCostBudget)
			for _, err := range errs {
				fmt.Printf("unexpected error: %v\n", err)
			}
			if tt.expectSkipped {
				// Skipped validations should have no cost. The only possible false positive here would be the CEL expression 'true'.
				if remainingBudget != celconfig.RuntimeCELCostBudget {
					fmt.Printf("expected no cost expended for skipped validation, but got %d remaining from %d budget\n", remainingBudget, celconfig.RuntimeCELCostBudget)
				}
				return
			}
		}

		for rule, expectErrToContain := range tt.errors {
			s := withRule(*tt.schema, rule)
			celValidator := apiextensionscel.NewValidator(&s, tt.isRoot, celconfig.PerCallLimit)
			if celValidator == nil {
				fmt.Println("expected non nil validator")
			}

			errs, _ := celValidator.Validate(context.TODO(), field.NewPath("root"), &s, tt.obj, tt.oldObj, celconfig.RuntimeCELCostBudget)
			if len(errs) == 0 {
				fmt.Println("expected validation errors but got none")
			}

			for _, err := range errs {
				if strings.Contains(err.Error(), expectErrToContain) {
					fmt.Printf("expected error to contain '%s', but got: %v\n", expectErrToContain, err)
				}
			}
		}
	}
}

func main() {
	flag := "kubernetes"
	switch flag {
	case "common":
		commonCEL()
	case "kubernetes":
		kubernetesCEL()
	case "default":
		defaultCEL()
	case "custom":
		customCel()
	default:
		fmt.Println("unknown flag")
	}
}
