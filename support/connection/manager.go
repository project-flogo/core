package connection

type Manager interface {
	Type()

	GetConnection() interface{}

	//ReleaseConnection(connection interface{})
}

type ManagerFactory interface {
	Type()

	NewManager(settings map[string]interface{}) (Manager, error)
}

