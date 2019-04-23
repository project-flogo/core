package sample

import (
	"context"

	"github.com/project-flogo/core/data/coerce"
	"github.com/project-flogo/core/data/metadata"
	"github.com/project-flogo/core/trigger"
)

var triggerMd = trigger.NewMetadata(&Settings{}, &HandlerSettings{}, &Output{}, &Reply{})

func init() {
	_ = trigger.Register(&Trigger{}, &Factory{})
}

type Trigger struct {
	settings *Settings
	id       string
}

type Factory struct {
}

func (*Factory) New(config *trigger.Config) (trigger.Trigger, error) {
	s := &Settings{}
	err := metadata.MapToStruct(config.Settings, s, true)
	if err != nil {
		return nil, err
	}

	return &Trigger{id: config.Id, settings: s}, nil
}

func (f *Factory) Metadata() *trigger.Metadata {
	return triggerMd
}

// Metadata implements trigger.Trigger.Metadata
func (t *Trigger) Metadata() *trigger.Metadata {
	return triggerMd
}

func (t *Trigger) Initialize(ctx trigger.InitContext) error {

	logger := ctx.Logger()

	aSetting := t.settings.ASetting
	logger.Debug("Setting 'aSetting' = %s", aSetting)

	// Init handlers
	for _, handler := range ctx.GetHandlers() {

		s := &HandlerSettings{}
		err := metadata.MapToStruct(handler.Settings(), s, true)
		if err != nil {
			return err
		}

		//init handler using setting
		aSetting := s.ASetting
		logger.Debug("Handler setting 'aSetting' = %s", aSetting)

		registerDummyEventHandler(aSetting, newActionHandler(handler))
	}

	return nil
}

// Start implements util.Managed.Start
func (t *Trigger) Start() error {
	//start servers/services if necessary
	return nil
}

// Stop implements util.Managed.Stop
func (t *Trigger) Stop() error {
	//stop servers/services if necessary
	return nil
}

// registerDummyEventHandler is used for dummy event handler registration, this should be replaced
// with the appropriate event handling mechanism for the trigger.  Some form of a discriminator
// should be used for dispatching to different handlers.  For example a REST based trigger might
// dispatch based on the method and path.
func registerDummyEventHandler(discriminator string, onEvent dummyOnEvent) {
	//ignore
}

// dummyOnEvent is a dummy event handler for our dummy event source
type dummyOnEvent func(interface{})

func newActionHandler(handler trigger.Handler) dummyOnEvent {

	return func(data interface{}) {

		strData, _ := coerce.ToString(data)
		output := &Output{AnOutput: strData}

		results, err := handler.Handle(context.Background(), output.ToMap())
		if err != nil {
			//handle error
		}
		reply := &Reply{}
		_ = reply.FromMap(results)

		//do something with the reply
	}
}
