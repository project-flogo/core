package sample

import (
	"context"
	"fmt"

	"github.com/project-flogo/core/action"
	"github.com/project-flogo/core/api"
	"github.com/project-flogo/core/app/resource"
	"github.com/project-flogo/core/data/metadata"
	"github.com/project-flogo/core/support/log"
)

func init() {
	_ = action.Register(&Action{}, &ActionFactory{})
	resource.RegisterLoader("action", &Manager{})
}

var actionMd = action.ToMetadata(&Settings{}, &Input{}, &Output{})

// Manager loads the action definition resource
type Manager struct {
}

// LoadResource loads the action definition
func (m *Manager) LoadResource(config *resource.Config) (*resource.Resource, error) {
	return resource.New("action", "test data"), nil
}

type ActionFactory struct {
}

func (f *ActionFactory) Initialize(ctx action.InitContext) error {
	return nil
}

func (f *ActionFactory) New(config *action.Config) (action.Action, error) {

	s := &Settings{}
	err := metadata.MapToStruct(config.Settings, s, true)
	if err != nil {
		return nil, err
	}

	log.RootLogger().Debugf("Setting: %s", s.ASetting)

	act := &Action{settings: s}

	return act, nil

}

// Generate generates Flogo API code for this action
func (f *ActionFactory) Generate(settingsName string, imports *api.Imports, config *action.Config) (code string, err error) {
	return fmt.Sprintf("var %s = %#v\n", settingsName, config.Settings), nil
}

type Action struct {
	settings *Settings
}

// Metadata implements action.Action.Metadata
func (a *Action) Metadata() *action.Metadata {
	return actionMd
}

// IOMetadata implements action.Action.IOMetadata
func (a *Action) IOMetadata() *metadata.IOMetadata {
	return nil
}

// Run implements action.SyncAction.Run
func (a *Action) Run(ctx context.Context, inputValues map[string]interface{}) (map[string]interface{}, error) {

	input := &Input{}
	err := input.FromMap(inputValues)
	if err != nil {
		return nil, err
	}

	log.RootLogger().Infof("Input: %s", input.AnInput)

	output := &Output{AnOutput: input.AnInput}

	return output.ToMap(), nil
}
