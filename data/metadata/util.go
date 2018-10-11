package metadata

import (
	"fmt"
	"github.com/project-flogo/core/data/coerce"
	"github.com/project-flogo/core/data/resolve"
	"github.com/project-flogo/core/support/logger"
	"reflect"

	"github.com/project-flogo/core/data"
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

				fv.Set(reflect.ValueOf(val))
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

		//tv := data.NewTypedValue(details.Type, fv.Interface())
		values[details.Label] = fv.Interface()
	}

	return values
}

func ResolveSettingValue(setting string, value interface{}, settingsMd map[string]data.TypedValue) (interface{}, error) {

	strVal, ok := value.(string)

	toType := data.TypeUnknown

	if settingsMd != nil {
		tv := settingsMd[setting]
		if tv != nil {
			toType = tv.Type()
		}
	}

	if ok && len(strVal) > 0 && strVal[0] == '$' {
		v, err := resolve.GetBasicResolver().Resolve(strVal, nil)
		if err == nil {

			v, err = coerce.ToType(v, toType)
			if err != nil {
				return nil, err
			}

			logger.Debugf("Resolved setting [%s: %s] to : %v", setting, value, v)
			return v, nil
		}
	}

	return coerce.ToType(value, toType)
}
