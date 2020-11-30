package app

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestController(t *testing.T) {
	testApp := &App{stopOnError: true, name: "test app", version: "v1.0.0"}
	testApp.initFlowController()
	c := GetFlowController()
	err := c.StartControl()
	assert.Nil(t, err)
	// Start another control
	err = c.StartControl()
	assert.Equal(t, AlreadyControlled, err.Error())
	err = c.ReleaseControl()
	assert.Nil(t, err)

	//Start again
	err = c.StartControl()
	assert.Nil(t, err)
	err = c.ReleaseControl()
	assert.Nil(t, err)
}
