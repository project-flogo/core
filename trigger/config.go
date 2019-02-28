package trigger

import (
	"encoding/json"
	"github.com/project-flogo/core/action"
	"github.com/project-flogo/core/data/expression"
	"github.com/project-flogo/core/data/mapper"
	"github.com/project-flogo/core/data/metadata"
	"github.com/project-flogo/core/data/resolve"
)

// Config is the configuration for a Trigger
type Config struct {
	Id       string                 `json:"id"`
	Type     string                 `json:"type"` //an alias to the ref, can be used if imported
	Ref      string                 `json:"ref"`
	Settings map[string]interface{} `json:"settings"`
	Handlers []*HandlerConfig       `json:"handlers"`
}

func (c *Config) FixUp(md *Metadata) error {

	ef := expression.NewFactory(resolve.GetBasicResolver())

	//fix up settings
	if len(c.Settings) > 0 {
		var err error
		mdSettings := md.Settings
		for name, value := range c.Settings {
			c.Settings[name], err = metadata.ResolveSettingValue(name, value, mdSettings, ef)
			if err != nil {
				return err
			}
		}
	}

	// fix up handler settings
	for _, hc := range c.Handlers {
		hc.parent = c

		if len(hc.Settings) > 0 {
			var err error
			mdSettings := md.Settings
			for name, value := range hc.Settings {
				hc.Settings[name], err = metadata.ResolveSettingValue(name, value, mdSettings, ef)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

type HandlerConfig struct {
	parent   *Config
	Name     string                 `json:"name,omitempty"`
	Settings map[string]interface{} `json:"settings"`
	Actions  []*ActionConfig        `json:"actions"`

	//handle complex object
}

// UnmarshalJSON overrides the default UnmarshalJSON for TaskInst
func (hc *HandlerConfig) UnmarshalJSON(d []byte) error {
	ser := &struct {
		Name     string                 `json:"name,omitempty"`
		Settings map[string]interface{} `json:"settings"`
		Actions  []*ActionConfig        `json:"actions"`
		Action   *ActionConfig          `json:"action"`
	}{}

	if err := json.Unmarshal(d, ser); err != nil {
		return err
	}

	hc.Name = ser.Name
	hc.Settings = ser.Settings

	if ser.Action != nil {
		hc.Actions = []*ActionConfig{ser.Action}
	} else {
		hc.Actions = ser.Actions
	}

	return nil
}

// ActionConfig is the configuration for the Action
type ActionConfig struct {
	*action.Config
	If     string                 `json:"if,omitempty"`
	Input  map[string]interface{} `json:"input,omitempty"`
	Output map[string]interface{} `json:"output,omitempty"`

	Act action.Action
}

// UnmarshalJSON overrides the default UnmarshalJSON for TaskInst
func (ac *ActionConfig) UnmarshalJSON(d []byte) error {
	ser := &struct {
		If     string                 `json:"if,omitempty"`
		Input  map[string]interface{} `json:"input,omitempty"`
		Output map[string]interface{} `json:"output,omitempty"`

		Ref      string                 `json:"ref"`
		Type     string                 `json:"type"`
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

	ac.Config = &action.Config{}

	ac.Ref = ser.Ref
	ac.Type = ser.Type
	ac.Id = ser.Id
	ac.If = ser.If
	ac.Input = ser.Input
	ac.Output = ser.Output
	ac.Settings = ser.Settings

	if ac.Settings == nil {
		ac.Settings = make(map[string]interface{})
	}

	if ser.Data != nil {
		for key, value := range ser.Data {
			ac.Settings[key] = value
		}
	}

	input, output, err := mapper.ConvertLegacyMappings(ser.Mappings, resolve.GetBasicResolver())
	if err != nil {
		return err
	}

	if ac.Input == nil {
		ac.Input = input
	} else {
		for key, value := range input {
			ac.Input[key] = value
		}
	}

	if ac.Output == nil {
		ac.Output = output
	} else {
		for key, value := range output {
			ac.Output[key] = value
		}
	}

	return nil
}
