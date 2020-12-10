package support

import (
	"reflect"
	"strings"
)

func GetRef(contrib interface{}) string {
	v := reflect.ValueOf(contrib)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	ref := v.Type().PkgPath()

	return fixPkgPathVendoring(ref)
}

type HasRef interface {
	Ref() string
}

// fixes vendored paths
func fixPkgPathVendoring(pkgPath string) string {
	const vendor = "/vendor/"
	if i := strings.LastIndex(pkgPath, vendor); i != -1 {
		return pkgPath[i+len(vendor):]
	}
	return pkgPath
}
