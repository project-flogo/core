package property

import (
	"encoding/json"
	"github.com/project-flogo/core/support/logger"
	"io/ioutil"
	"strings"
)

func init() {
	SetDefaultManager(NewManager(make(map[string]interface{})))
}

var defaultManager *Manager

func SetDefaultManager(manager *Manager) {
	defaultManager = manager
}

func DefaultManager() *Manager {
	return defaultManager
}

func NewManager(properties map[string]interface{}) *Manager {

	manager := &Manager{properties: properties}
	return manager
}

type Manager struct {
	properties map[string]interface{}
}

func (m *Manager) GetProperty(name string) (interface{}, bool) {
	val, exists := m.properties[name]
	return val, exists
}

func (m *Manager) AddExternalProperties(providerId string, overrides string, processors ...PostProcessor) error {

	provider := GetProvider(providerId)
	newProps, err := loadExternalProperties(provider, overrides)
	if err != nil {
		return err
	}

	for _, processor := range processors {
		processor(newProps)
	}

	for key, value := range newProps {
		m.properties[key] = value
	}

	return nil
}

type PostProcessor func(properties map[string]interface{}) error

func loadExternalProperties(provider Provider, overrides string) (map[string]interface{}, error) {

	props := make(map[string]interface{})

	if overrides != "" {
		if strings.HasSuffix(overrides, ".json") {
			// Override through file

			propFile := overrides
			file, e := ioutil.ReadFile(propFile)
			if e != nil {
				return nil, e
			}
			e = json.Unmarshal(file, &props)

			if e != nil {
				return nil, e
			}
		} else if strings.ContainsRune(overrides, '=') {
			// Override through P1=V1,P2=V2
			for _, pair := range strings.Split(overrides, ",") {
				kv := strings.Split(pair, "=")
				if len(kv) == 2 && kv[0] != "" {
					key := strings.TrimSpace(kv[0])
					value := strings.TrimSpace(kv[1])
					props[key] = value
				} else {
					logger.Warnf("'%s' is not valid override value. It must be in PropName=PropValue format.", pair)
				}
			}
		}
	}

	//look for properties that need to go to the provider
	if provider != nil {
		for key, value := range props {
			if strVal, ok := value.(string); ok && strVal[0] == '$' {
				val := provider.GetProperty(strVal[1:])
				props[key] = val
			}
		}
	}

	return props, nil
}
