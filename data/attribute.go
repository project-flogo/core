package data

import (
	"encoding/json"
	"github.com/project-flogo/core/data/schema"
)

// TODO this should probably go away
func NewAttribute(name string, dataType Type, value interface{}) *Attribute {
	return &Attribute{name: name, dataType: dataType, value: value}
}

// Attribute is a simple structure used to define a data Attribute/property
type Attribute struct {
	name     string
	dataType Type
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

//todo temporary work around for attribute deserialization

// UnmarshalJSON implements json.Unmarshaler.UnmarshalJSON
func (a *Attribute) UnmarshalJSON(data []byte) error {

	ser := &struct {
		Name   string      `json:"name"`
		Type   string      `json:"type"`
		Value  interface{} `json:"value"`
		Schema interface{} `json:"schema,omitempty"`

		//KeyType  string      `json:"keyType,omitempty"`
		//ElemType string      `json:"elemType,omitempty"`
	}{}

	if err := json.Unmarshal(data, ser); err != nil {
		return err
	}
	a.name = ser.Name

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

func getSchema(schemaDef interface{}) (schema.Schema, error) {
	return schema.FindOrCreate(schemaDef)
}
