package app

import (
	"encoding/json"
	"testing"

	_ "github.com/project-flogo/core/examples/action"
	_ "github.com/project-flogo/core/examples/trigger"
	"github.com/stretchr/testify/assert"
)

var app = `{
	"name": "_APP_NAME_",
	"type": "flogo:app",
	"version": "0.0.1",
	"description": "My flogo application description",
	"appModel": "1.1.0",
	"imports": [
	  "github.com/project-flogo/core/examples/trigger",
	  "github.com/project-flogo/core/examples/action"
	],
	"triggers": [
	  {
		"id": "my_trigger",
		"ref": "github.com/project-flogo/core/examples/trigger",
		"settings": {
		  "aSetting": 2
		},
		"handlers": [
		  {
			"settings": {
			  "aSetting": 2
			},
			"actions": [
			  {
				"ref": "github.com/project-flogo/core/examples/action",
				"settings": {
				  "aSetting": "a"
				},
				"input": {
				  "in": "=$.anOutput"
				}
		  }
		]
		  }
		]
	  }
	]
  }
  `

func TestApp(t *testing.T) {
	var cfg *Config
	err := json.Unmarshal([]byte(app), &cfg)
	assert.Nil(t, err)

	app, err := New(cfg, nil, ContinueOnError)
	assert.Nil(t, err)
	assert.NotNil(t, app)

	assert.Equal(t, "_APP_NAME_", app.Name())
	assert.Equal(t, "0.0.1", app.Version())

	err = app.Start()
	assert.Nil(t, err)

	err = app.Stop()
	assert.Nil(t, err)
}
