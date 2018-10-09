package mapper

import (
	"fmt"
)

//Deprecated
func ConvertMappingValue(mappingType string, mappingValue interface{}) (interface{}, error) {
	switch mappingType {
	case "assign", "1":
		return mappingValue, nil
	case "literal", "2":
		return mappingValue, nil
	case "expression", "3":
		expr, ok := mappingValue.(string)
		if !ok {
			return nil, fmt.Errorf("invalid expression mapping: '%v'", mappingValue)
		}
		if expr[0] != '=' {
			return "=" + expr, nil
		}

		return expr, nil
	case "object", "4":
		return mappingValue, nil
	case "array", "5":
		return mappingValue, nil
	default:
		return 0, fmt.Errorf("unsupported mapping type: %s", mappingType)
	}
}
