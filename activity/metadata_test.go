package activity

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

var mdJSON = `
	{
		"settings" : [
			{
				"name":"a",
				"type":"string"
			}
		],
		"input" : [
			{
				"name":"a",
				"type":"string"
			}
		],
		"output" : [
			{
				"name":"a",
				"type":"string"
			}
		]	
	}
`
var mdErrJSON = `
{
	"input" : {
		"name":"a",
		"type":"string"
	}
	
}
`

func TestToMetadata(t *testing.T) {

	md := ToMetadata(&Settings{}, &Input{}, &Output{})

	assert.NotNil(t, md)

}
func TestUnMarshalJSON(t *testing.T) {

	md := ToMetadata(&Settings{}, &Input{}, &Output{})

	err := md.UnmarshalJSON([]byte(mdJSON))

	assert.Nil(t, err)

	err = md.UnmarshalJSON([]byte(mdErrJSON))

	assert.NotNil(t, err)
}
