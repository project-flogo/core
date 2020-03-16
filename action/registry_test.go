package action

import (
	"context"
	"errors"
	"testing"

	"github.com/project-flogo/core/data"
	"github.com/project-flogo/core/data/metadata"
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

func TestLegacyRegister(t *testing.T) {
	inputs := []struct {
		ref    string
		f      Factory
		result error
	}{
		{"", nil, errors.New("'action ref' must be specified when registering")},
		{"sample", nil, errors.New("cannot register action with 'nil' action factory")},
		{"sample", &MockFactory{}, nil},
		{"sample", &MockFactory{}, errors.New("action already registered: sample")},
	}

	for _, in := range inputs {
		assert.Equal(t, in.result, LegacyRegister(in.ref, in.f))
	}
}

func TestGetFactory(t *testing.T) {
	assert.NotNil(t, GetFactory("sample"))
	assert.Nil(t, GetFactory("github.com/TIBCOSoftware/flogo-contrib/action/flow"))
}
