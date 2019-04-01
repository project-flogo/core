package metadata

import (
	"encoding/json"
	"reflect"
	"strings"

	"github.com/project-flogo/core/data"
)

//
//type Metadata struct {
//	Settings map[string]typed.TypedValue
//	Input    map[string]typed.TypedValue
//	Output   map[string]typed.TypedValue
//}
//
//func New(mdStructs ...interface{}) *Metadata {
//
//	var settings map[string]typed.TypedValue
//	var input map[string]typed.TypedValue
//	var output map[string]typed.TypedValue
//
//	for _, mdStruct := range mdStructs {
//		typedMap := StructToTypedMap(mdStruct)
//		name := GetStructName(mdStruct)
//
//		switch name {
//		case "settings":
//			settings = typedMap
//		case "input":
//			input = typedMap
//		case "output":
//			output = typedMap
//		}
//	}
//
//	return &Metadata{Settings: settings, Input: input, Output: output}
//}

func GetStructName(mdStruct interface{}) string {
	if t := reflect.TypeOf(mdStruct); t.Kind() == reflect.Ptr {
		return strings.ToLower(t.Elem().Name())
	} else {
		return strings.ToLower(t.Name())
	}
}

type IOMetadata struct {
	Input  map[string]data.TypedValue
	Output map[string]data.TypedValue
}

func (md *IOMetadata) MarshalJSON() ([]byte, error) {
	var mdInputs []*data.Attribute
	var mdOutputs []*data.Attribute

	for _, v := range md.Input {
		if attr, ok := v.(*data.Attribute); ok {
			mdInputs = append(mdInputs, attr)
		}
	}
	for _, v := range md.Output {
		if attr, ok := v.(*data.Attribute); ok {
			mdOutputs = append(mdOutputs, attr)
		}
	}

	return json.Marshal(&struct {
		Input  []*data.Attribute `json:"input,omitempty"`
		Output []*data.Attribute `json:"output,omitempty"`
	}{
		Input:  mdInputs,
		Output: mdOutputs,
	})
}

func (md *IOMetadata) UnmarshalJSON(b []byte) error {

	ser := &struct {
		Input  []*data.Attribute `json:"input"`
		Output []*data.Attribute `json:"output"`
	}{}

	if err := json.Unmarshal(b, ser); err != nil {
		return err
	}

	md.Input = make(map[string]data.TypedValue, len(ser.Input))
	md.Output = make(map[string]data.TypedValue, len(ser.Output))

	for _, attr := range ser.Input {
		md.Input[attr.Name()] = attr
	}

	for _, attr := range ser.Output {
		md.Output[attr.Name()] = attr
	}

	return nil
}
