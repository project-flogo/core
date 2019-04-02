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

func NewMockActivity() *MockActivity {
	return &MockActivity{metadata: &Metadata{}}
}

//TestRegisterNilActivity
func TestRegisterNilActivity(t *testing.T) {

	orig := activities
	activities = make(map[string]Activity)
	defer func() { activities = orig }()

	err := registerNil()
	if assert.NotNil(t, err) {
		assert.Equal(t, "cannot register 'nil' activity", err.Error())
	}
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

	err = Register(nil)
	return err
}

//TestRegisterDupActivity
func TestRegisterDupActivity(t *testing.T) {

	orig := activities
	activities = make(map[string]Activity)
	defer func() { activities = orig }()

	err := registerDup()
	if assert.NotNil(t, err) {
		assert.Equal(t, "activity already registered: github.com/project-flogo/core/activity", err.Error())
	}
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

	act := NewMockActivity()
	_ = Register(act)
	err = Register(act)
	return err
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

	act := NewMockActivity()
	_ = Register(act)
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

	act := NewMockActivity()
	_ = Register(act, testFactory)

	// Get factory
	f := GetFactory("github.com/project-flogo/core/activity")
	assert.NotNil(t, f)
}

func testFactory(ctx InitContext) (Activity, error) {
	return nil, nil
}
