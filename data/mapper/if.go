package mapper

type IfElseMapping struct {
	Mapping interface{} `json:"mapping"`
}

func GetIfElseMapping(value interface{}) (interface{}, bool) {
	switch t := value.(type) {
	case *IfElseMapping:
		return t.Mapping, true
	case map[string]interface{}:
		if mapping, ok := t["@if"]; ok {
			return mapping, true
		}
		return nil, false
	default:
		return nil, false
	}
}
