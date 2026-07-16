package app

import (
	"encoding/json"
	"testing"

	_ "github.com/project-flogo/core/examples/action"
	_ "github.com/project-flogo/core/examples/trigger"
	"github.com/project-flogo/core/support"
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

func TestResolveAliasConflict_GrandparentUnique(t *testing.T) {
	ref := "github.com/org/contrib/activity/query"
	result := resolveAliasConflict(ref, "query", "type1")
	assert.Equal(t, "contrib_query", result)
}

func TestResolveAliasConflict_GrandparentAlsoTaken(t *testing.T) {
	_ = support.RegisterAlias("type2", "contrib_query", "some/other/ref")

	ref := "github.com/org/contrib/activity/query"
	result := resolveAliasConflict(ref, "query", "type2")
	assert.Equal(t, "query_2", result)
}

func TestResolveAliasConflict_NumericFallbackSkipsTaken(t *testing.T) {
	_ = support.RegisterAlias("type3", "contrib_query", "some/other/ref")
	_ = support.RegisterAlias("type3", "query_2", "another/ref")
	_ = support.RegisterAlias("type3", "query_3", "yet/another/ref")

	ref := "github.com/org/contrib/activity/query"
	result := resolveAliasConflict(ref, "query", "type3")
	assert.Equal(t, "query_4", result)
}

func TestResolveAliasConflict_ShallowRefNoGrandparent(t *testing.T) {
	ref := "query"
	result := resolveAliasConflict(ref, "query", "type4")
	assert.Equal(t, "query_2", result)
}

func TestResolveAliasConflict_SingleSegmentParent(t *testing.T) {
	ref := "activity/query"
	result := resolveAliasConflict(ref, "query", "type5")
	assert.Equal(t, "query_2", result)
}

func TestResolveAliasConflict_DifferentContribTypes(t *testing.T) {
	_ = support.RegisterAlias("type6_other", "activity_log", "some/trigger/ref")

	ref := "github.com/org/contrib/activity/log"
	result := resolveAliasConflict(ref, "log", "type6")
	assert.Equal(t, "contrib_log", result)
}

func TestResolveAliasConflict_NumericStartsAt2(t *testing.T) {
	_ = support.RegisterAlias("type7", "contrib_query", "some/ref")

	ref := "github.com/org/contrib/activity/query"
	result := resolveAliasConflict(ref, "query", "type7")
	assert.Equal(t, "query_2", result)
}

func TestResolveAliasConflict_EmptyRef(t *testing.T) {
	result := resolveAliasConflict("", "query", "type8")
	assert.Equal(t, "query_2", result)
}
