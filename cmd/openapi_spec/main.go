package main

import (
	"consoledot-go-template/internal/payloads"
	"encoding/json"
	"os"
	"reflect"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3gen"
	"gopkg.in/yaml.v3"
)

// addPayloads helps isolate all the schema types registration
func addPayloads(spec *APISpec) {
	// payloads - MAKE SURE THE TYPE HAS JSON/YAML Go STRUCT TAGS (or "map key XXX not found" error occurs)
	spec.addTypeSchema("v1.HelloRequest", &payloads.HelloRequest{})
	spec.addTypeSchema("v1.HelloResponse", &payloads.HelloResponse{})
}

func addErrors(spec *APISpec) {
	// error payloads
	spec.addTypeSchema("v1.ErrorResponse", &payloads.ErrorResponse{})

	// general error responses
	spec.addResponse("NotFound", "The requested resource was not found", "#/components/schemas/v1.ErrorResponse")
	spec.addResponse("InternalError", "The server encountered an internal error", "#/components/schemas/v1.ErrorResponse")
	spec.addResponse("BadRequest", "The request's parameters are invalid", "#/components/schemas/v1.ErrorResponse")
}

// Enables nullable fields in OpenAPI spec by go tag nullable: "true".
// Keep in mind, that this generates OpenAPI 3 type nullable,
// which is backward incompatible with Swagger (OpenAPI 2).
var enableNullableCustomizer = openapi3gen.SchemaCustomizer(
	func(_name string, _t reflect.Type, tag reflect.StructTag, schema *openapi3.Schema) error {
		if tag.Get("nullable") == "true" {
			schema.Nullable = true
		}
		return nil
	},
)

type APISpec struct {
	Components openapi3.Components `json:"components,omitempty" yaml:"components,omitempty"`
	Servers    openapi3.Servers    `json:"servers,omitempty" yaml:"servers,omitempty"`
}

func NewSpec() APISpec {
	spec := APISpec{}
	spec.Servers = openapi3.Servers{
		&openapi3.Server{
			Description: "Local development",
			URL:         "http://0.0.0.0:{port}/api/{applicationName}",
			Variables: map[string]*openapi3.ServerVariable{
				"applicationName": {Default: "template"},
				"port":            {Default: "8000"},
			},
		},
	}
	spec.Components = openapi3.NewComponents()
	spec.Components.Schemas = make(map[string]*openapi3.SchemaRef)
	spec.Components.Responses = make(map[string]*openapi3.ResponseRef)
	return spec
}

func (spec APISpec) addTypeSchema(name string, model interface{}) {
	schema, err := openapi3gen.NewSchemaRefForValue(model, spec.Components.Schemas, enableNullableCustomizer)
	if err != nil {
		panic(err)
	}
	spec.Components.Schemas[name] = schema
}

func (spec APISpec) addResponse(name string, description string, ref string) {
	response := openapi3.NewResponse().WithDescription(description).WithJSONSchemaRef(&openapi3.SchemaRef{Ref: ref})
	spec.Components.Responses[name] = &openapi3.ResponseRef{Value: response}
}

func main() {
	spec := NewSpec()
	addPayloads(&spec)
	addErrors(&spec)

	bufferYAML, err := os.ReadFile("./cmd/openapi_spec/paths.yml")
	if err != nil {
		panic(err)
	}

	schemas, err := yaml.Marshal(&spec)
	if err != nil {
		panic(err)
	}
	bufferYAML = append(bufferYAML, schemas...)

	// Load the final spec and dump it to JSON
	// We also validate the YAML spec by doing this.
	loadedSchema, err := openapi3.NewLoader().LoadFromData(bufferYAML)
	if err != nil {
		panic(err)
	}
	bufferJson, err := json.MarshalIndent(loadedSchema, "", "  ")
	if err != nil {
		panic(err)
	}

	if err = os.WriteFile("./api/openapi.gen.yml", bufferYAML, 0o644); err != nil {
		panic(err)
	}
	if err = os.WriteFile("./api/openapi.gen.json", bufferJson, 0o644); err != nil {
		panic(err)
	}
}
