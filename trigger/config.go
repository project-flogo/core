package trigger

import (
	"encoding/json"
	"strconv"

	"github.com/project-flogo/core/action"
	"github.com/project-flogo/core/data/expression"
	"github.com/project-flogo/core/data/metadata"
	"github.com/project-flogo/core/data/resolve"
)

// Config is the configuration for a Trigger
type Config struct {
	Id       string                 `json:"id"`
	Ref      string                 `json:"ref"`
	Settings map[string]interface{} `json:"settings"`
	Handlers []*HandlerConfig       `json:"handlers"`

	//DEPRECATED
	Type string `json:"type,omitempty"`
}

func (c *Config) FixUp(md *Metadata, resolver resolve.CompositeResolver) error {

	ef := expression.NewFactory(resolver)

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
	for i, hc := range c.Handlers {
		hc.parent = c

		if hc.Name == "" {
			// Set name: <triggername>_handler<index+1> e.g. mytrigger_handler1
			hc.Name = c.Id + "_handler" + strconv.Itoa(i+1)
		}

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
	Schemas  *SchemaConfig          `json:"schemas,omitempty"`
}

type SchemaConfig struct {
	Output map[string]interface{} `json:"output,omitempty"`
	Reply  map[string]interface{} `json:"reply,omitempty"`
}

// UnmarshalJSON overrides the default UnmarshalJSON for TaskInst
func (hc *HandlerConfig) UnmarshalJSON(d []byte) error {
	ser := &struct {
		Name     string                 `json:"name,omitempty"`
		Settings map[string]interface{} `json:"settings"`
		Actions  []*ActionConfig        `json:"actions"`
		Action   *ActionConfig          `json:"action"`
		Schemas  *SchemaConfig          `json:"schemas,omitempty"`
	}{}

	if err := json.Unmarshal(d, ser); err != nil {
		return err
	}

	hc.Name = ser.Name
	hc.Settings = ser.Settings
	hc.Schemas = ser.Schemas

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

	Act action.Action `json:"-,omitempty"`
}