package activity

import (
	"encoding/json"

	"github.com/project-flogo/core/data"
	"github.com/project-flogo/core/data/metadata"
)

type Metadata struct {
	*metadata.IOMetadata
	Settings map[string]data.TypedValue
}

// UnmarshalJSON overrides the default UnmarshalJSON for TaskEnv
func (md *Metadata) UnmarshalJSON(b []byte) error {

	ser := &struct {
		Settings []*data.Attribute `json:"settings"`
		Input    []*data.Attribute `json:"input"`
		Output   []*data.Attribute `json:"output"`
	}{}

	if err := json.Unmarshal(b, ser); err != nil {
		return err
	}

	md.IOMetadata = &metadata.IOMetadata{}

	md.Settings = make(map[string]data.TypedValue, len(ser.Settings))
	md.Input = make(map[string]data.TypedValue, len(ser.Input))
	md.Output = make(map[string]data.TypedValue, len(ser.Output))

	for _, attr := range ser.Settings {
		md.Settings[attr.Name()] = attr
	}

	for _, attr := range ser.Input {
		md.Input[attr.Name()] = attr
	}

	for _, attr := range ser.Output {
		md.Output[attr.Name()] = attr
	}

	return nil
}

func ToMetadata(mdStructs ...interface{}) *Metadata {

	var settings map[string]data.TypedValue
	var input map[string]data.TypedValue
	var output map[string]data.TypedValue

	for _, mdStruct := range mdStructs {
		typedMap := metadata.StructToTypedMap(mdStruct)
		name := metadata.GetStructName(mdStruct)

		switch name {
		case "settings":
			settings = typedMap
		case "input":
			input = typedMap
		case "output":
			output = typedMap
		}
	}

	return &Metadata{Settings: settings, IOMetadata: &metadata.IOMetadata{Input: input, Output: output}}
}
