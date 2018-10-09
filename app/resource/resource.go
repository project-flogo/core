package resource

type Resource struct {
	resType string
	resObj  interface{}
}

func New(resType string, resObj interface{}) *Resource {
	return &Resource{resType: resType, resObj: resObj}
}

func (r *Resource) Type() string {
	return r.resType
}

func (r *Resource) Object() interface{} {
	return r.resObj
}
