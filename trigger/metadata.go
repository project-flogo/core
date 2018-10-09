package trigger

import (
	"github.com/project-flogo/core/data"
	"github.com/project-flogo/core/data/metadata"
	"reflect"
	"strings"
)

type Metadata struct {
	Settings        map[string]data.TypedValue
	HandlerSettings map[string]data.TypedValue
	Output          map[string]data.TypedValue
	Reply           map[string]data.TypedValue

	//DEPRECATED
	ID string
}

func NewMetadata(mdStructs ...interface{}) *Metadata {

	if len(mdStructs) == 1 {
		if mdJson, ok := mdStructs[0].(string); ok {
			return metadataFromOldJson(mdJson)
		}
	}

	var settings map[string]data.TypedValue
	var handlerSettings map[string]data.TypedValue
	var output map[string]data.TypedValue
	var reply map[string]data.TypedValue

	for _, mdStruct := range mdStructs {
		typedMap := metadata.StructToTypedMap(mdStruct)
		name := getName(mdStruct)

		switch name {
		case "settings":
			settings = typedMap
		case "handlersettings":
			handlerSettings = typedMap
		case "output":
			output = typedMap
		case "reply":
			reply = typedMap
		}
	}

	return &Metadata{Settings: settings, HandlerSettings: handlerSettings, Output: output, Reply: reply}
}

func getName(mdStruct interface{}) string {
	if t := reflect.TypeOf(mdStruct); t.Kind() == reflect.Ptr {
		return strings.ToLower(t.Elem().Name())
	} else {
		return strings.ToLower(t.Name())
	}
}

//func FromStructMetadata(md *StructMetadata) *Metadata {
//
//	settings := metadata.StructToTypedMap(md.Settings)
//	handlerSettings := metadata.StructToTypedMap(md.HandlerSettings)
//	reply := metadata.StructToTypedMap(md.Reply)
//	outputs := metadata.StructToTypedMap(md.Outputs)
//
//	return &Metadata{Settings: settings, HandlerSettings: handlerSettings, Reply: reply, Outputs: outputs}
//}

//// Metadata is the metadata for a Trigger
//type MetadataOld struct {
//	ID       string
//	Version  string
//	Handler  *HandlerMetadata
//	Settings map[string]*data.Attribute
//	Output   map[string]*data.Attribute
//	Reply    map[string]*data.Attribute
//}
//
//// EndpointMetadata is the metadata for a Trigger Endpoint
//type HandlerMetadata struct {
//	Settings []*data.Attribute `json:"settings"`
//}
//
//// NewMetadata creates a Metadata object from the json representation
////todo should return error instead of panic
//func NewMetadata(jsonMetadata string) *Metadata {
//	md := &Metadata{}
//	err := json.Unmarshal([]byte(jsonMetadata), md)
//	if err != nil {
//		panic("Unable to parse trigger metadata: " + err.Error())
//	}
//
//	return md
//}
//
//// UnmarshalJSON overrides the default UnmarshalJSON for Metadata
//func (md *Metadata) UnmarshalJSON(b []byte) error {
//
//	ser := &struct {
//		Name     string            `json:"name"`
//		Version  string            `json:"version"`
//		Ref      string            `json:"ref"`
//		Handler  *HandlerMetadata  `json:"handler"`
//		Settings []*data.Attribute `json:"settings"`
//		Output   []*data.Attribute `json:"output"`
//		Reply    []*data.Attribute `json:"reply"`
//
//		//for backwards compatibility
//		Endpoint *HandlerMetadata  `json:"endpoint"`
//		Outputs  []*data.Attribute `json:"outputs"`
//	}{}
//
//	if err := json.Unmarshal(b, ser); err != nil {
//		return err
//	}
//
//	if len(ser.Ref) > 0 {
//		md.ID = ser.Ref
//	} else {
//		// Added for backwards compatibility
//		// TODO remove and add a proper error once the BC is removed
//		md.ID = ser.Name
//	}
//
//	md.Version = ser.Version
//
//	if ser.Handler != nil {
//		md.Handler = ser.Handler
//	} else {
//		// Added for backwards compatibility
//		// TODO remove and add a proper error once the BC is removed
//
//		if ser.Endpoint != nil {
//			md.Handler = ser.Endpoint
//		}
//	}
//
//	md.Settings = make(map[string]*data.Attribute, len(ser.Settings))
//	md.Output = make(map[string]*data.Attribute, len(ser.Outputs))
//	md.Reply = make(map[string]*data.Attribute, len(ser.Reply))
//
//	for _, attr := range ser.Settings {
//		md.Settings[attr.Name()] = attr
//	}
//
//	if len(ser.Output) > 0 {
//		for _, attr := range ser.Output {
//			md.Output[attr.Name()] = attr
//		}
//	} else {
//		// for backwards compatibility
//		for _, attr := range ser.Outputs {
//			md.Output[attr.Name()] = attr
//		}
//	}
//
//	for _, attr := range ser.Reply {
//		md.Reply[attr.Name()] = attr
//	}
//
//	return nil
//}
//
//// OutputsToAttrs converts the supplied output data to attributes
//// todo remove coerce option, coercion now happens by default
//func (md *Metadata) OutputsToAttrs(outputData map[string]interface{}, coerce bool) ([]*data.Attribute, error) {
//
//	attrs := make([]*data.Attribute, 0, len(md.Output))
//
//	for k, a := range md.Output {
//		v, _ := outputData[k]
//
//		//if coerce {
//		//	var err error
//		//	v, err = data.CoerceToValue(v, a.Type)
//		//
//		//	if err != nil {
//		//		return nil, err
//		//	}
//		//}
//
//		var err error
//		attr, err := data.NewAttribute(a.Name(), a.Type(), v)
//		attrs = append(attrs, attr)
//
//		if err != nil {
//			return nil, err
//		}
//	}
//
//	return attrs, nil
//}
