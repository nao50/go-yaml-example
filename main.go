package main

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/goccy/go-yaml"
	"github.com/goccy/go-yaml/parser"
	"github.com/goccy/go-yaml/printer"
)

type API struct {
	A int
	B string
}

var yml = `---
a: 10
b: "small"
`

func main() {
	var api API

	// spec, _ := ioutil.ReadFile("./spec/api.yaml")
	// buf := bytes.NewBuffer(spec)

	validate := validator.New()
	validate.RegisterStructValidation(custom_validation, api)

	dec := yaml.NewDecoder(
		// buf,
		// yaml.RecursiveDir(true),
		// yaml.ReferenceDirs("spec"),
		strings.NewReader(yml),
		// yaml.Validator(validate), // <- ここで呼べないので
	)

	err := dec.Decode(&api)
	if err != nil {
		fmt.Println(yaml.FormatError(err, true, true))
	}

	// ここで呼ぶ
	err = validate.Struct(api)
	if err != nil {
		fmt.Println(err)
	}

}

//
func custom_validation(sl validator.StructLevel) {
	api := sl.Current().Interface().(API)
	if api.A > 5 && api.B == "small" {
		source, err := yamlSourceByPath(yml, "$.b")
		if err != nil {
			panic(err)
		}
		fmt.Printf("b value expected \"large\" but actual %s:\n%s\n", api.B, source)
	}
}

//
func yamlSourceByPath(originalSource string, pathStr string) (string, error) {
	file, err := parser.ParseBytes([]byte(originalSource), 0)
	if err != nil {
		return "", err
	}
	path, err := yaml.PathString(pathStr)
	if err != nil {
		return "", err
	}
	node, err := path.FilterFile(file)
	if err != nil {
		return "", err
	}
	var p printer.Printer
	return p.PrintErrorToken(node.GetToken(), true), nil
}
