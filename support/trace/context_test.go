package trace

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

type TestTracingContext struct {

}

func (t TestTracingContext) TraceObject() interface{} {
	return nil
}

func (t TestTracingContext) SetTags(tags map[string]interface{}) bool {
	return false
}

func (t TestTracingContext) SetTag(tagKey string, tagValue interface{}) bool {
	return false
}

func (t TestTracingContext) LogKV(kvs map[string]interface{}) bool {
	return false
}

func TestAppendTracingContext(t *testing.T) {

	tCtx := &TestTracingContext{}

	goCtx := AppendTracingContext(context.Background(), tCtx)
 	tc, ok := goCtx.Value(id).(TracingContext)

 	assert.True(t, ok)
 	assert.Equal(t, tCtx, tc)
}

func TestExtractTracingContext(t *testing.T) {

	tCtx := &TestTracingContext{}
	goCtx := context.WithValue(nil, id, tCtx)

	ttCtx := ExtractTracingContext(goCtx)
	assert.Equal(t, tCtx, ttCtx)

	ttCtx = ExtractTracingContext(context.Background())
	assert.Nil(t, ttCtx)
}