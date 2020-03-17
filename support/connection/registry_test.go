package connection

import (
	"reflect"
	"testing"
)

func TestManagerFactories(t *testing.T) {

	managerF1 := &TestManagerFactory{"testType"}
	managerFactories = make(map[string]ManagerFactory)
	managerFactories["testId1"] = managerF1
	managerF2 := &TestManagerFactory{"testType2"}
	managerFactories["testId2"] = managerF2

	tests := []struct {
		name string
		want map[string]ManagerFactory
	}{
		{"simple", map[string]ManagerFactory{"testId1": managerF1, "testId2": managerF2}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ManagerFactories(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ManagerFactories() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestManagers(t *testing.T) {

	manager1 := &TestManager{"testType"}
	managers = make(map[string]Manager)
	managers["testId1"] = manager1
	manager2 := &TestManager{"testType2"}
	managers["testId2"] = manager2

	tests := []struct {
		name string
		want map[string]Manager
	}{
		{"simple", map[string]Manager{"testId1": manager1, "testId2": manager2}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Managers(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Managers() = %v, want %v", got, tt.want)
			}
		})
	}
}


func TestRegisterManager(t *testing.T) {

	managers = make(map[string]Manager)
	manager1 := &TestManager{"testType"}

	type args struct {
		connectionId string
		manager      Manager
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"simple", args{"id1", manager1}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := RegisterManager(tt.args.connectionId, tt.args.manager); (err != nil) != tt.wantErr {
				t.Errorf("RegisterManager() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}


func TestGetManager(t *testing.T) {

	manager1 := &TestManager{"testType"}
	managers = make(map[string]Manager)
	managers["testId"] = manager1

	type args struct {
		id string
	}
	tests := []struct {
		name string
		args args
		want Manager
	}{
		{"simple", args{"testId"}, manager1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetManager(tt.args.id); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetManager() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetManagerFactory(t *testing.T) {

	managerF1 := &TestManagerFactory{"testType"}
	managerFactories = make(map[string]ManagerFactory)
	managerFactories["testId"] = managerF1

	type args struct {
		ref string
	}
	tests := []struct {
		name string
		args args
		want ManagerFactory
	}{
		{"simple", args{"testId"}, managerF1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetManagerFactory(tt.args.ref); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetManagerFactory() = %v, want %v", got, tt.want)
			}
		})
	}
}


func TestRegisterManagerFactory(t *testing.T) {

	managerF1 := &TestManagerFactory{"testType"}
	managerFactories = make(map[string]ManagerFactory)

	type args struct {
		factory ManagerFactory
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"simple", args{managerF1}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := RegisterManagerFactory(tt.args.factory); (err != nil) != tt.wantErr {
				t.Errorf("RegisterManagerFactory() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestReplaceManagerFactory(t *testing.T) {
	managerF1 := &TestManagerFactory{"testType"}
	managerFactories = make(map[string]ManagerFactory)
	managerFactories["factoryRef"] = managerF1
	managerF2 := &TestManagerFactory{"testType2"}

	type args struct {
		ref     string
		factory ManagerFactory
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"simple", args{"factoryRef", managerF2}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ReplaceManagerFactory(tt.args.ref, tt.args.factory); (err != nil) != tt.wantErr {
				t.Errorf("ReplaceManagerFactory() error = %v, wantErr %v", err, tt.wantErr)
			}
			if managerFactories[tt.args.ref] != tt.args.factory {
				t.Errorf("ReplaceManagerFactory() ref '%v' not replaced", tt.args.ref)
			}
		})
	}
}

