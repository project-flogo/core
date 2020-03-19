package trace

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type TestTracer struct {

}

func (t TestTracer) Start() error {
	return nil
}

func (t TestTracer) Stop() error {
	return nil

}

func (t TestTracer) Name() string {
	return "test tracer"
}

func (t TestTracer) Extract(format CarrierFormat, carrier interface{}) (TracingContext, error) {
	return nil, nil
}

func (t TestTracer) Inject(tCtx TracingContext, format CarrierFormat, carrier interface{}) error {
	return nil
}

func (t TestTracer) StartTrace(config Config, parent TracingContext) (TracingContext, error) {
	return nil, nil

}

func (t TestTracer) FinishTrace(tContext TracingContext, err error) error {
	return nil
}

func TestEnabled(t *testing.T) {

	tracer = nil
	assert.False(t, Enabled())
	tracer = &TestTracer{}
	assert.True(t, Enabled())
}


func TestGetTracer(t *testing.T) {

	tracer = nil
	assert.Nil(t, GetTracer())

	tr := &TestTracer{}
	tracer = tr

	ttr := GetTracer()
	assert.Equal(t, tr, ttr)
}

func TestRegisterTracer(t *testing.T) {

	tracer = nil
	tr := &TestTracer{}

	err := RegisterTracer(tr)
	assert.Nil(t, err)
	assert.Equal(t, tr, tracer)

	tr2 := &TestTracer{}
	err = RegisterTracer(tr2)
	assert.NotNil(t, err)
}