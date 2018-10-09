package resource

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

const resJSON = `
{
  "resources":
  [
    {
      "id": "flow:myflow",
      "data":{
      }
    },
    {
      "id": "schema:myschema",
      "data":{
      }
    },
    {
      "id": "connection:myConnection",
      "data":{
      }
    }
  ]
}
`

func TestDeserialize(t *testing.T) {

	defRep := &ResourcesConfig{}

	err := json.Unmarshal([]byte(resJSON), defRep)
	assert.Nil(t, err)

	fmt.Printf("Resources: %v", defRep)
}
