package event

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var sampleString string

type SampleListener struct {
}

func (sl *SampleListener) HandleEvent(ctx *Context) error {
	sampleString = "Listened"
	return nil
}

func TestRegisterListener(t *testing.T) {
	os.Setenv("FLOGO_PUBLISH_AUDIT_EVENTS", "true")
	defer func() {
		os.Unsetenv("FLOGO_PUBLISH_AUDIT_EVENTS")
	}()
	var err error
	err = RegisterListener("sample", &SampleListener{}, []string{"sample"})
	assert.Nil(t, err)

	err = RegisterListener("sample", nil, []string{"sample"})
	assert.NotNil(t, err)

	err = RegisterListener("", nil, []string{"sample"})
	assert.NotNil(t, err)
}

func TestPost(t *testing.T) {
	ctx := &Context{eventType: "sample", event: 2}
	Post("sample", ctx)
	time.Sleep(2 * time.Millisecond)
	assert.Equal(t, "Listened", sampleString)
	UnRegisterListener("sample", []string{"sample"})
}
