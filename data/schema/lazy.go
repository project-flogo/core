package schema

var toResolve []*schemaHolder

func ResolveSchemas() {
	for _, sh := range toResolve {
		s := Get(sh.id)
		if s == nil {
			sh.schema = emptySchema
		} else {
			sh.schema = s
		}
	}
}

type schemaHolder struct {
	id     string
	schema Schema
}

func (sh *schemaHolder) Type() string {
	return sh.schema.Type()
}

func (sh *schemaHolder) Value() string {
	return sh.schema.Value()
}

func (sh *schemaHolder) Validate(data interface{}) error {
	return sh.schema.Validate(data)
}
