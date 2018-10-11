package data

type StringsMap interface {
	Get(key string) string

	Iterate(itx func(key string, value string))
}

func NewFixedStringsMap(params map[string]string) StringsMap {

	fp := &stringsMapImpl{}
	fp.ssMap = make(map[string]string)

	for key, value := range params {
		fp.ssMap[key] = value
	}

	return fp
}

type stringsMapImpl struct {
	ssMap map[string]string
}

func (d *stringsMapImpl) Get(key string) string {
	return d.ssMap[key]
}

func (d *stringsMapImpl) Iterate(itx func(string, string)) {
	for key, value := range d.ssMap {
		itx(key, value)
	}
}

type StructValue interface {
	ToMap() map[string]interface{}
	FromMap(values map[string]interface{}) error
}
