package api

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/project-flogo/core/action"
	"github.com/project-flogo/core/app/resource"
	"github.com/project-flogo/core/data"
	"github.com/project-flogo/core/data/coerce"
	"github.com/project-flogo/core/data/metadata"
	"github.com/project-flogo/core/engine"
	"github.com/project-flogo/core/support"
	"github.com/project-flogo/core/trigger"
)

// App is the structure that defines an application
type App struct {
	properties map[string]data.TypedValue
	triggers   []*Trigger
	actions    map[string]*Action
	resources  []*resource.Config
}

// Trigger is the structure that defines a Trigger for the application
type Trigger struct {
	app      *App
	ref      string
	settings map[string]interface{}
	handlers []*Handler
}

// Handler is the structure that defines the handler for a Trigger
type Handler struct {
	app      *App
	settings map[string]interface{}
	actions  []*Action
	name     string
}

// HandlerFunc is the signature for a function to use as a handler for a Trigger
type HandlerFunc func(ctx context.Context, inputs map[string]interface{}) (map[string]interface{}, error)

// Action is the structure that defines the Action for a Handler
type Action struct {
	id             string
	ref            string
	act            action.Action
	settings       map[string]interface{}
	condition      string
	inputMappings  []string
	outputMappings []string
}

// NewApp creates a new Flogo application
func NewApp() *App {
	return &App{
		properties: make(map[string]data.TypedValue),
		actions:    make(map[string]*Action),
	}
}

// NewTrigger adds a new trigger to the application
func (a *App) NewTrigger(trg trigger.Trigger, settings interface{}) *Trigger {

	var settingsMap map[string]interface{}

	if settings != nil {
		if s, ok := settings.(map[string]interface{}); ok {
			settingsMap = s
		} else {
			settingsMap = metadata.StructToMap(settings)
		}
	}

	var ref string

	if hr, ok := trg.(support.HasRef); ok {
		ref = hr.Ref()
	} else {
		value := reflect.ValueOf(trg)
		value = value.Elem()
		ref = value.Type().PkgPath()
	}

	newTrg := &Trigger{app: a, ref: ref, settings: settingsMap}
	a.triggers = append(a.triggers, newTrg)

	return newTrg
}

// AddAction adds an action to the application
func (a *App) AddAction(id string, act action.Action, settings interface{}) error {

	newAct, err := newAction(act, settings)
	if err != nil {
		return err
	}
	a.actions[id] = newAct

	return nil
}

// AddProperty adds a shared property to the application
func (a *App) AddProperty(name string, dataType data.Type, value interface{}) error {
	tv, err := coerce.NewTypedValue(dataType, value)
	if err != nil {
		return err
	}
	a.properties[name] = tv
	return nil
}

// AddResource adds a Flogo resource to the application
func (a *App) AddResource(id string, data json.RawMessage) {

	res := &resource.Config{ID: id, Data: data}
	a.resources = append(a.resources, res)
}

// Properties gets the shared properties of the application
func (a *App) Properties() map[string]data.TypedValue {
	return a.properties
}

// Triggers gets the Triggers of the application
func (a *App) Triggers() []*Trigger {
	return a.triggers
}

// Triggers gets the Triggers of the application
func (a *App) Actions() map[string]*Action {
	return a.actions
}

// Settings gets the Trigger's settings
func (t *Trigger) Settings() map[string]interface{} {
	return t.settings
}

// NewHandler adds a new Handler to the Trigger
func (t *Trigger) NewHandler(settings interface{}) (*Handler, error) {

	var settingsMap map[string]interface{}

	if s, ok := settings.(map[string]interface{}); ok {
		settingsMap = s
	} else {
		settingsMap = metadata.StructToMap(settings)
	}

	newHandler := &Handler{
		app:      t.app,
		settings: settingsMap,
	}
	t.handlers = append(t.handlers, newHandler)

	return newHandler, nil
}

// Handlers gets the Trigger's Handlers
func (t *Trigger) Handlers() []*Handler {
	return t.handlers
}

// NewAction adds a new Action to the Handler
func (h *Handler) NewAction(handlerAction interface{}, settings ...interface{}) (act *Action, err error) {
	switch v := handlerAction.(type) {
	case HandlerFunc:
		act = &Action{act: NewProxyAction(v)}
	case func(ctx context.Context, inputs map[string]interface{}) (map[string]interface{}, error):
		act = &Action{act: NewProxyAction(v)}
	case action.Action:
		if len(settings) > 0 {
			act, err = newAction(v, settings[0])
		} else {
			act, err = newAction(v, nil)
		}
	case string:
		_, exists := h.app.actions[v]
		if !exists {
			return nil, fmt.Errorf("unknown action with id: %s", v)
		}
		act = &Action{id: v}
	default:
		return nil, fmt.Errorf("cannot create action from specified object: %+v", handlerAction)
	}
	h.actions = append(h.actions, act)
	return act, nil
}

// Settings gets the Handler's settings
func (h *Handler) Settings() map[string]interface{} {
	return h.settings
}

// Actions gets the Actions of the Handler
func (h *Handler) Actions() []*Action {
	return h.actions
}

// Settings gets the settings of the Action
func (a *Action) Settings() map[string]interface{} {
	return a.settings
}

// SetCondition sets the conditional expression which determines
// if the action is executed
func (a *Action) SetCondition(condition string) {
	a.condition = condition
}

// Condition returns the condition
func (a *Action) Condition() string {
	return a.condition
}

// SetInputMappings sets the input mappings for the Action, which maps
// the outputs of the Trigger to the inputs of the Action
func (a *Action) SetInputMappings(mappings ...string) {
	a.inputMappings = mappings
}

// SetOutputMappings sets the output mappings for the Action, which maps
// the outputs of the Action to the return of the Trigger
func (a *Action) SetOutputMappings(mappings ...string) {
	a.outputMappings = mappings
}

// InputMappings gets the Action's input mappings
func (a *Action) InputMappings() []string {
	return a.inputMappings
}

// OutputMappings gets the Action's output mappings
func (a *Action) OutputMappings() []string {
	return a.outputMappings
}

// NewEngine creates a new flogo Engine from the specified App
func NewEngine(a *App) (engine.Engine, error) {
	appConfig := toAppConfig(a)
	return engine.New(appConfig)
}

func newAction(act action.Action, settings interface{}) (*Action, error) {

	var settingsMap map[string]interface{}

	if settings, ok := settings.(map[string]interface{}); ok {
		settingsMap = settings
	} else {
		settingsMap = metadata.StructToMap(settings)
	}

	var ref string

	if hr, ok := act.(support.HasRef); ok {
		ref = hr.Ref()
	} else {
		value := reflect.ValueOf(act)
		value = value.Elem()
		ref = value.Type().PkgPath()
	}

	newAct := &Action{ref: ref, settings: settingsMap}

	return newAct, nil
}
