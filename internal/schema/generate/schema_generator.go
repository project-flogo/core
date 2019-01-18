// +build ignore

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/square-it/jsonschema"
	"github.com/project-flogo/core/app"
)

func main() {
	reflector := &jsonschema.Reflector{ExpandedStruct: true}
	schema := reflector.Reflect(&app.Config{})
	schemaJSON, err := json.MarshalIndent(schema, "", "    ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	err = ioutil.WriteFile("schema.json", schemaJSON, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
