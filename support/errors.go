package support

type RetriableError struct {
	actName   string
	errorCode string
	errorStr  string
	errorData interface{}
}

func (r *RetriableError) Error() string {

	return r.errorStr
}

func (r *RetriableError) Code() string {

	return r.errorCode
}

func (r *RetriableError) Data() interface{} {

	return r.errorData
}
