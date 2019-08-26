package tmp

//DEPRECATED
type PropertyManager interface {
	GetProperty(name string) (interface{}, bool)
}

var defaultManager PropertyManager

//DEPRECATED
func SetDefaultManager(manager PropertyManager) {
	defaultManager = manager
}

//DEPRECATED
func DefaultManager() PropertyManager {
	return defaultManager
}

