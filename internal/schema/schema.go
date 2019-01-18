//go:generate go run generate/schema_generator.go
//go:generate go-bindata -pkg schema -o assets.go schema.json

package schema

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/xeipuuv/gojsonschema"
)

var schema *gojsonschema.Schema

func init() {
	jsonSchema, err := Asset("schema.json")
	if err != nil {
		panic(err)
	}
	schemaLoader := gojsonschema.NewStringLoader(string(jsonSchema))
	schema, err = gojsonschema.NewSchema(schemaLoader)
	if err != nil {
		panic(err)
	}
}

// Validate validates the provided JSON against the v2 JSON schema.
func Validate(JSON []byte) error {
	JSONLoader := gojsonschema.NewStringLoader(string(JSON))
	result, err := schema.Validate(JSONLoader)

	if err != nil {
		return err
	}

	if result.Valid() {
		return err
	}
	var msg bytes.Buffer

	msg.WriteString("The JSON is not valid. See errors:\n")
	for _, desc := range result.Errors() {
		msg.WriteString(fmt.Sprintf("- %s\n", desc))
	}
	return errors.New(msg.String())
}
