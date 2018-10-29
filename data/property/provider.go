package property

import (
	"fmt"
	"github.com/project-flogo/core/support/log"
)

var (
	providers = make(map[string]Provider)
)

type Provider interface {
	GetProperty(name string) interface{}
}

func GetProvider(id string) Provider {
	return providers[id]
}

func RegisterProvider(id string, provider Provider) error {

	if id == "" {
		return fmt.Errorf("'ref' must be specified when registering provider")
	}

	if provider == nil {
		return fmt.Errorf("cannot register 'nil' provider")
	}

	if _, dup := providers[id]; dup {
		return fmt.Errorf("provider already registered: %s", id)
	}

	log.RootLogger().Debugf("Registering provider [ %s ]", id)

	providers[id] = provider

	return nil
}
