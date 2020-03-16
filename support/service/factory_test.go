package service

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type TestServiceFactory struct {

}

func (t *TestServiceFactory) NewService(config *Config) (Service, error) {
	return nil, nil
}

func TestGetFactory(t *testing.T) {

	ts := &TestServiceFactory{}

	serviceFactories = make(map[string]Factory)
	serviceFactories["test"] = ts

	f := GetFactory("test")
	assert.Equal(t, ts, f)

	e := GetFactory("nothere")
	assert.Nil(t, e)
}


func TestRegisterFactory(t *testing.T) {

	ts := &TestServiceFactory{}
	err := RegisterFactory(ts)
	assert.Nil(t, err)

	f := serviceFactories["github.com/project-flogo/core/support/service"]
	assert.Equal(t, ts, f)

	err = RegisterFactory(ts)
	assert.NotNil(t, err)
}
