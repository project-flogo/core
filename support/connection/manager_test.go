package connection

import (
	"github.com/project-flogo/core/support"
	"reflect"
	"testing"
)

type TestManager struct {
	mType string
}

func (t *TestManager) Type() string {
	return t.mType
}

func (t *TestManager) GetConnection() interface{} {
	return nil
}

func (t *TestManager) ReleaseConnection(connection interface{}) {
}

type TestManagerFactory struct {
	mfType string
}

func (t *TestManagerFactory) Type() string {
	return t.mfType
}

func (t *TestManagerFactory) NewManager(settings map[string]interface{}) (Manager, error) {
	return &TestManager{mType: t.mfType}, nil
}

func init() {
	support.RegisterAlias("connection", "connRef", "long/connRef")
}

func TestIsShared(t *testing.T) {

	manager1 := &TestManager{"testType"}
	managers = make(map[string]Manager)
	managers["testId"] = manager1
	manager2 := &TestManager{"testType2"}

	type args struct {
		manager Manager
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"shared-true", args{manager1}, true},
		{"shared-false", args{manager2}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsShared(tt.args.manager); got != tt.want {
				t.Errorf("IsShared() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewManager(t *testing.T) {

	managerF1 := &TestManagerFactory{"testType"}
	managerFactories = make(map[string]ManagerFactory)
	managerFactories["factoryRef"] = managerF1

	type args struct {
		config *Config
	}
	tests := []struct {
		name    string
		args    args
		want    Manager
		wantErr bool
	}{
		{"registered", args{&Config{Ref: "factoryRef"}}, &TestManager{mType: "testType"}, false},
		{"not-registered", args{&Config{Ref: "badFactoryRef"}}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewManager(tt.args.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewManager() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewManager() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewSharedManager(t *testing.T) {

	managerF1 := &TestManagerFactory{"testType"}
	managerFactories = make(map[string]ManagerFactory)
	managerFactories["factoryRef"] = managerF1

	type args struct {
		id     string
		config *Config
	}
	tests := []struct {
		name    string
		args    args
		want    Manager
		wantErr bool
	}{
		{"registered", args{"id1", &Config{Ref: "factoryRef"}}, &TestManager{mType: "testType"}, false},
		{"not-registered", args{"id2", &Config{Ref: "badFactoryRef"}}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewSharedManager(tt.args.id, tt.args.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewSharedManager() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSharedManager() got = %v, want %v", got, tt.want)
			}
			if tt.want != nil && managers[tt.args.id] == nil {
				t.Errorf("NewSharedManager() id  '%v' not shared", tt.args.id)
			}

		})
	}
}
