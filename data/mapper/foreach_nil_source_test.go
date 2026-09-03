package mapper

import (
	"encoding/json"
	"testing"

	"github.com/project-flogo/core/data"
	"github.com/stretchr/testify/assert"
)

// Regression for FLOGO-19539: @foreach over a source-array path that resolves to
// nil/undefined must be treated as an empty iteration, not a failure.
//
// The customer app (tc8-arraymappings-check) maps `facts` with
// @foreach($activity[input_array].output.attribute, attribute); the upstream
// input_array mapper has only schemas (no input mapping), so the source path is
// unresolvable at runtime and raises "unable to evaluate path: attribute".
//
// This @foreach lives inside a literal-array mapping. ObjectMapper.Eval's
// literal-array branch used to swallow the element error (return err, nil), which
// masked the not-found error and made the array come out empty. FLOGO-19268 fixed
// that swallow to propagate errors, which un-masked the not-found error and turned
// it into a hard failure. The fix makes @foreach tolerant of a not-found source so
// the pre-existing (intended) empty-iteration behavior is restored without
// re-masking genuine evaluation errors.

// TestForeachUnresolvedSourceIsEmpty covers the source path resolving to a
// not-found error at the top level of an array element.
func TestForeachUnresolvedSourceIsEmpty(t *testing.T) {
	// Top-level literal array whose single element maps `facts` via @foreach over
	// an unresolvable source path ($.result.attribute — nothing is in scope).
	mappingValue := `{"mapping": [
   {
      "facts": {
         "@foreach($.result.attribute, attribute)":{
            "x":"=$loop.y"
         }
      }
   }
]}`

	var arrayMapping interface{}
	err := json.Unmarshal([]byte(mappingValue), &arrayMapping)
	assert.Nil(t, err)

	mappings := map[string]interface{}{"target": arrayMapping}
	factory := NewFactory(resolver)
	mapper, err := factory.NewMapper(mappings)
	assert.Nil(t, err)

	// Empty scope: $.result.attribute is unresolvable.
	scope := data.NewSimpleScope(map[string]interface{}{}, nil)

	results, err := mapper.Apply(scope)
	// Must not error — an absent source array is an empty iteration.
	assert.Nil(t, err)

	arr, ok := results["target"].([]interface{})
	assert.True(t, ok)
	assert.Equal(t, 1, len(arr))
	facts := arr[0].(map[string]interface{})["facts"]
	// facts came from @foreach over an absent source, so it is empty.
	if facts != nil {
		assert.Equal(t, 0, len(facts.([]interface{})))
	}
}

// TestForeachPresentSourceStillIterates guards against over-tolerating: a source
// that DOES resolve must still iterate normally.
func TestForeachPresentSourceStillIterates(t *testing.T) {
	mappingValue := `{"mapping": [
   {
      "facts": {
         "@foreach($.result.attribute, attribute)":{
            "x":"=$loop.y"
         }
      }
   }
]}`

	var arrayMapping interface{}
	err := json.Unmarshal([]byte(mappingValue), &arrayMapping)
	assert.Nil(t, err)

	mappings := map[string]interface{}{"target": arrayMapping}
	factory := NewFactory(resolver)
	mapper, err := factory.NewMapper(mappings)
	assert.Nil(t, err)

	sourceData := `{"attribute":[{"y":"one"},{"y":"two"}]}`
	var result interface{}
	err = json.Unmarshal([]byte(sourceData), &result)
	assert.Nil(t, err)
	scope := data.NewSimpleScope(map[string]interface{}{"result": result}, nil)

	results, err := mapper.Apply(scope)
	assert.Nil(t, err)

	arr := results["target"].([]interface{})
	facts := arr[0].(map[string]interface{})["facts"].([]interface{})
	assert.Equal(t, 2, len(facts))
	assert.Equal(t, "one", facts[0].(map[string]interface{})["x"])
	assert.Equal(t, "two", facts[1].(map[string]interface{})["x"])
}

// TestForeachSourceEvalErrorStillFails guards the other direction: a genuine
// (non-not-found) evaluation error in the source expression must still fail, so
// the fix does not reintroduce the error-swallowing behavior.
func TestForeachSourceEvalErrorStillFails(t *testing.T) {
	assert.False(t, isSourceNotFoundError("some real evaluation failure"))
	assert.True(t, isSourceNotFoundError("unable to evaluate path: attribute"))
}
