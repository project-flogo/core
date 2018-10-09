package api

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestToMappings(t *testing.T) {

	mappings := []string{"in1=b", "in2= $.blah", "in3 = $.blah2"}

	//todo add additional tests when support for more mapping type is added
	defs, err := toMappings(mappings)

	assert.Nil(t, err)
	assert.Equal(t, 3, len(defs))

	v,exists := defs["in1"]
	assert.True(t, exists)
	assert.Equal(t, "=b",v)

	v,exists = defs["in2"]
	assert.True(t, exists)
	assert.Equal(t, "= $.blah",v)

	v,exists = defs["in3"]
	assert.True(t, exists)
	assert.Equal(t, "= $.blah2",v)

}
