package trigger

import (
	"testing"

	"github.com/project-flogo/core/action"
	"github.com/stretchr/testify/assert"
)

type MockFactory struct {
}

func (f *MockFactory) Metadata() *Metadata {
	return nil
}

func (f *MockFactory) New(config *Config) (Trigger, error) {
	return &MockTrigger{}, nil
}

type MockTrigger struct {
}

func (t *MockTrigger) Initialize(ctx InitContext) error {
	//ignore
	return nil
}

func (t *MockTrigger) Init(actionRunner action.Runner) {
	//Noop
}

func (t *MockTrigger) Start() error        { return nil }
func (t *MockTrigger) Stop() error         { return nil }
func (t *MockTrigger) Metadata() *Metadata { return nil }

//TestRegisterFactoryEmptyRef
func TestRegisterNilTrigger(t *testing.T) {

	orig := triggerFactories
	triggerFactories = make(map[string]Factory)
	defer func() { triggerFactories = orig }()

	// Register factory
	err := Register(nil, nil)

	assert.NotNil(t, err)
	assert.Equal(t, "'trigger' must be specified when registering", err.Error())
}

//TestRegisterFactoryNilFactory
func TestRegisterFactoryNilFactory(t *testing.T) {

	orig := triggerFactories
	triggerFactories = make(map[string]Factory)
	defer func() { triggerFactories = orig }()

	// Register factory
	err := Register(&MockTrigger{}, nil)

	assert.NotNil(t, err)
	assert.Equal(t, "cannot register trigger with 'nil' trigger factory", err.Error())
}

//TestAddFactoryDuplicated
func TestAddFactoryDuplicated(t *testing.T) {

	orig := triggerFactories
	triggerFactories = make(map[string]Factory)
	defer func() { triggerFactories = orig }()

	f := &MockFactory{}

	// Register factory: this time should pass
	err := Register(&MockTrigger{}, f)
	assert.Nil(t, err)

	// Register factory: this time should fail, duplicated
	err = Register(&MockTrigger{}, f)
	assert.NotNil(t, err)
	assert.Equal(t, "trigger already registered for ref github.com/project-flogo/core/trigger", err.Error())
}

//TestAddFactoryOk
func TestAddFactoryOk(t *testing.T) {

	orig := triggerFactories
	triggerFactories = make(map[string]Factory)
	defer func() { triggerFactories = orig }()

	f := &MockFactory{}

	// Add factory
	err := LegacyRegister("github.com/mock", f)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(triggerFactories))
}

//TestGetFactoriesOk
func TestGetFactoriesOk(t *testing.T) {

	orig := triggerFactories
	triggerFactories = make(map[string]Factory)
	defer func() { triggerFactories = orig }()

	f := &MockFactory{}

	// Add factory
	err := LegacyRegister("github.com/mock", f)
	assert.Nil(t, err)

	// Get factory
	fs := Factories()
	assert.Equal(t, 1, len(fs))
}
