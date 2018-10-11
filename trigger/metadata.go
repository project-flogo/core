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
