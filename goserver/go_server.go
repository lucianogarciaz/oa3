package goserver

import (
	"bytes"
	"fmt"
	"go/format"
	"sort"
	"strings"
	"text/template"
	"unicode"

	"github.com/aarondl/oa3/generator"
	"github.com/aarondl/oa3/openapi3spec"
	"github.com/aarondl/oa3/templates"
)

const (
	// DefaultPackage name for go
	DefaultPackage = "oa3gen"
	// Disclaimer printed to the top of Go files
	Disclaimer = `// Code generated by oa3 (https://github.com/aarondl/oa3). DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.`
)

// templates for generation
var tpls = []string{
	"schema.tpl",
	"schema_top.tpl",
}

// funcs to use for generation
var funcs = map[string]interface{}{
	"camelSnake":        camelSnake,
	"named":             named,
	"primitive":         primitive,
	"isInlinePrimitive": isInlinePrimitive,
}

// templateData for go templates
type templateData struct {
	Name   string
	Object interface{}
	*templates.TemplateData
}

// generator generates templates for Go
type gen struct {
	tpl *template.Template
}

// New go generator
func New() generator.Interface {
	return &gen{}
}

// Load templates
func (g *gen) Load(dir string) error {
	var err error
	g.tpl, err = templates.Load(funcs, dir, tpls...)
	return err
}

// Do generation for Go.
func (g *gen) Do(spec *openapi3spec.OpenAPI3, params map[string]string) ([]generator.File, error) {
	var files []generator.File
	f, err := generateTopLevelSchemas(spec, params, g.tpl)
	if err != nil {
		return nil, err
	}

	files = append(files, f...)

	for i, f := range files {
		formatted, err := format.Source(f.Contents)
		if err != nil {
			return nil, fmt.Errorf("failed to format file(%s): %w\n%s", f.Name, err, f.Contents)
		}

		files[i].Contents = formatted
	}

	return files, nil
}

// generateSchemas creates files for the topLevel-level referenceable types
func generateTopLevelSchemas(spec *openapi3spec.OpenAPI3, params map[string]string, tpl *template.Template) ([]generator.File, error) {
	if spec.Components == nil {
		return nil, nil
	}

	keys := make([]string, 0, len(spec.Components.Schemas))
	for k := range spec.Components.Schemas {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	topLevelStructs := make([]generator.File, 0, len(keys))

	for _, k := range keys {
		v := spec.Components.Schemas[k]

		// Don't generate arrays as additional types
		if v.Type == "array" && v.Items != nil && len(v.Items.Ref) != 0 {
			continue
		}

		filename := "schema_" + camelSnake(k) + ".go"

		tData := templates.NewTemplateData(spec, params)
		data := templateData{
			TemplateData: tData,
			Name:         k,
			Object:       v,
		}

		buf := new(bytes.Buffer)
		if err := tpl.ExecuteTemplate(buf, "schema_top", data); err != nil {
			return nil, fmt.Errorf("failed rendering template %q: %w", "schema", err)
		}

		fileBytes := new(bytes.Buffer)
		pkg := DefaultPackage
		if pkgParam := params["package"]; len(pkgParam) > 0 {
			pkg = pkgParam
		}

		fileBytes.WriteString(Disclaimer)
		fmt.Fprintf(fileBytes, "\npackage %s\n", pkg)
		if imps := imports(data.Imports); len(imps) != 0 {
			fileBytes.WriteByte('\n')
			fileBytes.WriteString(imports(data.Imports))
			fileBytes.WriteByte('\n')
		}
		fileBytes.WriteByte('\n')
		fileBytes.Write(buf.Bytes())

		topLevelStructs = append(topLevelStructs, generator.File{Name: filename, Contents: fileBytes.Bytes()})
	}

	return topLevelStructs, nil
}

func isInlinePrimitive(schema *openapi3spec.Schema) bool {
	if schema.Type == "object" {
		return schema.AdditionalProperties != nil
	}

	return true
}

func primitive(tdata templateData, schema *openapi3spec.Schema) (string, error) {
	if schema.Nullable {
		return primitiveNil(tdata, schema)
	}

	return primitiveNonNil(tdata, schema)
}

func primitiveNonNil(tdata templateData, schema *openapi3spec.Schema) (string, error) {
	switch schema.Type {
	case "integer":
		if schema.Format != nil {
			switch *schema.Format {
			case "int32":
				return "int32", nil
			case "int64":
				return "int64", nil
			}
		}

		return "int", nil
	case "number":
		if schema.Format != nil {
			switch *schema.Format {
			case "float":
				return "float32", nil
			case "double":
				return "float64", nil
			}
		}

		return "float64", nil
	case "string":
		return "string", nil
	case "boolean":
		return "bool", nil
	}

	return "", fmt.Errorf("schema had unsupported type: %s", schema.Type)
}

func primitiveNil(tdata templateData, schema *openapi3spec.Schema) (string, error) {
	switch schema.Type {
	case "integer":
		tdata.Import("github.com/volatiletech/null")

		if schema.Format != nil {
			switch *schema.Format {
			case "int32":
				return "null.Int32", nil
			case "int64":
				return "null.Int64", nil
			}
		}

		return "null.Int", nil
	case "number":
		tdata.Import("github.com/volatiletech/null")

		if schema.Format != nil {
			switch *schema.Format {
			case "float":
				return "null.Float32", nil
			case "double":
				return "null.Float64", nil
			}
		}

		return "null.Float64", nil
	case "string":
		tdata.Import("github.com/volatiletech/null")
		return "null.String", nil
	case "boolean":
		tdata.Import("github.com/volatiletech/null")
		return "null.Bool", nil
	}

	return "", fmt.Errorf("schema had unsupported nil type: %s", schema.Type)
}

func imports(imps map[string]struct{}) string {
	if len(imps) == 0 {
		return ""
	}

	var std, third []string
	for imp := range imps {
		splits := strings.Split(imp, "/")
		if len(splits) > 0 && strings.ContainsRune(splits[0], '.') {
			third = append(third, imp)
			continue
		}

		std = append(std, imp)
	}

	sort.Strings(std)
	sort.Strings(third)

	buf := new(bytes.Buffer)
	buf.WriteString("import (")
	for _, imp := range std {
		fmt.Fprintf(buf, "\n\t\"%s\"", imp)
	}
	if len(std) != 0 && len(third) != 0 {
		buf.WriteByte('\n')
	}
	for _, imp := range third {
		fmt.Fprintf(buf, "\n\t\"%s\"", imp)
	}
	buf.WriteString("\n)")

	return buf.String()
}

// schema_UserIDProfile -> schema_user_id_profile
// ID -> id
func camelSnake(filename string) string {
	build := new(strings.Builder)

	var upper bool

	in := []rune(filename)
	for i, r := range []rune(in) {
		if !unicode.IsLetter(r) {
			upper = false
			build.WriteRune(r)
			continue
		}

		if !unicode.IsUpper(r) {
			upper = false
			build.WriteRune(r)
			continue
		}

		addUnderscore := false
		if upper {
			if i+1 < len(in) && unicode.IsLower(in[i+1]) {
				addUnderscore = true
			}
		} else {
			if i-1 > 0 && unicode.IsLetter(in[i-1]) {
				addUnderscore = true
			}
		}

		if addUnderscore {
			build.WriteByte('_')
		}

		upper = true
		build.WriteRune(unicode.ToLower(r))
	}

	return build.String()
}

func named(tplData templateData, nextName string, nextObj interface{}) templateData {
	return templateData{
		TemplateData: tplData.TemplateData,
		Name:         tplData.Name + nextName,
		Object:       nextObj,
	}
}
