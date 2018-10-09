package activity

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockActivity struct {
	metadata *Metadata
}

func (MockActivity) Eval(ctx Context) (done bool, err error) {
	return true, nil
}

func (m *MockActivity) Metadata() *Metadata {
	return m.metadata
}

func NewMockActivity(id string) *MockActivity {
	return &MockActivity{metadata: &Metadata{ID: id}}
}

//TestRegisterNilActivity
func TestRegisterNilActivity(t *testing.T) {

	orig := activities
	activities = make(map[string]Activity)
	defer func() { activities = orig }()

	err := registerNil()
	assert.NotNil(t, err)
	assert.Equal(t, "cannot register 'nil' activity", err.Error())
}

func registerNil() (err error) {
	defer func() {
		if r := recover(); r != nil {
			var ok bool
			err, ok = r.(error)
			if !ok {
				err = fmt.Errorf("%v", r)
			}
		}
	}()

	Register(nil)
	return nil
}

//TestRegisterDupActivity
func TestRegisterDupActivity(t *testing.T) {

	orig := activities
	activities = make(map[string]Activity)
	defer func() { activities = orig }()

	err := registerDup()
	assert.NotNil(t, err)
	assert.Equal(t, "activity already registered: github.com/mock", err.Error())
}

func registerDup() (err error) {
	defer func() {
		if r := recover(); r != nil {
			var ok bool
			err, ok = r.(error)
			if !ok {
				err = fmt.Errorf("%v", r)
			}
		}
	}()

	act := NewMockActivity("github.com/mock")
	Register(act)
	Register(act)
	return nil
}

//TestRegisterActivityOk
func TestRegisterActivityOk(t *testing.T) {

	orig := activities
	activities = make(map[string]Activity)
	defer func() {
		activities = orig
		r := recover()
		assert.Nil(t, r)
	}()

	act := NewMockActivity("github.com/mock")
	Register(act)
	assert.Equal(t, 1, len(activities))
}

//TestGetFactoriesOk
func TestGetFactoriesOk(t *testing.T) {

	orig := activities
	activities = make(map[string]Activity)
	defer func() {
		activities = orig
		r := recover()
		assert.Nil(t, r)
	}()

	act := NewMockActivity("github.com/mock")
	Register(act)

	// Get factory
	as := Activities()
	assert.Equal(t, 1, len(as))
}
