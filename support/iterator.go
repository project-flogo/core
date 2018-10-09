package support

type Iterator interface {
	Next() interface{}
	HasNext() bool
}

//type MySlice []interface{}

type sliceIterator struct {
	slice []interface{}
	index int
}

func (i *sliceIterator) Next() interface{} {
	i.index++
	return i.slice[i.index-1]
}

func (i *sliceIterator) HasNext() bool {
	return i.index < len(i.slice)
}

//func (s *MySlice) GetIterator() Iterator {
//	return &sliceIterator{s, 0}
//}

type FixedDetails struct {
	data map[string]string
}

func (d *FixedDetails) Get(key string) string {
	return d.data[key]
}

func (d *FixedDetails) Iterate(itx func(string, string)) {
	for key, value := range d.data {
		itx(key, value)
	}
}
