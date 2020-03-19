package resource

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetTypeFromID(t *testing.T) {

	resType, err := GetTypeFromID("flow:myflow")
	assert.Nil(t, err)
	assert.Equal(t, "flow", resType)
}

type SampleResourceLoader struct {
}
type Data struct {
	Name string `json:"name"`
}

var config = `{"name": "a" } `

func (rl SampleResourceLoader) LoadResource(config *Config) (*Resource, error) {
	var data Data
	//fmt.Println("Config...", string(config.Data))
	err := json.Unmarshal(config.Data, &data)
	if err != nil {
		return nil, err
	}
	return New("sample", data), nil
}

func TestResourceLoader(t *testing.T) {
	var sampleLoader SampleResourceLoader
	err := RegisterLoader("sample", sampleLoader)
	assert.Nil(t, err)

	err = RegisterLoader("sample", sampleLoader)
	assert.NotNil(t, err)

	assert.NotNil(t, GetLoader("sample"))

}

func TestResource(t *testing.T) {
	var sampleLoader SampleResourceLoader

	sampleResource, err := sampleLoader.LoadResource(&Config{ID: "sample", Data: []byte(config)})
	assert.Nil(t, err)
	assert.NotNil(t, sampleResource)

	assert.NotNil(t, sampleResource.Type())
	assert.NotNil(t, sampleResource.Object())

	resManager := NewManager(map[string]*Resource{"sample": sampleResource})
	assert.NotNil(t, resManager)

	assert.Equal(t, resManager.GetResource("sample"), resManager.GetResource("res://sample"))
}
