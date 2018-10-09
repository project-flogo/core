package expression

import (
	"fmt"
	"os"
	"testing"

	"github.com/project-flogo/core/data"
	"github.com/project-flogo/core/data/resolve"
	"github.com/stretchr/testify/assert"
)

var resolver = resolve.NewCompositeResolver(map[string]resolve.Resolver{"env": &resolve.EnvResolver{}, ".": &TestResolver{"."}})

type TestResolver struct {
	retVal string
}

func (*TestResolver) GetResolverInfo() *resolve.ResolverInfo {
	return resolve.NewResolverInfo(false, false)
}

//EnvResolver Environment Resolver $env[item]
func (*TestResolver) Resolve(scope data.Scope, item string, field string) (interface{}, error) {
	// Environment resolution
	value, exists := os.LookupEnv(item)
	if !exists {
		err := fmt.Errorf("failed to resolve Environment Variable: '%s', ensure that variable is configured", item)
		return "", err
	}

	return value, nil
}

func TestNewExpressionFactory(t *testing.T) {
	f := NewFactory(resolver)
	assert.NotNil(t, f)
}

func TestNewExpression(t *testing.T) {
	f := NewFactory(resolver)
	assert.NotNil(t, f)

	e, err := f.NewExpr("1234")
	assert.Nil(t, err)
	_, ok := e.(*literalExpr)
	assert.True(t, ok)

	v, err := e.Eval(nil)
	assert.Nil(t, err)
	assert.Equal(t, 1234, v)

	e, err = f.NewExpr("$env[PATH]")
	assert.Nil(t, err)
	_, ok = e.(*literalExpr)
	assert.True(t, ok)

	v, err = e.Eval(nil)
	assert.Nil(t, err)
	_, ok = v.(string)
	assert.True(t, ok)

	//e, err = f.NewExpr("$test")
	//assert.Nil(t, err)
	//_, ok = e.(*resolveExpr)
	//assert.True(t, ok)

	e, err = f.NewExpr("$.test")
	assert.Nil(t, err)
	_, ok = e.(*resolutionExpr)
	assert.True(t, ok)

	e, err = f.NewExpr("{{\"test\"}}")
	assert.Nil(t, err)
	_, ok = e.(*templateExpr)
	assert.True(t, ok)

	v, err = e.Eval(nil)
	assert.Nil(t, err)
	s, ok := v.(string)
	assert.Equal(t, "test", s)
}
