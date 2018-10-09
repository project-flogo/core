package app

var propertyProvider *PropertyProvider

func init() {
	propertyProvider = &PropertyProvider{properties: make(map[string]interface{})}
}

func GetPropertyProvider() *PropertyProvider {
	return propertyProvider
}

type PropertyProvider struct {
	properties map[string]interface{}
}

func (pp *PropertyProvider) GetProperty(property string) (interface{}, bool) {
	prop, exists := pp.properties[property]
	return prop, exists
}

func (pp *PropertyProvider) SetProperty(property string, value interface{}) {
	pp.properties[property] = value
}

func (pp *PropertyProvider) SetProperties(value map[string]interface{}) {
	pp.properties = value
}
