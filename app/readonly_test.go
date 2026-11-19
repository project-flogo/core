package app

import (
	"encoding/json"
	"testing"

	"github.com/project-flogo/core/data/property"
	"github.com/stretchr/testify/assert"
)

type overrideResolver struct{}

func (overrideResolver) Name() string { return "override-all" }

func (overrideResolver) LookupValue(string) (interface{}, bool) { return "override", true }

func TestNew_ReadOnlyPropertiesIgnoreOverrides(t *testing.T) {
	previousResolvers := property.EnabledResolvers
	previousManager := property.DefaultManager()
	property.EnabledResolvers = []property.ExternalResolver{overrideResolver{}}
	t.Cleanup(func() {
		property.EnabledResolvers = previousResolvers
		property.SetDefaultManager(previousManager)
	})

	var cfg Config
	assert.NoError(t, json.Unmarshal([]byte(`{"name":"ro","version":"1.0.0","properties":[
		{"name":"locked","type":"string","value":"app","readOnly":true},
		{"name":"legacyLocked","type":"string","value":"app","overridable":false},
		{"name":"noFlag","type":"string","value":"app"},
		{"name":"unlocked","type":"string","value":"app","readOnly":false,"overridable":false},
		{"name":"legacyUnlocked","type":"string","value":"app","overridable":true}]}`), &cfg))

	_, err := New(&cfg, nil)
	assert.NoError(t, err)
	assert.NoError(t, property.DefaultManager().Finalize(property.ExternalResolverProcessor))

	assert.Equal(t, map[string]interface{}{"locked": "app", "legacyLocked": "app", "noFlag": "override",
		"unlocked": "override", "legacyUnlocked": "override"}, property.DefaultManager().GetProperties())
}
