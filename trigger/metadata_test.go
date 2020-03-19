package trigger

import (
	"testing"

	"github.com/project-flogo/core/data"
	"github.com/stretchr/testify/assert"
)

type TVal struct {
	val interface{}
}

func (t TVal) Type() data.Type {
	return data.TypeString
}
func (t TVal) Value() interface{} {
	return "sample"
}

func TestMetadata(t *testing.T) {
	var tval TVal
	md := &Metadata{Settings: map[string]data.TypedValue{"settings": tval}, HandlerSettings: map[string]data.TypedValue{"input": tval},
		Output: map[string]data.TypedValue{"output": tval}, Reply: map[string]data.TypedValue{"reply": tval}}

	assert.NotNil(t, NewMetadata(md))
}
