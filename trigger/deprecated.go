package trigger

import (
	"encoding/json"
	"fmt"

	"github.com/project-flogo/core/data"
	"github.com/project-flogo/core/support/log"
)

//DEPRECATED
func LegacyRegister(ref string, f Factory) error {

	if ref == "" {
		return fmt.Errorf("'ref' must be specified when registering")
	}

	if f == nil {
		return fmt.Errorf("cannot register trigger with 'nil' trigger factory")
	}

	if triggerFactories[ref] != nil {
		return fmt.Errorf("trigger already registered for ref '%s'", ref)
	}

	log.RootLogger().Debugf("Registering legacy trigger [ %s ]", ref)

	triggerFactories[ref] = f

	return nil
}

////////////////////////////
// OldMetadata

// OldMetadata is the oldMetadata for a Trigger
type oldMetadata struct {
	Ref      string              `json:"ref"`
	Handler  *oldHandlerMetadata `json:"handler"`
	Settings []*data.Attribute   `json:"settings"`
	Output   []*data.Attribute   `json:"output"`
	Reply    []*data.Attribute   `json:"reply"`
}

// EndpointOldMetadata is the oldMetadata for a Trigger Endpoint
type oldHandlerMetadata struct {
	Settings []*data.Attribute `json:"settings"`
}

// NewOldMetadata creates a OldMetadata object from the json representation
//DEPRECATED
func metadataFromOldJson(jsonMetadata string) *Metadata {
	oldMd := &oldMetadata{}
	err := json.Unmarshal([]byte(jsonMetadata), oldMd)
	if err != nil {
		panic("Unable to parse trigger oldMetadata: " + err.Error())
	}

	md := &Metadata{}

	md.ID = oldMd.Ref

	if len(oldMd.Settings) > 0 {
		md.Settings = make(map[string]data.TypedValue, len(oldMd.Settings))
		for _, attr := range oldMd.Settings {
			md.Settings[attr.Name()] = data.NewTypedValue(attr.Type(), attr.Value())
		}
	}
	if len(oldMd.Output) > 0 {
		md.Output = make(map[string]data.TypedValue, len(oldMd.Output))
		for _, attr := range oldMd.Output {
			md.Output[attr.Name()] = data.NewTypedValue(attr.Type(), attr.Value())
		}
	}

	if len(oldMd.Reply) > 0 {
		md.Reply = make(map[string]data.TypedValue, len(oldMd.Reply))
		for _, attr := range oldMd.Reply {
			md.Reply[attr.Name()] = data.NewTypedValue(attr.Type(), attr.Value())
		}
	}
	if oldMd.Handler != nil && len(oldMd.Handler.Settings) > 0 {
		md.HandlerSettings = make(map[string]data.TypedValue, len(oldMd.Handler.Settings))
		for _, attr := range oldMd.Handler.Settings {
			md.HandlerSettings[attr.Name()] = data.NewTypedValue(attr.Type(), attr.Value())
		}
	}

	return md
}
