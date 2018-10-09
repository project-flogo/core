package support

import "reflect"

func GetRef(contrib interface{}) string {
	v := reflect.ValueOf(contrib)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	ref := v.Type().PkgPath()

	return ref
}

type HasRef interface {
	Ref() string
}
