package connection

import (
	"os"
	"reflect"
	"testing"

	appresolve "github.com/project-flogo/core/app/resolve"
	"github.com/project-flogo/core/data/resolve"
)

func TestResolveConfig(t *testing.T) {

	os.Setenv("TCVAL", "foo")
	defer func() {
		os.Unsetenv("TCVAL")
	}()
	cfg1 := &Config{Ref: "testRef1", Settings: map[string]interface{}{"setting1": "=$env[TCVAL]"}}
	cfg2 := &Config{Ref: "#connRef", Settings: map[string]interface{}{"setting1": "foo"}}

	resolver := resolve.NewCompositeResolver(map[string]resolve.Resolver{
		"env": &resolve.EnvResolver{},
	})
	appresolve.SetAppResolver(resolver)

	type args struct {
		config *Config
	}
	tests := []struct {
		name     string
		args     args
		wantErr  bool
		wantRef  string
		wantSVal string
	}{
		{"resolve-setting", args{cfg1}, false, "", "foo"},
		{"resolve-ref", args{cfg2}, false, "long/connRef", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ResolveConfig(tt.args.config); (err != nil) != tt.wantErr {
				t.Errorf("ResolveConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
			if got := tt.args.config.Settings["setting1"]; tt.wantSVal != "" && tt.wantSVal != got {
				t.Errorf("ResolveConfig() got = %v, want %v", got, tt.wantSVal)
			}
			if got := tt.args.config.Ref; tt.wantRef != "" && tt.wantRef != got {
				t.Errorf("ResolveConfig() got = %v, want %v", got, tt.wantRef)
			}
		})
	}

	appresolve.SetAppResolver(nil)
}

func TestToConfig(t *testing.T) {

	mapCfg1 := map[string]interface{}{"ref": "testRef1", "settings": map[string]interface{}{"setting1": "foo"}}
	cfg1 := &Config{Ref: "testRef1", Settings: map[string]interface{}{"setting1": "foo"}}

	type args struct {
		config map[string]interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    *Config
		wantErr bool
	}{
		{"simple", args{mapCfg1}, cfg1, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ToConfig(tt.args.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToConfig() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_resolveRef(t *testing.T) {

	cfg1 := &Config{Ref: "#connRef", Settings: map[string]interface{}{"setting1": "foo"}}

	type args struct {
		config *Config
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		wantRef string
	}{
		{"simple", args{cfg1}, false, "long/connRef"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := resolveRef(tt.args.config); (err != nil) != tt.wantErr {
				t.Errorf("resolveRef() error = %v, wantErr %v", err, tt.wantErr)
			}
			if got := tt.args.config.Ref; tt.wantRef != "" && tt.wantRef != got {
				t.Errorf("resolveRef() got = %v, want %v", got, tt.wantRef)
			}
		})
	}
}
