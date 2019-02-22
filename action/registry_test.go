package action

import (
	"context"
	"testing"

	"github.com/project-flogo/core/data/metadata"

	"github.com/project-flogo/core/data"
	"github.com/stretchr/testify/assert"
)

type MockFactory struct {
}

func (f *MockFactory) Initialize(ctx InitContext) error {
	return nil
}

func (f *MockFactory) New(config *Config) (Action, error) {
	return &MockAction{}, nil
}

type MockAction struct {
}

func (t *MockAction) Metadata() *Metadata {
	return nil
}

func (t *MockAction) IOMetadata() *metadata.IOMetadata {
	return nil
}

func (t *MockAction) Run(context context.Context, inputs map[string]*data.Attribute, handler ResultHandler) error {
	return nil
}

//TestRegisterFactoryEmptyRef
func TestRegisterFactoryEmptyRef(t *testing.T) {

	orig := actionFactories
	actionFactories = make(map[string]Factory)
	defer func() { actionFactories = orig }()

	// Register factory
	err := Register(nil, nil)

	assert.NotNil(t, err)
	assert.Equal(t, "'action' must be specified when registering", err.Error())
}

//TestRegisterFactoryNilFactory
func TestRegisterFactoryNilFactory(t *testing.T) {

	orig := actionFactories
	actionFactories = make(map[string]Factory)
	defer func() { actionFactories = orig }()

	// Register factory
	err := Register(&MockAction{}, nil)

	assert.NotNil(t, err)
	assert.Equal(t, "cannot register action with 'nil' action factory", err.Error())
}

//TestAddFactoryDuplicated
func TestAddFactoryDuplicated(t *testing.T) {

	orig := actionFactories
	actionFactories = make(map[string]Factory)
	defer func() { actionFactories = orig }()

	f := &MockFactory{}

	// Register factory: this time should pass
	err := Register(&MockAction{}, f)
	assert.Nil(t, err)

	// Register factory: this time should fail, duplicated
	err = Register(&MockAction{}, f)
	assert.NotNil(t, err)
	//assert.Equal(t, "action factory already registered for ref 'github.com/mock'", err.Error())
}

//TestAddFactoryOk
func TestAddFactoryOk(t *testing.T) {

	orig := actionFactories
	actionFactories = make(map[string]Factory)
	defer func() { actionFactories = orig }()

	f := &MockFactory{}

	// Add factory
	err := Register(&MockAction{}, f)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(actionFactories))
}

//TestGetFactoriesOk
func TestGetFactoriesOk(t *testing.T) {

	orig := actionFactories
	actionFactories = make(map[string]Factory)
	defer func() { actionFactories = orig }()

	f := &MockFactory{}

	// Add factory
	err := Register(&MockAction{}, f)
	assert.Nil(t, err)

	// Get factory
	fs := Factories()
	assert.Equal(t, 1, len(fs))
}
