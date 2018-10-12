package metadata

import (
	"fmt"
	"strings"

	"github.com/project-flogo/core/data"
	"github.com/project-flogo/core/data/coerce"
)

type FieldDetails struct {
	Name     string
	Label    string
	Type     data.Type
	Allowed  []string
	Required bool
}

func (d *FieldDetails) AllowedToString() string {

	if len(d.Allowed) == 0 {
		return "null"
	}

	vals := make([]string, len(d.Allowed))
	for i, value := range d.Allowed {
		s := "\"" + value + "\""
		vals[i] = s
	}

	return fmt.Sprintf("[%s]", strings.Join(vals, ","))
}

func (d *FieldDetails) Validate(value interface{}) error {
	valid := true

	if d.Required {
		if value == nil || value == "" {
			valid = false
		}
	}

	if len(d.Allowed) > 0 {

		valid = false
		for _, av := range d.Allowed {
			//todo handler error
			allowedValue, _ := coerce.ToType(av, d.Type)
			if d.Type == data.TypeString {
				strVal, ok := value.(string)
				if ok && strings.EqualFold(strVal, av) {
					valid = true
					break
				}
			} else {
				if value == allowedValue {
					valid = true
					break
				}
			}
		}
	}

	if !valid {
		return fmt.Errorf("value '%v' is not valid", value)
	}

	return nil
}

func NewFieldDetails(name string, dType string, mdTag string) *FieldDetails {

	details := &FieldDetails{}
	details.Name = name
	details.Type = data.ToTypeFromGoRep(dType)

	if len(mdTag) == 0 {
		details.Label = name
		return details
	}

	components := deconstructTag(mdTag)

	if components[0] == "-" {
		return nil
	}

	details.Label = components[0]

	if len(components) > 0 {
		for i := 0; i < len(components); i++ {
			applyTagComponent(details, components[i])
		}
	}

	return details
}

func applyTagComponent(details *FieldDetails, component string) {

	//process sets
	if strings.HasPrefix(component, "allowed(") {
		values := component[8 : len(component)-1]
		details.Allowed = strings.Split(values, ",")
		return
	}

	//process flags
	switch component {
	case "required":
		details.Required = true
	}
}

func deconstructTag(str string) []string {

	var parts []string

	start := 0
	ignore := false
	for i := 0; i < len(str); i++ {

		switch str[i] {
		case '(':
			ignore = true
		case ')':
			ignore = false
		case ',':
			if !ignore {
				parts = append(parts, str[start:i])
				start = i + 1
			}
		}
	}

	parts = append(parts, str[start:])

	return parts
}
