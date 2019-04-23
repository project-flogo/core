package data

import "github.com/project-flogo/core/data/schema"

type TypedValue interface {
	Type() Type
	Value() interface{}
}

//type CompoundValue interface {
//	TypedValue
//	KeyType() Type
//	ElemType() Type
//}

//func ToCompoundValue(value TypedValue) (CompoundValue, bool) {
//	if value.Type() == TypeArray || value.Type() == TypeMap {
//		ctv, ok := value.(CompoundValue)
//		return ctv, ok
//	}
//
//	return nil, false
//}

func GetSchema(value TypedValue) (schema.Schema, bool) {
	if value.Type() == TypeObject || value.Type() == TypeArray || value.Type() == TypeMap {
		ws, ok := value.(schema.HasSchema)
		return ws.Schema(), ok
	}

	return nil, false
}

func ToTypedValue(value interface{}) TypedValue {
	dataType, err := GetType(value)

	if err != nil {
		dataType = TypeAny
	}

	return &valueImpl{dataType: dataType, value: value}
}

func NewTypedValueWithConversion(dataType Type, value interface{}) (TypedValue, error) {
	newVal, err := typeConverter(value, dataType)
	if err != nil {
		return nil, err
	}
	return &valueImpl{dataType: dataType, value: newVal}, nil
}

func NewTypedValue(dataType Type, value interface{}) TypedValue {
	return &valueImpl{dataType: dataType, value: value}
}

type valueImpl struct {
	dataType Type
	value    interface{}
}

func (v *valueImpl) Type() Type {
	return v.dataType
}

func (v *valueImpl) Value() interface{} {
	return v.value
}

func NewTypedValueFromAttr(attr *Attribute) TypedValue {
	if attr.Type().IsSimple() {
		return &refTypedValue{attr: attr}
	} else {
		return &refCompoundValue{attr: attr}
	}
}

type refTypedValue struct {
	attr  *Attribute
	value interface{}
}

func (ta *refTypedValue) Type() Type {
	return ta.attr.Type()
}

func (ta *refTypedValue) Value() interface{} {
	return ta.value
}

func (ta *refTypedValue) Schema() interface{} {
	return ta.attr.Schema()
}

type refCompoundValue struct {
	attr  *Attribute
	value interface{}
}

func (ta *refCompoundValue) Type() Type {
	return ta.attr.Type()
}

func (ta *refCompoundValue) Value() interface{} {
	//todo get default from attr?
	return ta.value
}

//func (ta *refCompoundValue) KeyType() Type {
//	return ta.attr.KeyType()
//}
//
//func (ta *refCompoundValue) ElemType() Type {
//	return ta.attr.ElemType()
//}

func (ta *refCompoundValue) Schema() schema.Schema {
	return ta.attr.Schema()
}
