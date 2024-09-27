package spec

import (
	"fmt"
)

func New(specDef *Def) (Spec, error) {
	s := defaultSpecImpl{def: specDef}
	return s, nil
}

type defaultSpecImpl struct {
	def *Def
}

// Returns the name of spec
func (s defaultSpecImpl) Name() string {
	return s.def.Name
}

// Returns the type of spec
func (s defaultSpecImpl) Type() string {
	return s.def.Type
}

// Returns content of spec
func (s defaultSpecImpl) Value() string {
	return s.def.Value
}

func CreateSpec(specRep interface{}) (Spec, error) {

	switch t := specRep.(type) {
	case map[string]string:
		def := &Def{}
		if sName, ok := t["name"]; ok {
			def.Name = sName
		} else {
			if sName, ok := t["filename"]; ok {
				def.Name = sName
			} else {
				return nil, fmt.Errorf("invalid spec definition, name not specified: %+v", t)
			}
		}

		if sType, ok := t["type"]; ok {
			def.Type = sType
		} else {
			return nil, fmt.Errorf("invalid schema definition, type not specified: %+v", t)
		}

		if sValue, ok := t["content"]; ok {
			def.Value = sValue
		} else {
			return nil, fmt.Errorf("invalid schema definition, value not specified: %+v", t)
		}
		return New(def)
	case map[string]interface{}:
		def := &Def{}
		if sName, ok := t["name"]; ok {
			def.Name = sName.(string)
		} else {
			if sName, ok := t["filename"]; ok {
				def.Name = sName.(string)
			} else {
				return nil, fmt.Errorf("invalid spec definition, name not specified: %+v", t)
			}
		}

		if sValue, ok := t["content"]; ok {
			def.Value, ok = sValue.(string)
			if !ok {
				return nil, fmt.Errorf("invalid schema definition, value is not a string specified: %+v", sValue)
			}
		} else {
			return nil, fmt.Errorf("invalid schema definition, value not specified: %+v", t)
		}
		return New(def)
	case nil:
		return nil, nil
	default:
		return nil, fmt.Errorf("invalid spec definition, %v", t)
	}
}

var emptySpec = &emptySpecImpl{}

type emptySpecImpl struct {
}

// Name implements Spec.
func (*emptySpecImpl) Name() string {
	return ""
}

func (*emptySpecImpl) Type() string {
	return ""
}

func (*emptySpecImpl) Value() string {
	return ""
}
