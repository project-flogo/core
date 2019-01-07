package propertyresolver

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnvValueResolver(t *testing.T) {
	os.Setenv("Test", "Test")
	os.Setenv("TEST_PROP", "test.Prop")

	defer func() {
		os.Unsetenv("Test")
		os.Unsetenv("TEST_PROP")
	}()
	resolver := &EnvVariableValueResolver{}

	resolvedVal, found := resolver.LookupValue("Test")
	assert.True(t, true, found)
	assert.Equal(t, "Test", resolvedVal)

	resolvedVal, found = resolver.LookupValue("test.Prop")
	assert.True(t, true, found)
	assert.Equal(t, "test.Prop", resolvedVal)

}
