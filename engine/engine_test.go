package engine

import (
	"github.com/project-flogo/core/app"
	"github.com/stretchr/testify/assert"
	"testing"
)

//TestNewEngineErrorNoApp
func TestNewEngineErrorNoApp(t *testing.T) {
	_, err := New(nil)

	assert.NotNil(t, err)
	assert.Equal(t, "no App configuration provided", err.Error())
}

//TestNewEngineErrorNoAppName
func TestNewEngineErrorNoAppName(t *testing.T) {
	app := &app.Config{}

	_, err := New(app)

	assert.NotNil(t, err)
	assert.Equal(t, "no App name provided", err.Error())
}

//TestNewEngineErrorNoAppVersion
func TestNewEngineErrorNoAppVersion(t *testing.T) {
	app := &app.Config{Name: "MyApp"}

	_, err := New(app)

	assert.NotNil(t, err)
	assert.Equal(t, "no App version provided", err.Error())
}
