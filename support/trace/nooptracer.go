package trace

type nooptracer struct{}

func (nt *nooptracer) Name() string     { return "noop-tracer" }
func (nt *nooptracer) configure() error { return nil }
func (nt *nooptracer) Inject(flowID string, taskInstID string, format CarrierFormat, carrier interface{}) error { return nil }
func (nt *nooptracer) Extract(format CarrierFormat, data interface{}) (interface{}, error) { return nil, nil }
func (nt *nooptracer) SetTag(spanKey string, TagKey string, TagValue interface{}) bool { return false }
func (nt *nooptracer) LogKV(spanKey string, alternatingKeyValues ...interface{}) bool  { return false }
