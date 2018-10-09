package data

var propertyProvider PropertyProvider

func init() {
	propertyProvider = &DefaultPropertyProvider{}
}

type PropertyProvider interface {
	GetProperty(property string) (value interface{}, exists bool)
}

func SetPropertyProvider(provider PropertyProvider) {
	propertyProvider = provider
}

func GetPropertyProvider() PropertyProvider {
	return propertyProvider
}

// DefaultPropertyProvider empty property provider
type DefaultPropertyProvider struct {
}

func (pp *DefaultPropertyProvider) GetProperty(property string) (value interface{}, exists bool) {
	return nil, false
}
