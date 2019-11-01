package engine

import (
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

//TestNewPooledConfigOk
func TestNewPooledConfigDefault(t *testing.T) {
	pooledConfig := NewPooledRunnerConfig()

	// assert Success
	assert.Equal(t, DefaultRunnerWorkers, pooledConfig.NumWorkers)
	assert.Equal(t, DefaultRunnerQueueSize, pooledConfig.WorkQueueSize)
}

//TestNewPooledConfigOk
func TestNewPooledConfigOverride(t *testing.T) {
	previousWorkers := os.Getenv(EnvKeyRunnerWorkers)
	defer os.Setenv(EnvKeyRunnerWorkers, previousWorkers)
	previousQueue := os.Getenv(EnvKeyRunnerQueueSize)
	defer os.Setenv(EnvKeyRunnerQueueSize, previousQueue)

	newWorkersValue := 6
	newQueueValue := 60

	// Change values
	_ = os.Setenv(EnvKeyRunnerWorkers, strconv.Itoa(newWorkersValue))
	_ = os.Setenv(EnvKeyRunnerQueueSize, strconv.Itoa(newQueueValue))

	pooledConfig := NewPooledRunnerConfig()

	// assert Success
	assert.Equal(t, newWorkersValue, pooledConfig.NumWorkers)
	assert.Equal(t, newQueueValue, pooledConfig.WorkQueueSize)
}
