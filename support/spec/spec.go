package spec

type Spec interface {
	Type() string
	Value() string
	Name() string
}

type Def struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Value string `json:"content"`
}
