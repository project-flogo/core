package metadata

import (
	"fmt"
	"reflect"

	"github.com/project-flogo/core/data"
	"github.com/project-flogo/core/data/coerce"
	"github.com/project-flogo/core/data/expression"
)

const metadataTag = "md"

func StructToTypedMap(object interface{}) map[string]data.TypedValue {

	if object == nil {
		return nil
	}

	v := reflect.ValueOf(object)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	values := make(map[string]data.TypedValue, v.NumField())

	for i := 0; i < v.NumField(); i++ {
		fv := v.Field(i)
		ft := v.Type().Field(i)

		tag := ft.Tag.Get(metadataTag)

		details := NewFieldDetails(ft.Name, fv.Type().String(), tag)

		tv := data.NewTypedValue(details.Type, fv.Interface())
		values[details.Label] = tv
	}

	return values
}

func TypedMapToStruct(m map[string]data.TypedValue, object interface{}, validate bool) error {

	v := reflect.ValueOf(object)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	for i := 0; i < v.NumField(); i++ {
		fv := v.Field(i)
		ft := v.Type().Field(i)

		if fv.CanSet() {
			tag := ft.Tag.Get(metadataTag)
			details := NewFieldDetails(ft.Name, fv.Type().String(), tag)

			mapVal, ok := m[details.Label]

			if ok {

				if validate {
					err := details.Validate(mapVal.Value())
					if err != nil {
						return err
					}
				}

				fv.Set(reflect.ValueOf(mapVal.Value()))
			} else {
				if validate && details.Required {
					return fmt.Errorf("field '%s' is required", details.Label)
				}
			}
		}
	}

	return nil
}

func MapToStruct(m map[string]interface{}, object interface{}, validate bool) error {

	v := reflect.ValueOf(object)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	for i := 0; i < v.NumField(); i++ {
		fv := v.Field(i)
		ft := v.Type().Field(i)

		if fv.CanSet() {
			tag := ft.Tag.Get(metadataTag)
			details := NewFieldDetails(ft.Name, fv.Type().String(), tag)

			mapVal, ok := m[details.Label]

			if ok {

				if validate {
					err := details.Validate(mapVal)
					if err != nil {
						return err
					}
				}

				val, err := coerce.ToType(mapVal, details.Type)
				if err != nil {
					return err
				}

				if IsZeroOfUnderlyingType(val) {
					fv.Set(reflect.Zero(fv.Type()))
				} else {
					fv.Set(reflect.ValueOf(val))
				}
			} else {
				if validate && details.Required {
					return fmt.Errorf("field '%s' is required", details.Label)
				}
			}
		}
	}

	return nil
}

func StructToMap(object interface{}) map[string]interface{} {

	if object == nil {
		return nil
	}

	v := reflect.ValueOf(object)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	values := make(map[string]interface{}, v.NumField())

	for i := 0; i < v.NumField(); i++ {
		fv := v.Field(i)
		ft := v.Type().Field(i)

		tag := ft.Tag.Get(metadataTag)
		details := NewFieldDetails(ft.Name, fv.Type().String(), tag)

		values[details.Label] = fv.Interface()
	}

	return values
}

func ResolveSettingValue(setting string, value interface{}, settingsMd map[string]data.TypedValue, ef expression.Factory) (interface{}, error) {

	strVal, ok := value.(string)

	toType := data.TypeUnknown

	if settingsMd != nil {
		tv := settingsMd[setting]
		if tv != nil {
			toType = tv.Type()
		}
	}

	if ok && len(strVal) > 0 && strVal[0] == '=' && ef != nil {
		expr, err := ef.NewExpr(strVal[1:])
		if err != nil {
			return nil, err
		}

		value, err = expr.Eval(nil)
		if err != nil {
			return nil, err
		}
	}

	return coerce.ToType(value, toType)
}

func IsZeroOfUnderlyingType(x interface{}) bool {
	if x == nil {
		return true
	}
	typ := reflect.TypeOf(x)
	zero := reflect.Zero(typ).Interface()
	if typ.Comparable() {
		return x == zero
	}
	return reflect.DeepEqual(x, zero)
}
