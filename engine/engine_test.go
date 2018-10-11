package engine

import (
	"testing"

	"github.com/project-flogo/core/app"
	"github.com/stretchr/testify/assert"
)

//TestNewEngineErrorNoApp
func TestNewEngineErrorNoApp(t *testing.T) {
	_, err := New(nil)

	assert.NotNil(t, err)
	assert.Equal(t, "no App configuration provided", err.Error())
}

//TestNewEngineErrorNoAppName
func TestNewEngineErrorNoAppName(t *testing.T) {
	appConfig := &app.Config{}

	_, err := New(appConfig)

	assert.NotNil(t, err)
	assert.Equal(t, "no App name provided", err.Error())
}

//TestNewEngineErrorNoAppVersion
func TestNewEngineErrorNoAppVersion(t *testing.T) {
	appConfig := &app.Config{Name: "MyApp"}

	_, err := New(appConfig)

	assert.NotNil(t, err)
	assert.Equal(t, "no App version provided", err.Error())
}
