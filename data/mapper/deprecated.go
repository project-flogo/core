package mapper

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/project-flogo/core/data/coerce"
	"github.com/project-flogo/core/data/expression"
	"github.com/project-flogo/core/data/resolve"
	"sort"
	"strconv"
	"strings"
	"text/scanner"
	"unicode"
)

var testExprfactory = expression.NewFactory(resolve.GetBasicResolver())

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
	return convertMapperValue(mapping.Value, strType)
}

func ToObjectMapping(mappings []*objectMappings) (interface{}, error) {

	sort.Slice(mappings, func(i, j int) bool {
		return len(mappings[i].paths) < len(mappings[j].paths)
	})

	//obj := make(map[string]interface{})
	//toObjectFromPath(obj)
	return nil, nil
}

func toObjectFromPath(fields []string, value interface{}, object map[string]interface{}) interface{} {
	path := fields[0]
	fieldName := getFieldName(path)
	if strings.Index(path, "[") >= 0 && strings.HasSuffix(path, "]") {
		//Make sure the index are integer
		index, err := strconv.Atoi(getNameInsideBrancket(path))
		if err == nil {
			//Array
			if array, exist := object[fieldName]; !exist {
				tmpArrray := make([]interface{}, index+1)
				tmpArrray[index] = make(map[string]interface{})
				if len(fields) == 1 {
					tmpArrray[index] = value
					object[fieldName] = tmpArrray
					return object
				} else {
					object[fieldName] = tmpArrray
					return toObjectFromPath(fields[1:], value, tmpArrray[index].(map[string]interface{}))
				}
			} else {
				//exist
				if av, ok := array.([]interface{}); ok {
					if len(fields) <= 1 {
						av = Insert(av, index, value)
						object[fieldName] = av
						return object
					} else {
						if index >= len(av) {
							av = Insert(av, index, make(map[string]interface{}))
						} else {
							av[index] = make(map[string]interface{})
						}
						object[fieldName] = av
						return toObjectFromPath(fields[1:], value, object[fieldName].([]interface{})[index].(map[string]interface{}))
					}
				} else {
					return fmt.Errorf("not an array")
				}
			}

		} else {
			//Not an array
			if len(fields) == 1 {
				object[fieldName] = value
				return object
			} else {
				if _, exist := object[fieldName]; !exist {
					object[fieldName] = map[string]interface{}{}
				}
			}

			if len(fields) > 1 {
				toObjectFromPath(fields[1:], value, object[fieldName].(map[string]interface{}))
			}
		}

	} else {
		//Not an array
		if len(fields) <= 1 {
			object[fieldName] = value
			return object
		} else {
			if _, exist := object[fieldName]; !exist {
				object[fieldName] = map[string]interface{}{}
			}
		}

		toObjectFromPath(fields[1:], value, object[path].(map[string]interface{}))
	}
	return nil
}

func Insert(slice []interface{}, index int, value interface{}) []interface{} {
	if index >= len(slice) {
		// add to the end of slice in case of index >= len(slice)
		tmpArray := make([]interface{}, index+1)
		tmpArray[index] = value
		copy(tmpArray, slice)
		return tmpArray
	}
	slice[index] = value
	return slice
}

func getNameInsideBrancket(fieldName string) string {
	if strings.Index(fieldName, "[") >= 0 {
		index := fieldName[strings.Index(fieldName, "[")+1 : strings.Index(fieldName, "]")]
		return index
	}

	return ""
}

type objectMappings struct {
	fieldName string
	paths     []string
	mapping   *LegacyMapping
}

func getSortedMapto(mappings []*LegacyMapping) (map[string]interface{}, error) {
	input := make(map[string]interface{})
	fieldNameMap := make(map[string][]*objectMappings)
	for _, m := range mappings {
		field, err := ParseMappingField(m.MapTo)
		if err != nil {
			return nil, err
		}
		objectLevelFields := field.Getfields()
		fieldName := getFieldName(objectLevelFields[0])
		fieldNameMap[fieldName] = append(fieldNameMap[fieldName], &objectMappings{fieldName: fieldName, paths: objectLevelFields[1:], mapping: m})
	}

	for k, v := range fieldNameMap {
		if len(v) == 1 {
			val, err := convertMappingValue(v[0].mapping)
			if err != nil {
				return nil, err
			}
			input[k] = val
		}

		if len(v) > 1 {

		}
	}
	return nil, nil
}

func getFieldName(fieldname string) string {
	if strings.Index(fieldname, "[") > 0 && strings.Index(fieldname, "]") > 0 {
		return fieldname[:strings.Index(fieldname, "[")]
	}
	return fieldname
}

type LegacyArrayMapping struct {
	From   interface{}           `json:"from"`
	To     string                `json:"to"`
	Type   string                `json:"type"`
	Fields []*LegacyArrayMapping `json:"fields,omitempty"`
}

func ParseArrayMapping(arrayDatadata interface{}) (*LegacyArrayMapping, error) {
	amapping := &LegacyArrayMapping{}
	switch t := arrayDatadata.(type) {
	case string:
		err := json.Unmarshal([]byte(t), amapping)
		if err != nil {
			return nil, err
		}
	case interface{}:
		s, err := coerce.ToString(t)
		if err != nil {
			return nil, fmt.Errorf("Convert array mapping value to string error, due to [%s]", err.Error())
		}
		err = json.Unmarshal([]byte(s), amapping)
		if err != nil {
			return nil, err
		}
	}
	return amapping, nil
}

func ToNewArray(mapping *LegacyArrayMapping) (interface{}, error) {
	var newMapping interface{}
	var fieldsMapping map[string]interface{}
	if mapping.From == "NEWARRAY" {
		fieldsMappings := make([]interface{}, 1)
		fieldsMappings[0] = make(map[string]interface{})
		fieldsMapping = fieldsMappings[0].(map[string]interface{})
		newMapping = fieldsMappings
	} else {
		newMapping = make(map[string]interface{})
		fieldsMapping = make(map[string]interface{})
		newMapping.(map[string]interface{})[fmt.Sprintf("@foreach(%s)", mapping.From)] = fieldsMapping
	}

	var err error
	for _, field := range mapping.Fields {
		if field.Type == "foreach" {
			//Check to see if it is a new array
			fieldsMapping[ToNewMapto(field.To)], err = ToNewArray(field)
		} else {
			fieldsMapping[ToNewMapto(field.To)], err = convertMapperValue(field.From, field.Type)
			if err != nil {
				return nil, err
			}
		}
	}
	return newMapping, nil
}

func ToNewMapto(old string) string {
	if strings.HasPrefix(old, "$.") || strings.HasPrefix(old, "$$") {
		return old[2:]
	}
	return old
}

//func ToMaptoObject(mapto string) (interface{}, error) {
//	m := make(map[string]interface{})
//	err := path.SetValue(m, mapto, m)
//	if err != nil {
//		return nil, fmt.Errorf("construct object from map to [%s] failed, due to [%s]", mapto, err.Error())
//	}
//	return m, nil
//}

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

func convertMapperValue(value interface{}, typ string) (interface{}, error) {
	switch typ {
	case "assign", "1":
		if v, ok := value.(string); ok {
			if !ResolvableExpr(v) {
				return v, nil
			}
			return "=" + v, nil
		}
		return value, nil
	case "literal", "2":
		return value, nil
	case "expression", "3":
		expr, ok := value.(string)
		if !ok {
			return nil, fmt.Errorf("invalid expression mapping: '%v'", value)
		}
		return "=" + expr, nil
	case "object", "4":
		return value, nil
	case "array", "5":
		arrayMapping, err := ParseArrayMapping(value)
		if err != nil {
			return nil, err
		}
		return ToNewArray(arrayMapping)
	case "primitive":
		//This use to handle very old array mapping type
		if priValue, ok := value.(string); ok {
			if !ResolvableExpr(priValue) {
				//Not an expr, just return is as value
				return priValue, nil
			}
			return "=" + priValue, nil
		}

		return value, nil
	default:
		return 0, errors.New("unsupported mapping type: " + typ)
	}
}

func ResolvableExpr(expr string) bool {
	_, err := testExprfactory.NewExpr(expr)
	if err != nil {
		//Not an expr, just return is as value
		return false
	}
	return true
}

type MappingField struct {
	fields []string
	ref    string
	s      *scanner.Scanner
}

func NewMappingField(fields []string) *MappingField {
	return &MappingField{fields: fields}
}

func ParseMappingField(mRef string) (*MappingField, error) {
	//Remove any . at beginning
	if strings.HasPrefix(mRef, ".") {
		mRef = mRef[1:]
	}
	g := &MappingField{ref: mRef}

	err := g.Start(mRef)
	if err != nil {
		return nil, fmt.Errorf("parse mapping [%s] failed, due to %s", mRef, err.Error())
	}
	return g, nil
}

func (m *MappingField) GetRef() string {
	return m.ref
}

func (m *MappingField) Getfields() []string {
	return m.fields
}

func (m *MappingField) paserName() error {
	fieldName := ""
	switch ch := m.s.Scan(); ch {
	case '.':
		return m.Parser()
	case '[':
		//Done
		if fieldName != "" {
			m.fields = append(m.fields, fieldName)
		}
		m.s.Mode = scanner.ScanInts
		nextAfterBracket := m.s.Scan()
		if nextAfterBracket == '"' || nextAfterBracket == '\'' {
			//Special charactors
			m.s.Mode = scanner.ScanIdents
			return m.handleSpecialField(nextAfterBracket)
		} else {
			//HandleArray
			if m.fields == nil || len(m.fields) <= 0 {
				m.fields = append(m.fields, "["+m.s.TokenText()+"]")
			} else {
				m.fields[len(m.fields)-1] = m.fields[len(m.fields)-1] + "[" + m.s.TokenText() + "]"
			}
			ch := m.s.Scan()
			if ch != ']' {
				return fmt.Errorf("Inliad array format")
			}
			m.s.Mode = scanner.ScanIdents
			return m.Parser()
		}
	case scanner.EOF:
		if fieldName != "" {
			m.fields = append(m.fields, fieldName)
		}
	default:
		fieldName = fieldName + m.s.TokenText()
		if fieldName != "" {
			m.fields = append(m.fields, fieldName)
		}
		return m.Parser()
	}

	return nil
}

func (m *MappingField) handleSpecialField(startQutoes int32) error {
	fieldName := ""
	run := true

	for run {
		switch ch := m.s.Scan(); ch {
		case startQutoes:
			//Check if end with startQutoes]
			nextAfterQuotes := m.s.Scan()
			if nextAfterQuotes == ']' {
				//end specialfield, startover
				m.fields = append(m.fields, fieldName)
				run = false
				return m.Parser()
			} else {
				fieldName = fieldName + string(startQutoes)
				fieldName = fieldName + m.s.TokenText()
			}
		default:
			fieldName = fieldName + m.s.TokenText()
		}
	}
	return nil
}

func (m *MappingField) Parser() error {
	switch ch := m.s.Scan(); ch {
	case '.':
		return m.paserName()
	case '[':
		m.s.Mode = scanner.ScanInts
		nextAfterBracket := m.s.Scan()
		if nextAfterBracket == '"' || nextAfterBracket == '\'' {
			//Special charactors
			m.s.Mode = scanner.ScanIdents
			return m.handleSpecialField(nextAfterBracket)
		} else {
			//HandleArray
			if m.fields == nil || len(m.fields) <= 0 {
				m.fields = append(m.fields, "["+m.s.TokenText()+"]")
			} else {
				m.fields[len(m.fields)-1] = m.fields[len(m.fields)-1] + "[" + m.s.TokenText() + "]"
			}
			//m.handleArray()
			ch := m.s.Scan()
			if ch != ']' {
				return fmt.Errorf("Inliad array format")
			}
			m.s.Mode = scanner.ScanIdents
			return m.Parser()
		}
	case scanner.EOF:
		//Done
		return nil
	default:
		m.fields = append(m.fields, m.s.TokenText())
		return m.paserName()
	}
	return nil
}

func (m *MappingField) Start(jsonPath string) error {
	m.s = new(scanner.Scanner)
	m.s.IsIdentRune = IsIdentRune
	m.s.Init(strings.NewReader(jsonPath))
	m.s.Mode = scanner.ScanIdents
	//Donot skip space and new line
	m.s.Whitespace ^= 1<<'\t' | 1<<'\n' | 1<<'\r' | 1<<' '
	return m.Parser()
}

func IsIdentRune(ch rune, i int) bool {
	return ch == '$' || ch == '-' || ch == '_' || unicode.IsLetter(ch) || unicode.IsDigit(ch) && i > 0
}
