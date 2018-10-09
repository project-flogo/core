package trigger

import (
	"encoding/json"
	"github.com/project-flogo/core/action"
	"github.com/project-flogo/core/data/coerce"
	"github.com/project-flogo/core/data/mapper"
)

// Config is the configuration for a Trigger
type Config struct {
	Id       string                 `json:"id"`
	Ref      string                 `json:"ref"`
	Settings map[string]interface{} `json:"settings"`
	Handlers []*HandlerConfig       `json:"handlers"`
}

func (c *Config) FixUp(metadata *Metadata) {

	// fix up top-level outputs
	//for name, value := range c.Output {
	//
	//	attr, ok := metadata.Output[name]
	//
	//	if ok {
	//		newValue, err := coerce.ToType(value, attr.Type())
	//
	//		if err != nil {
	//			//todo handle error
	//		} else {
	//			c.Output[name] = newValue
	//		}
	//	}
	//}

	// fix up handler outputs
	for _, hc := range c.Handlers {

		hc.parent = c

		////for backwards compatibility
		//if hc.ActionId == "" {
		//	hc.ActionId = strconv.Itoa(time.Now().Nanosecond())
		//}

		//// fix up outputs
		//for name, value := range hc.Output {
		//
		//	attr, ok := metadata.Output[name]
		//
		//	if ok {
		//		newValue, err := coerce.ToType(value, attr.Type())
		//
		//		if err != nil {
		//			//todo handle error
		//		} else {
		//			hc.Output[name] = newValue
		//		}
		//	}
		//}
	}
}

func (c *Config) GetSetting(setting string) string {

	//val, exists := data.GetValueWithResolver(c.Settings, setting)

	val, exists := c.Settings[setting]
	if !exists {
		return ""
	}

	strVal, err := coerce.ToString(val)
	if err != nil {
		return ""
	}

	return strVal
}

type HandlerConfig struct {
	parent   *Config
	Name     string                 `json:"name,omitempty"`
	Settings map[string]interface{} `json:"settings"`
	Action   *ActionConfig

	//Output   map[string]interface{} `json:"output"`
	//handle complex object
}

// ActionConfig is the configuration for the Action
type ActionConfig struct {
	*action.Config

	Input  map[string]interface{} `json:"input,omitempty"`
	Output map[string]interface{} `json:"output,omitempty"`

	Act action.Action
}

func (hc *HandlerConfig) GetSetting(setting string) string {

	//val, exists := data.GetValueWithResolver(c.Settings, setting)

	val, exists := hc.Settings[setting]
	if !exists {
		return ""
	}

	strVal, err := coerce.ToString(val)
	if err != nil {
		return ""
	}

	return strVal
}

// UnmarshalJSON overrides the default UnmarshalJSON for TaskInst
func (ti *ActionConfig) UnmarshalJSON(d []byte) error {
	ser := &struct {
		Input  map[string]interface{} `json:"input,omitempty"`
		Output map[string]interface{} `json:"output,omitempty"`

		Ref      string                 `json:"ref"`
		Settings map[string]interface{} `json:"settings"`

		//referenced action
		Id string `json:"id"`

		//DEPRECATED
		Data map[string]interface{} `json:"data"`
		//DEPRECATED
		Mappings *mapper.LegacyMappings `json:"mappings"`
	}{}

	if err := json.Unmarshal(d, ser); err != nil {
		return err
	}

	ti.Config = &action.Config{}

	ti.Ref = ser.Ref
	ti.Id = ser.Id
	ti.Input = ser.Input
	ti.Output = ser.Output
	ti.Settings = ser.Settings

	if ti.Settings == nil {
		ti.Settings = make(map[string]interface{})
	}

	if ser.Data != nil {
		for key, value := range ser.Data {
			ti.Settings[key] = value
		}
	}

	input, output, err := mapper.ConvertLegacyMappings(ser.Mappings)
	if err != nil {
		return err
	}

	if ti.Input == nil {
		ti.Input = input
	} else {
		for key, value := range input {
			ti.Input[key] = value
		}
	}

	if ti.Output == nil {
		ti.Output = output
	} else {
		for key, value := range output {
			ti.Output[key] = value
		}
	}

	return nil
}
