package mapper

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
)

//DEPRECATED
type LegacyMappings struct {
	Input  []*LegacyMapping `json:"input,omitempty"`
	Output []*LegacyMapping `json:"output,omitempty"`
}

// MappingDef is a simple structure that defines a mapping
//DEPRECATED
type LegacyMapping struct {
	Type  interface{} `json:"type"`
	Value interface{} `json:"value"`
	MapTo string      `json:"mapTo"`
}

func ConvertLegacyMappings(mappings *LegacyMappings) (input map[string]interface{}, output map[string]interface{}, err error) {

	if mappings == nil {
		return nil, nil, nil
	}

	input = make(map[string]interface{}, len(mappings.Input))
	output = make(map[string]interface{}, len(mappings.Output))

	if mappings.Input != nil {
		for _, mapping := range mappings.Input {
			val, err := convertMappingValue(mapping)
			if err != nil {
				return nil, nil, err
			}
			input[mapping.MapTo] = val
		}
	}

	if mappings.Output != nil {
		for _, mapping := range mappings.Output {
			val, err := convertMappingValue(mapping)
			if err != nil {
				return nil, nil, err
			}
			output[mapping.MapTo] = val
		}
	}

	return input, output, nil
}

func convertMappingValue(mapping *LegacyMapping) (interface{}, error) {
	strType, _ := toString(mapping.Type)
	switch strType {
	case "assign", "1":
		return mapping.Value, nil
	case "literal", "2":
		return mapping.Value, nil
	case "expression", "3":
		expr, ok := mapping.Value.(string)
		if !ok {
			return nil, fmt.Errorf("invalid expression mapping: '%v'", mapping.Value)
		}
		return expr, nil
	case "object", "4":
		return mapping.Value, nil
	case "array", "5":
		return mapping.Value, nil
	default:
		return 0, errors.New("unsupported mapping type: " + strType)
	}
}

// ToString coerce a value to a string
func toString(val interface{}) (string, error) {

	switch t := val.(type) {
	case string:
		return t, nil
	case int:
		return strconv.Itoa(t), nil
	case float64:
		return strconv.FormatFloat(t, 'f', -1, 64), nil
	case json.Number:
		return t.String(), nil
	default:
		return "", nil
	}
}
