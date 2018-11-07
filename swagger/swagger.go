package swagger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

// Generate generates a Swagger 2.0 document based off of the provided endpoints.
func Generate(host string, name string, description string, version string, endpoints []Endpoint) ([]byte, error) {
	fmt.Println("Inside Generate")
	fmt.Println("endpoints received:", endpoints)
	paths := map[string]interface{}{}

	for _, endpoint := range endpoints {
		path := map[string]interface{}{}
		fmt.Println("Before scrubbedpath")
		parameters, scrubbedPath := swaggerParametersExtractor(endpoint.Path, endpoint.BeginDelim, endpoint.EndDelim)
		fmt.Println("After swag returns")
		ok := map[string]interface{}{
			"description": endpoint.Description,
		}
		fmt.Println("path map")
		path[strings.ToLower(endpoint.Method)] = map[string]interface{}{
			"description": endpoint.Description,
			"tags":        []interface{}{endpoint.Name},
			"parameters":  parameters,
			"responses": map[string]interface{}{
				"200": ok,
				"default": map[string]interface{}{
					"description": "error",
				},
			},
		}
		paths[scrubbedPath] = path
	}

	swagger := map[string]interface{}{
		"swagger": "2.0",
		"info": map[string]interface{}{
			"version":     version,
			"title":       name,
			"description": description,
		},
		"host":  host,
		"paths": paths,
	}
	fmt.Println("marshalllll")
	docs, err := json.MarshalIndent(&swagger, "", "    ")
	if err != nil {
		return nil, err
	}
	fmt.Println("before return generate")
	return docs, err
}

func swaggerParametersExtractor(path string, beginDelim rune, endDelim rune) ([]interface{}, string) {
	fmt.Println("Inside Swaggerparams")
	parameters := []interface{}{}
	routePath := []rune(path)
	for i := 0; i < len(routePath); i++ {
		fmt.Println("Inside for")
		if routePath[i] == beginDelim {
			key := bytes.Buffer{}
			for i++; i < len(routePath) && routePath[i] != endDelim; i++ {
				if routePath[i] != ' ' && routePath[i] != '\t' {
					key.WriteRune(routePath[i])
				}
			}
			if beginDelim == ':' {
				path = strings.Replace(path, fmt.Sprintf(":%s", key.String()), fmt.Sprintf("{%s}", key.String()), 1)
			}
			parameter := map[string]interface{}{
				"name":     key.String(),
				"in":       "path",
				"required": true,
				"type":     "string",
			}
			parameters = append(parameters, parameter)
		}
	}
	fmt.Println("Before swag returns")
	return parameters, path
}
