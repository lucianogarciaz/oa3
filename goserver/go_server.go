package goserver

import (
	"bytes"
	"errors"
	"fmt"
	"go/format"
	"io/fs"
	"sort"
	"strconv"
	"strings"
	"text/template"
	"unicode"

	"github.com/aarondl/oa3/generator"
	"github.com/aarondl/oa3/openapi3spec"
	"github.com/aarondl/oa3/templates"
	"github.com/huandu/xstrings"
)

const (
	// DefaultPackage name for go
	DefaultPackage = "oa3gen"
	// Disclaimer printed to the top of Go files
	Disclaimer = `// Code generated by oa3 (https://github.com/aarondl/oa3). DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.`
)

// Constants for keys recognized in the parameters for the Go server
const (
	ParamKeyPackage     = "package"
	ParamKeyTimeType    = "timetype"
	ParamKeyUUIDType    = "uuidtype"
	ParamKeyDecimalType = "decimaltype"
)

// templates for generation
var TemplateList = []string{
	"api_interface.tpl",
	"api_methods.tpl",
	"responses.tpl",
	"schema.tpl",
	"schema_top.tpl",

	"validate_schema.tpl",
	"validate_field.tpl",
}

// TemplateFunctions to use for generation
var TemplateFunctions = map[string]any{
	"camelSnake":              camelSnake,
	"filterNonIdentChars":     filterNonIdentChars,
	"isInlinePrimitive":       isInlinePrimitive,
	"omitnullWrap":            omitnullWrap,
	"omitnullUnwrap":          omitnullUnwrap,
	"omitnullIsWrapped":       omitnullIsWrapped,
	"paramConvertFn":          paramConvertFn,
	"paramRequiresType":       paramRequiresType,
	"paramSchemaName":         paramSchemaName,
	"paramTypeName":           paramTypeName,
	"primitive":               primitive,
	"primitiveBits":           primitiveBits,
	"primitiveWrapped":        primitiveWrapped,
	"responseKind":            responseKind,
	"snakeToCamel":            snakeToCamel,
	"taggedPaths":             tagPaths,
	"hasComplexServers":       hasComplexServers,
	"hasJSONResponse":         hasJSONResponse,
	"responseTypeName":        responseTypeName,
	"responseNeedsWrap":       responseNeedsWrap,
	"responseNeedsCodeWrap":   responseNeedsCodeWrap,
	"responseNeedsHeaderWrap": responseNeedsHeaderWrap,
	"responseNeedsPtr":        responseNeedsPtr,
	"responseRefName":         responseRefName,

	// overrides of the defaults
	"mustValidate":        mustValidate,
	"mustValidateRecurse": mustValidateRecurse,
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
func (g *gen) Load(dir fs.FS) error {
	var err error
	g.tpl, err = templates.Load(TemplateFunctions, dir, TemplateList...)
	return err
}

// Do generation for Go.
func (g *gen) Do(spec *openapi3spec.OpenAPI3, params map[string]string) ([]generator.File, error) {
	if params == nil {
		params = make(map[string]string)
	}
	if pkg, ok := params[ParamKeyPackage]; !ok || len(pkg) == 0 {
		params[ParamKeyPackage] = DefaultPackage
	}

	var files []generator.File
	f, err := GenerateTopLevelSchemas(spec, params, g.tpl)
	if err != nil {
		return nil, fmt.Errorf("failed to generate schemas: %w", err)
	}

	files = append(files, f...)

	f, err = generateAPIInterface(spec, params, g.tpl)
	if err != nil {
		return nil, fmt.Errorf("failed to generate api interface: %w", err)
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

func generateAPIInterface(spec *openapi3spec.OpenAPI3, params map[string]string, tpl *template.Template) ([]generator.File, error) {
	if spec.Paths == nil {
		return nil, nil
	}

	files := make([]generator.File, 0)

	apiName := strings.Title(strings.ReplaceAll(spec.Info.Title, " ", "")) //nolint:staticcheck

	data := templates.NewTemplateDataWithObject(spec, params, apiName, nil, false)

	filename := generator.FilenameFromTitle(spec.Info.Title) + ".go"

	buf := new(bytes.Buffer)
	if err := tpl.ExecuteTemplate(buf, "api_interface", data); err != nil {
		return nil, fmt.Errorf("failed rendering template %q: %w", "schema", err)
	}

	fileBytes := new(bytes.Buffer)
	pkg := params["package"]

	fileBytes.WriteString(Disclaimer)
	fmt.Fprintf(fileBytes, "\npackage %s\n", pkg)
	if imps := Imports(data.Imports); len(imps) != 0 {
		fileBytes.WriteByte('\n')
		fileBytes.WriteString(Imports(data.Imports))
		fileBytes.WriteByte('\n')
	}
	fileBytes.WriteByte('\n')
	fileBytes.Write(buf.Bytes())

	content := make([]byte, len(fileBytes.Bytes()))
	copy(content, fileBytes.Bytes())
	files = append(files, generator.File{Name: filename, Contents: content})

	data = templates.NewTemplateDataWithObject(spec, params, apiName, nil, false)
	filename = generator.FilenameFromTitle(spec.Info.Title) + "_methods.go"

	buf.Reset()
	fileBytes.Reset()
	if err := tpl.ExecuteTemplate(buf, "api_methods", data); err != nil {
		return nil, fmt.Errorf("failed rendering template %q: %w", "schema", err)
	}

	fileBytes.WriteString(Disclaimer)
	fmt.Fprintf(fileBytes, "\npackage %s\n", pkg)
	if imps := Imports(data.Imports); len(imps) != 0 {
		fileBytes.WriteByte('\n')
		fileBytes.WriteString(Imports(data.Imports))
		fileBytes.WriteByte('\n')
	}
	fileBytes.WriteByte('\n')
	fileBytes.Write(buf.Bytes())

	files = append(files, generator.File{Name: filename, Contents: fileBytes.Bytes()})

	return files, nil
}

// GenerateSchemas creates files for the topLevel-level referenceable types
//
// Some supported Inline are also generated.
// Prefixed with their recursive names and Inline.
// components.responses[name].headers[headername].schema
// components.responses[name].content[mime-type].schema
// components.responses[name].content[mime-type].encoding[propname].headers[headername].schema
// components.parameters[name].schema
// components.requestBodies[name].content[mime-type].schema
// components.requestBodies[name].content[mime-type].encoding[propname].headers[headername].schema
// components.headers[name].schema
// paths.parameters[0].schema
// paths.(get|put...).parameters[0].schema
// paths.(get|put...).requestBody.content[mime-type].schema
// paths.(get|put...).responses[name].headers[headername].schema
// paths.(get|put...).responses[name].content[mime-type].schema
// paths.(get|put...).responses[name].content[mime-type].encoding[propname].headers[headername].schema
func GenerateTopLevelSchemas(spec *openapi3spec.OpenAPI3, params map[string]string, tpl *template.Template) ([]generator.File, error) {
	if spec.Components == nil {
		return nil, nil
	}

	var err error
	keys := make([]string, 0, len(spec.Components.Schemas))
	for k := range spec.Components.Schemas {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	topLevelStructs := make([]generator.File, 0, len(keys))

	for _, k := range keys {
		v := spec.Components.Schemas[k]
		filename := "schema_" + camelSnake(k) + ".go"
		topLevelStructs, err = recurseSchemas(spec, params, tpl, topLevelStructs, filename, k, v, false)
		if err != nil {
			return nil, err
		}
	}

	type opMap struct {
		Verb string
		Op   *openapi3spec.Operation
	}

	for _, p := range spec.Paths {
		opMaps := []opMap{
			{"GET", p.Get},
			{"POST", p.Post},
			{"PUT", p.Put},
			{"PATCH", p.Patch},
			{"TRACE", p.Trace},
			{"HEAD", p.Head},
			{"DELETE", p.Delete},
		}
		for _, o := range opMaps {
			if o.Op == nil {
				continue
			}

			for _, p := range o.Op.Parameters {
				if p == nil {
					continue
				}

				if p.Content != nil {
					return nil, errors.New("oa3 go server/client does not support content parameters yet")
				} else if p.Schema.Type == "object" || (p.Schema.Type == "array" && (p.Schema.Items.Type == "object" || p.Schema.Items.Type == "array")) {
					return nil, errors.New("oa3 go server/client does not support object serialization for parameters yet, only primitives or arrays of primitives (including enum)")
				}

				if !paramRequiresType(*p) {
					continue
				}

				filename := fmt.Sprintf("schema_%s_%s_%s_param.go", camelSnake(o.Op.OperationID), strings.ToLower(o.Verb), camelSnake(p.Name))
				topLevelStructs, err = recurseSchemas(spec, params, tpl, topLevelStructs, filename, paramSchemaName(o.Op.OperationID, o.Verb, p.Name), p.Schema, p.Required)
				if err != nil {
					return nil, err
				}
			}
		}
		for _, o := range opMaps {
			if o.Op == nil {
				continue
			}

			// If we have no request body ignore this op
			if o.Op.RequestBody != nil {
				if len(o.Op.RequestBody.Ref) != 0 {
					continue
				}

				json := o.Op.RequestBody.Content["application/json"]
				if json == nil {
					continue
				}

				schema := json.Schema
				// Refs are taken care of already
				if len(schema.Ref) != 0 {
					continue
				}

				filename := "schema_" + camelSnake(o.Op.OperationID) + "_reqbody.go"
				name := strings.Title(o.Op.OperationID) + "Inline" //nolint:staticcheck
				topLevelStructs, err = recurseSchemas(spec, params, tpl, topLevelStructs, filename, name, &schema, o.Op.RequestBody.Required)
				if err != nil {
					return nil, err
				}
			}
		}

		for _, o := range opMaps {
			if o.Op == nil {
				continue
			}

			for code, resp := range o.Op.Responses {
				if len(resp.Ref) != 0 {
					continue
				}
				if len(resp.Content) == 0 {
					continue
				}

				json := resp.Content["application/json"]
				if json == nil {
					continue
				}

				schema := json.Schema
				// Refs are taken care of already
				if len(schema.Ref) != 0 {
					continue
				}

				filename := "schema_" + camelSnake(o.Op.OperationID) + "_" + code + "_respbody.go"
				name := strings.Title(o.Op.OperationID) + strings.Title(code) + "Inline" //nolint:staticcheck
				topLevelStructs, err = recurseSchemas(spec, params, tpl, topLevelStructs, filename, name, &schema, true)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	for name, req := range spec.Components.RequestBodies {
		reqMedia, ok := req.Content["application/json"]
		if !ok {
			continue
		}

		schema := reqMedia.Schema
		if len(schema.Ref) != 0 {
			continue
		}

		filename := "schema_" + camelSnake(name) + "_reqbody.go"
		topLevelStructs, err = recurseSchemas(spec, params, tpl, topLevelStructs, filename, name+"Inline", &schema, req.Required)
		if err != nil {
			return nil, err
		}
	}

	for name, resp := range spec.Components.Responses {
		respMedia, ok := resp.Content["application/json"]
		if !ok {
			continue
		}

		if len(resp.Content) == 0 {
			continue
		}

		schema := respMedia.Schema
		if len(schema.Ref) != 0 {
			continue
		}

		filename := "schema_" + camelSnake(name) + "_respbody.go"
		topLevelStructs, err = recurseSchemas(spec, params, tpl, topLevelStructs, filename, name+"Inline", &schema, true)
		if err != nil {
			return nil, err
		}
	}

	return topLevelStructs, nil
}

// recurseSchemas checks for any embedded structs or enums that should
// be brought to the top level
func recurseSchemas(spec *openapi3spec.OpenAPI3, params map[string]string, tpl *template.Template, topLevelStructs []generator.File, filename, name string, ref *openapi3spec.SchemaRef, required bool) ([]generator.File, error) {
	if ref.IsRef() {
		return topLevelStructs, nil
	}

	var err error
	switch ref.Type {
	case "object":
		if ref.AdditionalProperties != nil && ref.AdditionalProperties.SchemaRef != nil {
			topLevelStructs, err = recurseSchemas(spec, params, tpl, topLevelStructs, "", name+"Item", ref.AdditionalProperties.SchemaRef, true)
			if err != nil {
				return nil, err
			}
		}

		for propname, prop := range ref.Properties {
			objname := name + snakeToCamel(strings.Title(propname)) //nolint:staticcheck
			topLevelStructs, err = recurseSchemas(spec, params, tpl, topLevelStructs, "", objname, prop, ref.IsRequired(propname))
			if err != nil {
				return nil, err
			}
		}
	case "array":
		topLevelStructs, err = recurseSchemas(spec, params, tpl, topLevelStructs, "", name+"Item", ref.Items, true)
		if err != nil {
			return nil, err
		}
	}

	if ref.Schema == nil {
		return topLevelStructs, nil
	}

	// If the filename has been set that means we're intent on writing this file
	// no matter what, else if there's an enum present or it's not an inline
	// primitive we have to render a type out for it.
	if len(filename) != 0 || ref.HasEnum() || !isInlinePrimitive(ref.Schema) {
		if len(filename) == 0 {
			filename = "schema_" + camelSnake(name) + ".go"
		}
		generated, err := makePseudoFile(spec, params, tpl, filename, name, ref, required)
		if err != nil {
			return nil, err
		}
		topLevelStructs = append(topLevelStructs, generated)
	}

	return topLevelStructs, nil
}

var (
	fileBuf   = new(bytes.Buffer)
	headerBuf = new(bytes.Buffer)
)

func makePseudoFile(spec *openapi3spec.OpenAPI3, params map[string]string, tpl *template.Template, filename string, name string, schema *openapi3spec.SchemaRef, required bool) (generator.File, error) {
	fileBuf.Reset()
	headerBuf.Reset()

	data := templates.NewTemplateDataWithObject(spec, params, name, schema, required)

	if err := tpl.ExecuteTemplate(fileBuf, "schema_top", data); err != nil {
		return generator.File{}, fmt.Errorf("failed rendering template %q: %w", "schema", err)
	}

	pkg := DefaultPackage
	if pkgParam := params["package"]; len(pkgParam) > 0 {
		pkg = pkgParam
	}

	headerBuf.WriteString(Disclaimer)
	fmt.Fprintf(headerBuf, "\npackage %s\n", pkg)
	if imps := Imports(data.Imports); len(imps) != 0 {
		headerBuf.WriteByte('\n')
		headerBuf.WriteString(Imports(data.Imports))
		headerBuf.WriteByte('\n')
	}
	headerBuf.WriteByte('\n')

	headerLen, fileLen := headerBuf.Len(), fileBuf.Len()
	contents := make([]byte, headerLen+fileLen)
	copy(contents, headerBuf.Bytes())
	copy(contents[headerLen:], fileBuf.Bytes())

	return generator.File{Name: filename, Contents: contents}, nil
}

func isInlinePrimitive(schema *openapi3spec.Schema) bool {
	if schema.Type == "object" {
		return schema.AdditionalProperties != nil
	}

	return true
}

func primitive(tdata templates.TemplateData, schema *openapi3spec.Schema) (string, error) {
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
		if schema.Format != nil {
			switch *schema.Format {
			case "date":
				if tdata.Params[ParamKeyTimeType] == "chrono" {
					tdata.Import("github.com/aarondl/chrono")
					return "chrono.Date", nil
				} else {
					tdata.Import("time")
					return "time.Time", nil
				}
			case "time":
				if tdata.Params[ParamKeyTimeType] == "chrono" {
					tdata.Import("github.com/aarondl/chrono")
					return "chrono.Time", nil
				} else {
					tdata.Import("time")
					return "time.Time", nil
				}
			case "date-time":
				if tdata.Params[ParamKeyTimeType] == "chrono" {
					tdata.Import("github.com/aarondl/chrono")
					return "chrono.DateTime", nil
				} else {
					tdata.Import("time")
					return "time.Time", nil
				}
			case "duration":
				tdata.Import("time")
				return "time.Duration", nil
			case "uuid":
				if tdata.Params[ParamKeyUUIDType] == "google" {
					tdata.Import("github.com/google/uuid")
					return "uuid.UUID", nil
				}
			case "decimal":
				if tdata.Params[ParamKeyDecimalType] == "shopspring" {
					tdata.Import("github.com/shopspring/decimal")
					return "decimal.Decimal", nil
				}
			}
		}
		return "string", nil
	case "boolean":
		return "bool", nil
	}

	return "", fmt.Errorf("schema expected primitive type (integer, number, string, boolean) but got: %s", schema.Type)
}

func primitiveBits(tdata templates.TemplateData, schema *openapi3spec.Schema) (string, error) {
	s, err := primitive(tdata, schema)
	if err != nil {
		return "", nil
	}

	ret := strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) {
			return -1
		}
		return r
	}, s)

	if len(ret) == 0 {
		return "64", nil
	} else {
		return ret, nil
	}
}

func primitiveWrapped(tdata templates.TemplateData, schema *openapi3spec.Schema, nullable bool, required bool) (string, error) {
	prim, err := primitive(tdata, schema)
	if err != nil {
		return "", err
	}

	return omitnullWrap(tdata, prim, nullable, required), nil
}

func omitnullWrap(tdata templates.TemplateData, typ string, nullable bool, required bool) string {
	var kind string
	switch {
	case !nullable && required:
		return typ
	case nullable && required:
		kind = "null"
	case nullable && !required:
		kind = "omitnull"
	case !nullable && !required:
		kind = "omit"
	}

	tdata.Import("github.com/aarondl/opt/" + kind)
	return kind + `.Val[` + typ + `]`
}

func omitnullUnwrap(name string, nullable bool, required bool) string {
	switch {
	case !nullable && required:
		return name
	default:
		return name + ".GetOrZero()"
	}
}

func omitnullIsWrapped(nullable bool, required bool) bool {
	return nullable || !required
}

// mustValidateRecurse checks to see if the current schema, or any sub-schema
// requires validation
func mustValidateRecurse(tdata templates.TemplateData, s *openapi3spec.Schema) bool {
	cycleMarkers := make(map[string]struct{})
	return mustValidateRecurseHelper(tdata, s, cycleMarkers)
}

func mustValidateRecurseHelper(tdata templates.TemplateData, s *openapi3spec.Schema, visited map[string]struct{}) bool {
	if mustValidate(tdata, s) {
		return true
	}

	if s.Type == "array" {
		if len(s.Items.Ref) != 0 {
			if _, ok := visited[s.Items.Ref]; ok {
				return false
			}
			visited[s.Items.Ref] = struct{}{}
		}

		return mustValidateRecurseHelper(tdata, s.Items.Schema, visited)
	} else if s.Type == "object" {
		mustV := false
		if s.AdditionalProperties != nil {
			if len(s.AdditionalProperties.Ref) != 0 {
				if _, ok := visited[s.AdditionalProperties.Ref]; !ok {
					visited[s.AdditionalProperties.Ref] = struct{}{}
					mustV = mustV || mustValidateRecurseHelper(tdata, s.AdditionalProperties.Schema, visited)
				}
			} else {
				mustV = mustV || mustValidateRecurseHelper(tdata, s.AdditionalProperties.Schema, visited)
			}
		}

		for _, v := range s.Properties {
			if len(v.Ref) != 0 {
				if _, ok := visited[v.Ref]; ok {
					continue
				}
				visited[v.Ref] = struct{}{}
			}

			mustV = mustV || mustValidateRecurseHelper(tdata, v.Schema, visited)
		}

		return mustV
	}

	return false
}

// mustValidate checks to see if the schema requires any kind of validation
// overrides the generic definition
//
// The general reason to have overridded this is because
// date/datetime/time/duration types are handled by using a type that validates
// it on parse/convert so there's no reason to generate validation after
// the conversion has already been done.
func mustValidate(tdata templates.TemplateData, s *openapi3spec.Schema) bool {
	return s.MultipleOf != nil ||
		s.Maximum != nil ||
		s.Minimum != nil ||
		s.MaxLength != nil ||
		s.MinLength != nil ||
		s.Pattern != nil ||
		s.MaxItems != nil ||
		s.MinItems != nil ||
		s.UniqueItems != nil ||
		s.MaxProperties != nil ||
		s.MinProperties != nil ||
		len(s.Enum) > 0 ||
		(s.Format != nil && *s.Format == "uuid" && !tdata.TemplateParamEquals("uuidtype", "google")) ||
		(s.Format != nil && *s.Format == "decimal" && !tdata.TemplateParamEquals("decimaltype", "shopspring"))
}

func Imports(imps map[string]struct{}) string {
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

func snakeToCamel(in string) string {
	build := new(strings.Builder)
	sawUnderscore := false
	for _, r := range in {
		if r == '_' {
			sawUnderscore = true
			continue
		}

		if sawUnderscore {
			sawUnderscore = false
			build.WriteRune(unicode.ToUpper(r))
		} else {
			build.WriteRune(r)
		}
	}

	return build.String()
}

func filterNonIdentChars(in string) string {
	build := new(strings.Builder)
	for _, c := range in {
		if unicode.IsLetter(c) || c == '_' {
			build.WriteRune(c)
		}
	}

	return build.String()
}

func responseNeedsWrap(op *openapi3spec.Operation, code string) bool {
	return responseNeedsHeaderWrap(op, code) || responseNeedsCodeWrap(op, code)
}

// checks to see if a response should be wrapped
func responseNeedsCodeWrap(op *openapi3spec.Operation, code string) bool {
	r := op.Responses[code]

	for responseCode, otherResponse := range op.Responses {
		if responseCode == code {
			continue
		}

		// Look for dupe schemas across all the codes
		for _, current := range r.Content {
			for _, other := range otherResponse.Content {

				if len(current.Schema.Ref) != 0 && len(other.Schema.Ref) != 0 && current.Schema.Ref == other.Schema.Ref {
					// If they're both refs to the same thing
					return true
				} else if len(current.Schema.Ref) == 0 && len(other.Schema.Ref) == 0 && current.Schema.Type == other.Schema.Type {
					// If they're both not refs to the same basic type
					if isInlinePrimitive(current.Schema.Schema) {
						return true
					}
				}
			}
		}
	}

	return false
}

func responseNeedsHeaderWrap(op *openapi3spec.Operation, code string) bool {
	return len(op.Responses[code].Headers) > 0
}

// responseKind returns the type of abstraction we need for a specific response
// code in an operation.
//
// The return can be one of three values: "wrapped" "empty" or ""
//
// Wrapped indicates it must be wrapped in a struct because it either has
// headers or it is a duplicate response type (say two strings) and we need
// to differentiate on code.
//
// Empty means there is no response body, and an empty response code type
// must be used in its place.
//
// An empty string means that no special handling is required and the type
// response type can be used directly.
func responseKind(op *openapi3spec.Operation, code string) string {
	r := op.Responses[code]
	if len(r.Headers) != 0 {
		return "wrapped"
	}

	// Return here since there's no point continuing if we can't find bodies to
	// collide with
	if len(r.Content) == 0 {
		return "empty"
	}

	json := r.Content["application/json"]
	if json == nil {
		return ""
	}

	body := r.Content["application/json"].Schema

	for respCode, resp := range op.Responses {
		if respCode == code {
			continue
		}

		// Don't compare bodies if the other one doesn't have one
		if len(resp.Content) == 0 {
			continue
		}

		otherBody := resp.Content["application/json"].Schema
		if len(body.Ref) != 0 && len(otherBody.Ref) != 0 && body.Ref == otherBody.Ref {
			return "wrapped"
		} else if len(body.Ref) == 0 && len(otherBody.Ref) == 0 && body.Type == otherBody.Type {
			if isInlinePrimitive(body.Schema) {
				return "wrapped"
			}
		}
	}

	return ""
}

func responseTypeName(op *openapi3spec.Operation, code string, ignoreWrap bool) string {
	opName := strings.Title(op.OperationID) //nolint:staticcheck

	if !ignoreWrap && responseNeedsWrap(op, code) {
		return opName + "WrappedResponse"
	}

	resp := op.Responses[code]
	if len(resp.Content) == 0 {
		n, err := strconv.Atoi(code)
		if err != nil {
			panic("failed to convert status to int")
		}
		return "HTTPStatus" + xstrings.ToCamelCase(templates.HTTPStatus(n))
	}

	for _, schema := range resp.Content {
		if len(schema.Schema.Ref) != 0 {
			return strings.Title(templates.RefName(schema.Schema.Ref)) //nolint:staticcheck
		} else {
			if len(resp.Ref) != 0 {
				return templates.RefName(resp.Ref) + "Inline"
			}
			return opName + strings.Title(code) + "Inline" //nolint:staticcheck
		}
	}

	return "UNKNOWN"
}

func responseRefName(op *openapi3spec.Operation) string {
	if len(op.Responses) > 1 {
		// We need the interface for > 1
		opName := strings.Title(op.OperationID) //nolint:staticcheck
		return opName + "Response"
	}

	for code := range op.Responses {
		return responseTypeName(op, code, false)
	}

	return ""
}

func responseNeedsPtr(op *openapi3spec.Operation) bool {
	if len(op.Responses) > 1 {
		return false
	}

	for _, resp := range op.Responses {
		// For non-json payloads we don't want a pointer
		if _, ok := resp.Content["application/json"]; !ok && len(resp.Content) > 0 {
			return false
		}
	}

	return true
}

func paramTypeName(tdata templates.TemplateData, operationID string, methodName string, param openapi3spec.ParameterRef) (string, error) {
	if len(param.Schema.Ref) > 0 {
		return templates.RefName(param.Schema.Ref), nil
	} else if paramRequiresType(param) {
		return paramSchemaName(operationID, methodName, param.Name), nil
	} else {
		return primitive(tdata, param.Schema.Schema)
	}
}

func paramSchemaName(operationID string, methodName string, paramName string) string {
	return snakeToCamel(strings.Title(operationID)) + strings.Title(strings.ToLower(methodName)) + strings.Title(snakeToCamel(paramName)) + "Param" //nolint:staticcheck
}

func paramRequiresType(param openapi3spec.ParameterRef) bool {
	return len(param.Schema.Schema.Enum) != 0 || param.Schema.Schema.Type == "object" || param.Schema.Schema.Type == "array"
}

// paramConvertFn returns a function that will be able to convert between
// two types given the situation
func paramConvertFn(tdata templates.TemplateData, param openapi3spec.ParameterRef, paramTypeName, rhs string) (string, error) {
	outerType := param.Schema.Type
	innerType := param.Schema.Type
	innerFormat := ""
	if innerType == "array" {
		innerType = param.Schema.Items.Type
		if param.Schema.Items.Format != nil {
			innerFormat = *param.Schema.Items.Format
		}
	} else if param.Schema.Format != nil {
		innerFormat = *param.Schema.Format
	}

	var innerConversion string

	switch innerType {
	case "string":
		tdata.Import("github.com/aarondl/oa3/support")

		switch innerFormat {
		case "date":
			if tdata.TemplateParamEquals("timetype", "chrono") {
				innerConversion = "support.StringToChronoDate"
			} else {
				innerConversion = "support.StringToDate"
			}
		case "date-time":
			if tdata.TemplateParamEquals("timetype", "chrono") {
				innerConversion = "support.StringToChronoDateTime"
			} else {
				innerConversion = "support.StringToDateTime"
			}
		case "time":
			if tdata.TemplateParamEquals("timetype", "chrono") {
				innerConversion = "support.StringToChronoTime"
			} else {
				innerConversion = "support.StringToTime"
			}
		case "uuid":
			if tdata.TemplateParamEquals("uuidtype", "google") {
				innerConversion = "support.StringToUUID"
			}
		case "decimal":
			if tdata.TemplateParamEquals("decimaltype", "shopspring") {
				innerConversion = "support.StringToDecimal"
			}
		case "duration":
			innerConversion = "support.StringToDuration"
		case "":
			if param.Schema.Items != nil && len(param.Schema.Items.Enum) > 0 {
				innerConversion = fmt.Sprintf("support.StringToString[string, %s]", paramTypeName+"Item")
			} else {
				innerConversion = "support.StringNoOp"
			}
		default:
			return "", fmt.Errorf("no conversion function available for %s", param.Name)
		}
	case "boolean":
		tdata.Import("github.com/aarondl/oa3/support")
		innerConversion = "support.StringToBool"
	case "integer":
		tdata.Import("github.com/aarondl/oa3/support")

		switch innerFormat {
		case "int32":
			innerConversion = "support.StringToInt[int32]"
		case "int64":
			innerConversion = "support.StringToInt[int64]"
		case "", "int":
			innerConversion = "support.StringToInt[int]"
		default:
			return "", fmt.Errorf("no conversion function available for %s", param.Name)
		}
	case "number":
		tdata.Import("github.com/aarondl/oa3/support")
		switch innerFormat {
		case "float":
			innerConversion = "support.StringToFloat[float32]"
		case "", "double":
			innerConversion = "support.StringToFloat[float64]"
		default:
			return "", fmt.Errorf("no conversion function available for %s", param.Name)
		}
	}

	if innerType == outerType {
		return fmt.Sprintf("%s(%s[0])", innerConversion, rhs), nil
	}

	var prim string
	if param.Schema.Items != nil && len(param.Schema.Items.Enum) > 0 {
		prim = paramTypeName + "Item"
	} else {
		var err error
		prim, err = primitive(tdata, param.Schema.Schema.Items.Schema)
		if err != nil {
			return "", fmt.Errorf("failed to get primitive for param (%s): %w", param.Name, err)
		}
	}

	var outerConversion string
	switch {
	case *param.Style == "form" && *param.Explode && outerType == "array":
		tdata.Import("github.com/aarondl/oa3/support")
		outerConversion = fmt.Sprintf("support.ExplodedFormArrayToSlice[%s]", prim)
	case *param.Style == "form" && !*param.Explode && outerType == "array":
		tdata.Import("github.com/aarondl/oa3/support")
		outerConversion = fmt.Sprintf("support.FlatFormArrayToSlice[%s]", prim)
	// case *param.Style == "form" && *param.Explode && param.Schema.Type == "object":
	// case *param.Style == "form" && !*param.Explode && param.Schema.Type == "object":

	case *param.Style == "simple" && *param.Explode && outerType == "array":
		tdata.Import("github.com/aarondl/oa3/support")
		outerConversion = fmt.Sprintf("support.FlatFormArrayToSlice[%s]", prim)
	case *param.Style == "simple" && !*param.Explode && outerType == "array":
		tdata.Import("github.com/aarondl/oa3/support")
		outerConversion = fmt.Sprintf("support.FlatFormArrayToSlice[%s]", prim)
	}

	return fmt.Sprintf("%s(%s, %s)", outerConversion, rhs, innerConversion), nil
}

func hasComplexServers(servers []openapi3spec.Server) bool {
	complicated := false
	for _, s := range servers {
		if len(s.Variables) > 0 {
			complicated = true
			break
		}
	}

	return len(servers) > 1 && complicated
}

func hasJSONResponse(op openapi3spec.Operation) bool {
	for _, r := range op.Responses {
		if _, ok := r.Content["application/json"]; ok {
			return true
		}
	}

	return false
}
