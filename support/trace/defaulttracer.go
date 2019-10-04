package trace

type nooptracer struct{}

type dummycontext struct{}


func (nt *nooptracer) Name() string {
	return "noop-tracer"
}


func (nt *nooptracer) Extract(format CarrierFormat, data interface{}) (TracingContext, error) {
	return dummycontext{}, nil
}

func (nt *nooptracer) StartTrace(config Config, parent TracingContext) (TracingContext, error) {
	return dummycontext{}, nil
}
func (nt *nooptracer) FinishTrace(tContext TracingContext, err error) error {
	return nil
}
func (nt *nooptracer) Start() error {
	return nil
}
func (nt *nooptracer) Stop() error {
	return nil
}

func (nt dummycontext) SetTag(TagKey string, TagValue interface{}) bool {
	return false
}

func (nt dummycontext) SetTags( map[string]interface{}) bool {
	return false
}
func (nt dummycontext) LogKV(kvs  map[string]interface{}) bool {
	return false
}

func (tc dummycontext) TraceObject() interface{} {
	return nil
}

func (nt dummycontext) Inject(format CarrierFormat, carrier interface{}) error {
	return nil
}
