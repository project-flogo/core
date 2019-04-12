package resource

import (
	"errors"
	"strings"

	"github.com/project-flogo/core/support"
	"github.com/project-flogo/core/support/log"
)

const UriScheme = "res://"

type Loader interface {
	LoadResource(config *Config) (*Resource, error)
}

var loaders = make(map[string]Loader)

// RegisterManager registers a resource manager for the specified type
func RegisterLoader(resourceType string, loader Loader) error {

	_, exists := loaders[resourceType]

	if exists {
		return errors.New("Resource Loader already registered for type: " + resourceType)
	}

	loaders[resourceType] = loader
	return nil
}

func GetLoader(resourceType string) Loader {
	return loaders[resourceType]
}

type Manager struct {
	resources map[string]*Resource
}

func NewManager(resources map[string]*Resource) *Manager {
	return &Manager{resources: resources}
}

func (m *Manager) GetResource(id string) *Resource {

	resId := id
	if strings.HasPrefix(id, UriScheme) {
		//is uri
		resId = id[6:]
	}

	return m.resources[resId]
}

func (m *Manager) CleanupResources() {

	for id, res := range m.resources {
		if needsCleanup, ok := res.resObj.(support.NeedsCleanup); ok {
			err := needsCleanup.Cleanup()
			if err != nil {
				log.RootLogger().Errorf("Error cleaning up resource '%s' : ", id, err)
			}
		}
	}
}

func GetTypeFromID(id string) (string, error) {

	idx := strings.Index(id, ":")

	if idx < 0 {
		return "", errors.New("Invalid resource id: " + id)
	}

	return id[:idx], nil
}
