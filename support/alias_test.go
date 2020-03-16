package support

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetAliasRef(t *testing.T) {

	aliases = make(map[string]map[string]string)
	aliases["activity"] = make(map[string]string)
	aliases["activity"]["alias"] = "fullRef"

	type args struct {
		contribType string
		alias       string
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 bool
	}{
		{"simple", args{"activity", "alias"}, "fullRef", true},
		{"dne", args{"activity", "alias2"}, "", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := GetAliasRef(tt.args.contribType, tt.args.alias)
			if got != tt.want {
				t.Errorf("GetAliasRef() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("GetAliasRef() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestRegisterAlias(t *testing.T) {
	aliases = make(map[string]map[string]string)

	type args struct {
		contribType string
		alias       string
		ref         string
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"simple", args{"activity", "alias", "fullRef"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := RegisterAlias(tt.args.contribType, tt.args.alias, tt.args.ref); (err != nil) != tt.wantErr {
				t.Errorf("RegisterAlias() error = %v, wantErr %v", err, tt.wantErr)
			}

			assert.NotNil(t, aliases[tt.args.contribType])
			assert.Equal(t, aliases[tt.args.contribType][tt.args.alias], tt.args.ref)
		})
	}
}
