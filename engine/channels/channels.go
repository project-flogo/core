package channels

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/project-flogo/core/support/log"
)

var channels = make(map[string]*channelImpl)
var active bool

type Channel interface {
	RegisterCallback(callback OnMessage) error
	Publish(msg interface{})
	PublishNoWait(msg interface{}) bool
}

type OnMessage func(msg interface{})

// Creates a new channel, channels have to be created before the engine starts
func New(name string, bufferSize int) (Channel, error) {

	if active {
		return nil, errors.New("cannot create channel after engine has been started")
	}

	if _, dup := channels[name]; dup {
		return nil, errors.New("channel already exists: " + name)
	}

	channel := &channelImpl{name: name, ch: make(chan interface{}, bufferSize)}
	channels[name] = channel

	return channel, nil

}

// Count returns the number of channels
func Count() int {
	return len(channels)
}

// Get gets the named channel
func Get(name string) Channel {
	if ch, ok := channels[name]; ok {
		return ch
	}
	return nil

	//note: can't use nil check with this code
	//return channels[name]
}

func Start() error {
	active = true

	var started []*channelImpl

	for _, channel := range channels {
		err := channel.Start()
		if err != nil {
			for _, startedChannel := range started {
				_ = startedChannel.Stop()
			}
			return fmt.Errorf("failed to start channel '%s', error: %s", channel.name, err.Error())
		}
		log.RootLogger().Debugf("Started Engine Channel: %s", channel.name)
		started = append(started, channel)
	}

	return nil
}

func Stop() error {
	for _, channel := range channels {
		err := channel.Stop()
		if err != nil {
			log.RootLogger().Warnf("error stopping channel '%s', error: %s", channel.name, err.Error())
		}
	}
	channels = make(map[string]*channelImpl)

	active = false

	return nil
}

type channelImpl struct {
	name      string
	callbacks []OnMessage
	ch        chan interface{}
	active    bool
}

func (c *channelImpl) Start() error {
	c.active = true
	go c.processEvents()

	return nil
}

func (c *channelImpl) Stop() error {
	close(c.ch)
	c.active = false

	return nil
}

func (c *channelImpl) RegisterCallback(callback OnMessage) error {

	if c.active {
		return errors.New("cannot add listener after channel has been started")
	}

	c.callbacks = append(c.callbacks, callback)
	return nil
}

func (c *channelImpl) Publish(msg interface{}) {
	c.ch <- msg
}

func (c *channelImpl) PublishNoWait(msg interface{}) bool {

	sent := false
	select {
	case c.ch <- msg:
		sent = true
	default:
		sent = false
	}

	return sent
}

func (c *channelImpl) processEvents() {

	for {
		select {
		case val, ok := <-c.ch:
			if !ok {
				//channel closed, so return
				return
			}

			for _, callback := range c.callbacks {
				go callback(val)
			}
		}
	}
}

// Decode decodes the channel descriptor
func Decode(channelDescriptor string) (string, int) {
	idx := strings.Index(channelDescriptor, ":")
	buffSize := 0
	chanName := channelDescriptor

	if idx > 0 {
		bSize, err := strconv.Atoi(channelDescriptor[idx+1:])
		if err != nil {
			log.RootLogger().Warnf("invalid channel buffer size '%s', defaulting to buffer size of %d", channelDescriptor[idx+1:], buffSize)
		} else {
			buffSize = bSize
		}

		chanName = channelDescriptor[:idx]
	}

	return chanName, buffSize
}
