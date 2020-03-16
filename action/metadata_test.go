package action

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type Settings struct {
	Name string `md:"name"`
}
type Input struct {
	Name string `md:"name"`
}
type Output struct {
	Name string `md:"name"`
}

func TestToMetadata(t *testing.T) {

	md := ToMetadata(&Settings{}, &Input{}, &Output{})

	assert.NotNil(t, md)

}
func TestMarshalJSON(t *testing.T) {

	md := ToMetadata(&Settings{}, &Input{}, &Output{})

	_, err := md.MarshalJSON()

	assert.Nil(t, err)

}
