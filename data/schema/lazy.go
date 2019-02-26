package schema

var toResolve []*schemaHolder
var emptySchema = &emptySchemaImpl{}

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

type emptySchemaImpl struct {
}

func (*emptySchemaImpl) Type() string {
	return ""
}

func (*emptySchemaImpl) Value() string {
	return ""
}

func (*emptySchemaImpl) Validate(data interface{}) error {
	return nil
}
