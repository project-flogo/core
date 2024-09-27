package spec

var toResolve []*specHolder

func ResolveSpecs() {
	for _, sh := range toResolve {
		s := Get(sh.id)
		if s == nil {
			sh.schema = emptySpec
		} else {
			sh.schema = s
		}
	}
}

type specHolder struct {
	id     string
	schema Spec
}

func (sh *specHolder) Type() string {
	return sh.schema.Type()
}

func (sh *specHolder) Value() string {
	return sh.schema.Value()
}

func (sh *specHolder) Name() string {
	return sh.schema.Name()
}
