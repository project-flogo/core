package data_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/project-flogo/core/data"
	_ "github.com/project-flogo/core/data/coerce"
	"github.com/stretchr/testify/assert"
)

func TestAttribute_UnmarshalReadOnly(t *testing.T) {
	tests := []struct {
		name     string
		flags    string
		readOnly bool
	}{
		{"no flag", ``, false},
		{"readOnly true", `,"readOnly":true`, true},
		{"readOnly false", `,"readOnly":false`, false},
		{"readOnly null", `,"readOnly":null`, false},
		{"legacy overridable false", `,"overridable":false`, true},
		{"legacy overridable true", `,"overridable":true`, false},
		{"legacy overridable null", `,"overridable":null`, false},
		{"readOnly false wins over legacy false", `,"readOnly":false,"overridable":false`, false},
		{"readOnly true wins over legacy true", `,"overridable":true,"readOnly":true`, true},
		{"readOnly null falls back to legacy", `,"readOnly":null,"overridable":false`, true},
		{"last duplicate wins", `,"readOnly":true,"readOnly":false`, false},
		{"key case folded", `,"READONLY":true`, true},
		{"legacy key case folded", `,"Overridable":false`, true},
		{"whitespace around the flag", ",\r\n\t\"readOnly\" :\r\n\ttrue\r\n", true},
		{"flag inside the value is not the flag", `,"value":{"readOnly":true}`, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var attr data.Attribute
			err := json.Unmarshal([]byte(`{"name":"p","type":"any"`+tt.flags+`}`), &attr)
			if assert.NoError(t, err) {
				assert.Equal(t, tt.readOnly, attr.ReadOnly())
				assert.Equal(t, !tt.readOnly, attr.Overridable())
			}
		})
	}
}

func TestAttribute_UnmarshalRejectsNonBooleanFlag(t *testing.T) {
	for _, flags := range []string{`"readOnly":"true"`, `"readOnly":1`, `"readOnly":{}`, `"readOnly":[]`,
		`"overridable":"false"`, `"overridable":0`} {
		var attr data.Attribute
		err := json.Unmarshal([]byte(`{"name":"p","type":"string",`+flags+`}`), &attr)
		if assert.Error(t, err, flags) {
			assert.Contains(t, err.Error(), `attribute "p"`, flags)
			var typeErr *json.UnmarshalTypeError
			assert.True(t, errors.As(err, &typeErr), flags)
		}
	}
}

func TestAttribute_UnmarshalResetsReadOnly(t *testing.T) {
	var attr data.Attribute
	assert.NoError(t, json.Unmarshal([]byte(`{"name":"p","type":"string","readOnly":true}`), &attr))
	assert.NoError(t, json.Unmarshal([]byte(`{"name":"p","type":"string"}`), &attr))
	assert.False(t, attr.ReadOnly())
}

func TestAttribute_MarshalReadOnly(t *testing.T) {
	tests := []struct {
		name string
		attr *data.Attribute
		want string
	}{
		{"editable writes no flag", data.NewAttribute("p", data.TypeString, "v"), `{"name":"p","type":"string","value":"v"}`},
		{"read-only writes readOnly", data.NewAttributeWithReadOnly("p", data.TypeString, true, "v"), `{"name":"p","type":"string","readOnly":true,"value":"v"}`},
		{"explicit editable writes no flag", data.NewAttributeWithReadOnly("p", data.TypeInt, false, 1), `{"name":"p","type":"int","value":1}`},
		{"deprecated constructor, locked", data.NewAttributeWithOverridable("p", data.TypeString, false, "v"), `{"name":"p","type":"string","readOnly":true,"value":"v"}`},
		{"deprecated constructor, editable", data.NewAttributeWithOverridable("p", data.TypeString, true, "v"), `{"name":"p","type":"string","value":"v"}`},
		{"schema attribute is editable", data.NewAttributeWithSchema("p", data.TypeObject, nil, nil), `{"name":"p","type":"object"}`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := json.Marshal(tt.attr)
			if assert.NoError(t, err) {
				assert.Equal(t, tt.want, string(b))
			}
		})
	}
}

func TestAttribute_LegacyLockSurvivesRoundTrip(t *testing.T) {
	var in data.Attribute
	assert.NoError(t, json.Unmarshal([]byte(`{"name":"p","type":"string","overridable":false,"value":"v"}`), &in))
	b, err := json.Marshal(&in)
	assert.NoError(t, err)
	assert.Equal(t, `{"name":"p","type":"string","readOnly":true,"value":"v"}`, string(b))
	var out data.Attribute
	assert.NoError(t, json.Unmarshal(b, &out))
	assert.True(t, out.ReadOnly())
}

// lockMembers reads an object's members token by token, independently of struct tags, and returns the last value of
// each lock key matched case-insensitively. ok is false for anything but one object, or a key that is not ASCII.
func lockMembers(in []byte) (readOnly, overridable json.RawMessage, ok bool) {
	dec := json.NewDecoder(bytes.NewReader(in))
	if tok, err := dec.Token(); err != nil || tok != json.Delim('{') {
		return nil, nil, false
	}
	for dec.More() {
		tok, err := dec.Token()
		key, isKey := tok.(string)
		if err != nil || !isKey {
			return nil, nil, false
		}
		for i := 0; i < len(key); i++ {
			if key[i] >= 0x80 {
				return nil, nil, false
			}
		}
		var value json.RawMessage
		if err := dec.Decode(&value); err != nil {
			return nil, nil, false
		}
		switch {
		case strings.EqualFold(key, "readOnly"):
			readOnly = value
		case strings.EqualFold(key, "overridable"):
			overridable = value
		}
	}
	return readOnly, overridable, true
}

// FuzzAttribute_ReadOnly checks the decoded lock against a token-level reading of the input that ignores struct tags:
// readOnly decides, else only overridable: false locks; null counts as absent.
func FuzzAttribute_ReadOnly(f *testing.F) {
	for _, seed := range []string{`{"name":"p","type":"string"}`, `{"name":"p","type":"string","readOnly":true}`,
		`{"name":"p","type":"string","overridable":false}`, `{"name":"p","type":"any","readOnly":null,"overridable":false}`,
		`{"name":"p","type":"any","READONLY":false,"overridable":false}`, `{"name":"p","type":"string","readOnly":"true"}`,
		`{"name":"p","type":"string","readOnly":true,"readOnly":null}`, `{"name":"p","type":"string","readOnly":true,"ReadOnly":null}`,
		`{"name":"p","type":"string","readOnly":true,"READONLY":false}`, `{"name":"p","type":"string","readOnly":"x","readOnly":true}`} {
		f.Add([]byte(seed))
	}
	f.Fuzz(func(t *testing.T, in []byte) {
		readOnly, overridable, ok := lockMembers(in)
		var attr data.Attribute
		err := json.Unmarshal(in, &attr)
		if !ok {
			return
		}
		bad := func(v json.RawMessage) bool {
			s := string(bytes.TrimSpace(v))
			return v != nil && s != "true" && s != "false" && s != "null"
		}
		if bad(readOnly) || bad(overridable) {
			if err == nil {
				t.Fatalf("%q: a non-boolean lock decodes", in)
			}
			return
		}
		if err != nil {
			return
		}
		want := string(bytes.TrimSpace(overridable)) == "false"
		if s := string(bytes.TrimSpace(readOnly)); s == "true" || s == "false" {
			want = s == "true"
		}
		if attr.ReadOnly() != want {
			t.Fatalf("%q: ReadOnly() = %v, the token scan reads %v", in, attr.ReadOnly(), want)
		}
		b, err := json.Marshal(&attr)
		if err != nil {
			return
		}
		var keys map[string]json.RawMessage
		if err := json.Unmarshal(b, &keys); err != nil {
			t.Fatalf("%q marshals to %s, which is not an object: %v", in, b, err)
		}
		for k := range keys {
			if strings.EqualFold(k, "overridable") {
				t.Fatalf("%q marshals to %s, which writes the legacy key", in, b)
			}
		}
		var back data.Attribute
		if err := json.Unmarshal(b, &back); err == nil && back.ReadOnly() != attr.ReadOnly() {
			t.Fatalf("%q: round trip %s reads ReadOnly() = %v, want %v", in, b, back.ReadOnly(), attr.ReadOnly())
		}
	})
}
