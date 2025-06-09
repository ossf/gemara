package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"reflect"
	"strings"
)

func main() {
	const filename = "./layer2/generated_types.go"

	src, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, filename, src, parser.ParseComments)
	if err != nil {
		panic(err)
	}

	ast.Inspect(file, func(n ast.Node) bool {
		ts, ok := n.(*ast.TypeSpec)
		if !ok {
			return true
		}
		st, ok := ts.Type.(*ast.StructType)
		if !ok {
			return true
		}
		for _, field := range st.Fields.List {
			if field.Tag == nil {
				continue
			}
			tag := reflect.StructTag(strings.Trim(field.Tag.Value, "`"))
			jsonTag := tag.Get("json")
			if jsonTag == "" {
				continue
			}
			if tag.Get("yaml") != "" {
				continue // already has yaml tag
			}
			// Compose new tag string
			tagParts := []string{}
			for _, part := range strings.Split(string(tag), " ") {
				if strings.HasPrefix(part, "yaml:") {
					continue
				}
				tagParts = append(tagParts, part)
			}
			tagParts = append(tagParts, fmt.Sprintf(`yaml:"%s"`, jsonTag))
			newTag := "`" + strings.Join(tagParts, " ") + "`"
			field.Tag.Value = newTag
		}
		return true
	})

	var buf bytes.Buffer
	if err := printer.Fprint(&buf, fset, file); err != nil {
		panic(err)
	}
	if err := os.WriteFile(filename, buf.Bytes(), 0644); err != nil {
		panic(err)
	}
	fmt.Println("YAML tags added to", filename)
}
