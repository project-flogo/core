package propertyresolver

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnvValueResolver(t *testing.T) {
	_ = os.Setenv("Test", "Test")
	_ = os.Setenv("TEST_PROP", "test.Prop")

	defer func() {
		_ = os.Unsetenv("Test")
		_ = os.Unsetenv("TEST_PROP")
	}()
	resolver := &EnvVariableValueResolver{autoMapping: true}

	resolvedVal, found := resolver.LookupValue("Test")
	assert.True(t, true, found)
	assert.Equal(t, "Test", resolvedVal)

	resolvedVal, found = resolver.LookupValue("test.Prop")
	assert.True(t, true, found)
	assert.Equal(t, "test.Prop", resolvedVal)

}
