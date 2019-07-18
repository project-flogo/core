package connection

type Manager interface {
	GetConnection() interface{}

	//ReleaseConnection(connection interface{})
}

type ManagerFactory interface {
	NewManager(settings map[string]string) (Manager, error)
}

