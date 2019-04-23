package channels

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNew_Started(t *testing.T) {
	channels = map[string]*channelImpl{}
	active = true
	channel, err := New("test1", 5)
	assert.Nil(t, channel)
	assert.NotNil(t, err)
}

func TestNew(t *testing.T) {
	channels = map[string]*channelImpl{}
	active = false
	channel, err := New("test1", 5)
	assert.Nil(t, err)
	assert.NotNil(t, channel)
}

func TestNilChannel(t *testing.T) {
	channels = map[string]*channelImpl{}
	active = false

	c := Get("dne")
	assert.Nil(t, c)
	assert.True(t, c == nil)
}

func TestNewDup(t *testing.T) {
	channels = map[string]*channelImpl{}
	active = false
	channel, err := New("test1", 5)
	assert.Nil(t, err)
	assert.NotNil(t, channel)
	channel, err = New("test1", 5)
	assert.Nil(t, channel)
	assert.NotNil(t, err)
}

func TestChannel_Name(t *testing.T) {
	channels = map[string]*channelImpl{}
	active = false
	channel, err := New("test1", 5)
	assert.Nil(t, err)
	cImpl := channel.(*channelImpl)
	assert.Equal(t, "test1", cImpl.name)
}

func TestStart(t *testing.T) {
	channels = map[string]*channelImpl{}
	channel, err := New("test1", 5)
	assert.Nil(t, err)
	assert.NotNil(t, channel)
	_ = Start()
	defer Stop()
	assert.True(t, active)

	channel2, err2 := New("test2", 5)
	assert.NotNil(t, err2)
	assert.Nil(t, channel2)

	cImpl := channel.(*channelImpl)
	assert.True(t, cImpl.active)
}

func TestChannel_Publish(t *testing.T) {
	channels = map[string]*channelImpl{}
	channel, err := New("test1", 1)
	assert.Nil(t, err)
	assert.NotNil(t, channel)

	channel.Publish(1)

	cImpl := channel.(*channelImpl)

	select {
	case msg := <-cImpl.ch:
		assert.Equal(t, 1, msg)
	default:
		assert.Fail(t, "no message received")
	}
}

func TestChannel_PublishNoWait(t *testing.T) {
	channels = map[string]*channelImpl{}
	channel, err := New("test1", 1)
	assert.Nil(t, err)
	assert.NotNil(t, channel)

	channel.Publish(1)
	sent := channel.PublishNoWait(2)

	assert.False(t, sent)
}

func TestChannel_AddListenerStarted(t *testing.T) {
	channels = map[string]*channelImpl{}
	channel, err := New("test1", 1)
	assert.Nil(t, err)
	assert.NotNil(t, channel)

	cImpl := channel.(*channelImpl)
	_ = cImpl.Start()

	err = cImpl.RegisterCallback(func(msg interface{}) {
		//dummy
	})
	assert.NotNil(t, err)
}

func TestChannel_AddListener(t *testing.T) {
	channels = map[string]*channelImpl{}
	channel, err := New("test1", 1)
	assert.Nil(t, err)
	assert.NotNil(t, channel)

	cImpl := channel.(*channelImpl)

	err = cImpl.RegisterCallback(func(msg interface{}) {
		//dummy
	})
	assert.Equal(t, 1, len(cImpl.callbacks))
}

type cbTester struct {
	called int
	val    interface{}
}

func (cbt *cbTester) onMessage(msg interface{}) {
	cbt.called++
	cbt.val = msg
}

func TestChannel_Callback(t *testing.T) {
	channels = map[string]*channelImpl{}
	channel, err := New("test1", 1)
	assert.Nil(t, err)
	assert.NotNil(t, channel)

	cImpl := channel.(*channelImpl)

	cbt := &cbTester{}
	err = cImpl.RegisterCallback(cbt.onMessage)
	_ = Start()

	channel.Publish(22)
	time.Sleep(100 * time.Millisecond)
	assert.Equal(t, 1, cbt.called)
	assert.Equal(t, 22, cbt.val)
}
