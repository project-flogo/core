package support

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHandlePanic(t *testing.T) {

}



func TestURLStringToFilePath(t *testing.T) {

	url1 := "/test/path"
	p, fixed := URLStringToFilePath(url1)
	assert.False(t, fixed)
	assert.Empty(t, p)

	url1 = "file:///test/path"
	p, fixed = URLStringToFilePath(url1)
	assert.True(t, fixed)
	assert.Equal(t, "/test/path",p)

	url1 = "file:///test/path%20with%20space"
	p, fixed = URLStringToFilePath(url1)
	assert.True(t, fixed)
	assert.Equal(t, "/test/path with space",p)
}
