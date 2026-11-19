package data

import (
	"encoding/json"
	"fmt"

	"github.com/project-flogo/core/data/schema"
)

func NewAttribute(name string, dataType Type, value interface{}) *Attribute {
	return &Attribute{name: name, dataType: dataType, value: value}
}

// NewAttributeWithReadOnly creates an attribute; a read-only app property keeps the app's value and ignores overrides.
func NewAttributeWithReadOnly(name string, dataType Type, readOnly bool, value interface{}) *Attribute {
	return &Attribute{name: name, dataType: dataType, readOnly: readOnly, value: value}
}

// Deprecated: use NewAttributeWithReadOnly.
func NewAttributeWithOverridable(name string, dataType Type, overridable bool, value interface{}) *Attribute {
	return NewAttributeWithReadOnly(name, dataType, !overridable, value)
}

func NewAttributeWithSchema(name string, dataType Type, value interface{}, schema schema.Schema) *Attribute {
	return &Attribute{name: name, dataType: dataType, value: value, schema: schema}
}

// Attribute is a simple structure used to define a data Attribute/property
type Attribute struct {
	name     string
	dataType Type
	readOnly bool
	value    interface{}

	//keyType  Type
	//elemType Type

	schema schema.Schema
}

func (a *Attribute) Name() string {
	return a.name
}

func (a *Attribute) Type() Type {
	return a.dataType
}

func (a *Attribute) ReadOnly() bool {
	return a.readOnly
}

// Deprecated: use ReadOnly.
func (a *Attribute) Overridable() bool {
	return !a.readOnly
}

func (a *Attribute) Value() interface{} {
	return a.value
}

//func (a *Attribute) KeyType() Type {
//	return a.keyType
//}
//
//func (a *Attribute) ElemType() Type {
//	return a.elemType
//}

func (a *Attribute) Schema() schema.Schema {
	return a.schema
}

func (a *Attribute) MarshalJSON() ([]byte, error) {

	val := a.value
	if vs, ok := val.(string); ok {
		if vs == "" {
			val = nil //hack to omit empty string value
		}
	}

	return json.Marshal(&struct {
		Name     string      `json:"name"`
		Type     string      `json:"type"`
		ReadOnly bool        `json:"readOnly,omitempty"`
		Value    interface{} `json:"value,omitempty"`
		Schema   interface{} `json:"schema,omitempty"`
	}{
		Name:     a.name,
		Type:     a.dataType.String(),
		ReadOnly: a.readOnly,
		Value:    val,
		Schema:   a.schema,
	})
}

// UnmarshalJSON implements json.Unmarshaler.UnmarshalJSON
func (a *Attribute) UnmarshalJSON(data []byte) error {

	ser := &struct {
		Name        string      `json:"name"`
		Type        string      `json:"type"`
		ReadOnly    *bool       `json:"readOnly,omitempty"`
		Overridable *bool       `json:"overridable,omitempty"`
		Value       interface{} `json:"value,omitempty"`
		Schema      interface{} `json:"schema,omitempty"`

		//KeyType  string      `json:"keyType,omitempty"`
		//ElemType string      `json:"elemType,omitempty"`
	}{}

	if err := json.Unmarshal(data, ser); err != nil {
		return attributeError(data, err)
	}
	a.name = ser.Name
	a.readOnly = isReadOnly(ser.ReadOnly, ser.Overridable)

	var err error
	a.schema, err = getSchema(ser.Schema)
	if err != nil {
		return err
	}

	dt, err := ToTypeEnum(ser.Type)
	if err != nil {
		return err
	}
	a.dataType = dt

	//if ser.ElemType != "" {
	//	dt, err := ToTypeEnum(ser.ElemType)
	//	if err != nil {
	//		return err
	//	}
	//	a.elemType = dt
	//}
	//
	//if ser.KeyType != "" {
	//	dt, err := ToTypeEnum(ser.KeyType)
	//	if err != nil {
	//		return err
	//	}
	//	a.keyType = dt
	//}

	val, err := typeConverter(ser.Value, a.dataType)

	if err != nil {
		return err
	} else {
		a.value = val
	}

	return nil
}

// attributeError names the attribute, so a bad entry among an app's properties or metadata can be found.
func attributeError(data []byte, err error) error {
	var named struct {
		Name string `json:"name"`
	}
	_ = json.Unmarshal(data, &named)
	return fmt.Errorf("attribute %q: %w", named.Name, err)
}

// isReadOnly applies the lock rule: readOnly decides; without it, only the legacy overridable: false locks.
func isReadOnly(readOnly, overridable *bool) bool {
	if readOnly != nil {
		return *readOnly
	}
	return overridable != nil && !*overridable
}

func getSchema(schemaDef interface{}) (schema.Schema, error) {
	return schema.FindOrCreate(schemaDef)
}
