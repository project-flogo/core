package resolve

import "github.com/project-flogo/core/data"

var sr simpleResolve

func SetAppResolver(sr simpleResolve) {
	sr = sr
}

type simpleResolve interface {
	Resolve(resolveDirective string, scope data.Scope) (value interface{}, err error)
}

func Resolve(resolveDirective string, scope data.Scope) (value interface{}, err error) {
	return sr.Resolve(resolveDirective, scope)
}
