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

func ConvertLegacyMappings(mappings *LegacyMappings, resolver resolve.CompositeResolver) (input map[string]interface{}, output map[string]interface{}, err error) {

	if mappings == nil {
		return nil, nil, nil
	}

	input = make(map[string]interface{}, len(mappings.Input))
	output = make(map[string]interface{}, len(mappings.Output))

	if mappings.Input != nil {
		input, err = HandleMappings(mappings.Input, resolver)
		if err != nil {
			return nil, nil, err
		}
	}

	if mappings.Output != nil {
		output, err = HandleMappings(mappings.Output, resolver)
		if err != nil {
			return nil, nil, err
		}
	}

	return input, output, nil
}

func HandleMappings(mappings []*LegacyMapping, resolver resolve.CompositeResolver) (map[string]interface{}, error) {

	input := make(map[string]interface{})

	fieldNameMap := make(map[string][]*objectMappings)
	for _, m := range mappings {
		//target is single field name
		if strings.Index(m.MapTo, ".") <= 0 && (strings.Index(m.MapTo, "[") <= 0 || strings.Index(m.MapTo, "]") <= 0) {
			typ, _ := toString(m.Type)
			val, err := convertMapperValue(m.Value, typ, resolver)
			if err != nil {
				return nil, err
			}
			input[m.MapTo] = val
		} else {
			//Handle multiple value to single field
			field, err := ParseMappingField(m.MapTo)
			if err != nil {
				return nil, err
			}
			mapToFields := field.Getfields()
			fieldName := getFieldName(mapToFields[0])
			objMapping := &objectMappings{fieldName: fieldName, mapping: m}

			if strings.Index(mapToFields[0], "[") >= 0 && strings.Index(mapToFields[0], "]") > 0 {
				mapToFields[0] = mapToFields[0][len(fieldName):]
				objMapping.targetFields = mapToFields
			} else {
				if len(mapToFields) > 1 {
					objMapping.targetFields = mapToFields[1:]
				}
			}
			fieldNameMap[fieldName] = append(fieldNameMap[fieldName], objMapping)
		}
	}

	for k, v := range fieldNameMap {
		sort.Slice(v, func(i, j int) bool {
			return len(v[i].targetFields) < len(v[j].targetFields)
		})

		var obj interface{}
		for _, objMapping := range v {
			typ, _ := toString(objMapping.mapping.Type)
			val, err := convertMapperValue(objMapping.mapping.Value, typ, resolver)
			if err != nil {
				return nil, err
			}

			if obj == nil && len(objMapping.targetFields) > 0 {
				if strings.Index(objMapping.targetFields[0], "[") >= 0 && strings.Index(objMapping.targetFields[0], "]") > 0 {
					obj = make([]interface{}, 1)
				} else {
					obj = make(map[string]interface{})
				}
			}

			if len(objMapping.targetFields) > 0 {
				obj, err = constructObjectFromPath(objMapping.targetFields, val, obj)
				if err != nil {
					return nil, err
				}
			} else {
				obj = val
			}
		}
		input[k] = obj
	}

	return input, nil
}

func (o *objectMappings) constructObject(value interface{}) (interface{}, error) {
	var obj interface{}
	if strings.Index(o.targetFields[0], "[") >= 0 && strings.Index(o.targetFields[0], "]") > 0 {
		obj = make([]interface{}, 1)
	} else {
		obj = make(map[string]interface{})
	}
	var err error
	obj, err = constructObjectFromPath(o.targetFields, value, obj)
	if err != nil {
		return nil, err
	}
	return obj, nil
}

func createArray(index int) []interface{} {
	tmpArrray := make([]interface{}, index+1)
	tmpArrray[index] = make(map[string]interface{})
	return tmpArrray
}

func constructObjectFromPath(fields []string, value interface{}, object interface{}) (interface{}, error) {
	fieldName := getFieldName(fields[0])
	//Has array
	if strings.Index(fields[0], "[") >= 0 && strings.HasSuffix(fields[0], "]") {
		//Make sure the index is integer
		index, err := strconv.Atoi(getNameInsideBrancket(fields[0]))
		if err == nil {
			object, err = handlePathArray(fieldName, index, fields, value, object)
			if err != nil {
				return nil, err
			}
		} else {
			if err := handlePathObject(fieldName, fields, value, object); err != nil {
				return nil, err
			}
		}
	} else {
		if err := handlePathObject(fieldName, fields, value, object); err != nil {
			return nil, err
		}
	}

	return object, nil
}

func handlePathObject(fieldName string, fields []string, value interface{}, object interface{}) error {
	//Not an array
	if len(fields) <= 1 {
		object.(map[string]interface{})[fieldName] = value
		return nil
	} else {
		if _, exist := object.(map[string]interface{})[fieldName]; !exist {
			object.(map[string]interface{})[fieldName] = map[string]interface{}{}
		}
	}
	var err error
	object.(map[string]interface{})[fieldName], err = constructObjectFromPath(fields[1:], value, object.(map[string]interface{})[fieldName].(map[string]interface{}))
	if err != nil {
		return err
	}
	return nil
}

func handlePathArray(fieldName string, index int, fields []string, value interface{}, object interface{}) (interface{}, error) {
	//Only root field with array no field name
	var err error
	if fieldName == "" {
		_, ok := object.([]interface{})
		if !ok {
			tmpArrray := createArray(index)
			object = tmpArrray
		} else {
			if index >= len(object.([]interface{})) {
				object = Insert(object.([]interface{}), index, map[string]interface{}{})
			}
		}
		if len(fields) == 1 {
			object.([]interface{})[index] = value
			return object, nil
		} else {
			if object.([]interface{})[index] == nil {
				object.([]interface{})[index] = make(map[string]interface{})
			}
			object.([]interface{})[index], err = constructObjectFromPath(fields[1:], value, object.([]interface{})[index])
			if err != nil {
				return nil, err
			}
		}
	} else {
		obj := object.(map[string]interface{})
		if array, exist := obj[fieldName]; !exist {
			obj[fieldName] = createArray(index)
			if len(fields) == 1 {
				obj[fieldName].([]interface{})[index] = value
				return obj, nil
			} else {
				obj[fieldName].([]interface{})[index], err = constructObjectFromPath(fields[1:], value, obj[fieldName].([]interface{})[index].(map[string]interface{}))
				if err != nil {
					return nil, err
				}
			}
		} else {
			//exist
			if av, ok := array.([]interface{}); ok {
				if len(fields) <= 1 {
					av = Insert(av, index, value)
					obj[fieldName] = av
					return obj, nil
				} else {
					if index >= len(av) {
						av = Insert(av, index, make(map[string]interface{}))
					} else {
						av[index] = make(map[string]interface{})
					}
					obj[fieldName] = av
					obj[fieldName].([]interface{})[index], err = constructObjectFromPath(fields[1:], value, obj[fieldName].([]interface{})[index].(map[string]interface{}))
					if err != nil {
						return nil, err
					}
				}
			} else {
				return nil, fmt.Errorf("not an array")
			}
		}
	}
	return object, nil
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
	fieldName    string
	targetFields []string
	mapping      *LegacyMapping
}

func getFieldName(fieldname string) string {
	if strings.Index(fieldname, "[") >= 0 && strings.Index(fieldname, "]") > 0 {
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

func ToNewArray(mapping *LegacyArrayMapping, resolver resolve.CompositeResolver) (interface{}, error) {
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
			fieldsMapping[ToNewArrayChildMapto(field.To)], err = ToNewArray(field, resolver)
		} else {
			fieldsMapping[ToNewArrayChildMapto(field.To)], err = convertMapperValue(field.From, field.Type, resolver)
			if err != nil {
				return nil, err
			}
		}
	}
	return newMapping, nil
}

func ToNewArrayChildMapto(old string) string {
	if strings.HasPrefix(old, "$.") || strings.HasPrefix(old, "$$") {
		return old[2:]
	}
	return old
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

func convertMapperValue(value interface{}, typ string, resolver resolve.CompositeResolver) (interface{}, error) {
	switch typ {
	case "assign", "1":
		if v, ok := value.(string); ok {
			if !ResolvableExpr(v, resolver) {
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
			return value, nil
		}
		return "=" + expr, nil
	case "object", "4":
		return value, nil
	case "array", "5":
		arrayMapping, err := ParseArrayMapping(value)
		if err != nil {
			return nil, err
		}
		return ToNewArray(arrayMapping, resolver)
	case "primitive":
		//This use to handle very old array mapping type
		if priValue, ok := value.(string); ok {
			if !ResolvableExpr(priValue, resolver) {
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

func ResolvableExpr(expr string, resolver resolve.CompositeResolver) bool {
	_, err := expression.NewFactory(resolver).NewExpr(expr)
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
