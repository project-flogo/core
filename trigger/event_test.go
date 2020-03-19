package trigger

import (
	"os"
	"testing"
	"time"

	"github.com/project-flogo/core/engine/event"
	"github.com/stretchr/testify/assert"
)

var sampleString string

type SampleListener struct {
}

func (sl *SampleListener) HandleEvent(ctx *event.Context) error {
	sampleString = "TriggerListened"
	return nil
}

func TestTriggerEvent(t *testing.T) {
	os.Setenv("FLOGO_PUBLISH_AUDIT_EVENTS", "true")
	var err error
	err = event.RegisterListener("sample", &SampleListener{}, []string{"triggerevent"})
	assert.Nil(t, err)
	PostTriggerEvent("Started", "name")
	time.Sleep(2 * time.Millisecond)
	assert.Equal(t, "TriggerListened", sampleString)

}
