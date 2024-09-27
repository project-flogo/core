package coerce

import (
	"fmt"
	"strings"

	"github.com/project-flogo/core/support/spec"
)

func ToSpec(val interface{}) (spec.Spec, error) {
	switch t := val.(type) {
	case string:
		if strings.HasPrefix(t, "spec://") {
			id := t[7:]
			s := spec.Get(id)
			if s == nil {
				return nil, fmt.Errorf("spec with id '%s' not found", t)
			}
			return s, nil
		} else {
			// Attempt to parse the string as JSON
			specData, err := ToObject(t)
			if err != nil {
				return nil, err
			}
			s, err := spec.CreateSpec(specData)
			if err != nil {
				return nil, err
			}
			return s, nil
		}
	case spec.Spec:
		return t, nil
	case map[string]string, map[string]interface{}:
		s, err := spec.CreateSpec(t)
		if err != nil {
			return nil, err
		}
		return s, nil
	default:
		// try to create config from map[string]interface{}
		return nil, fmt.Errorf("unable to create spec from '%#v'", val)
	}
}
