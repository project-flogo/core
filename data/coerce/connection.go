package coerce

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/project-flogo/core/support/connection"
)

func ToConnection(val interface{}) (connection.Manager, error) {

	switch t := val.(type) {
	case string:
		if strings.HasPrefix(t, "conn://") {
			id := t[7:]
			cm := connection.GetManager(id)

			if cm == nil {
				return nil, fmt.Errorf("connection with id '%s' not configured", t)
			}

			return cm, nil
		} else {
			cc := &connection.Config{}
			if t != "" {
				err := json.Unmarshal([]byte(t), cc)
				if err != nil {
					return nil, fmt.Errorf("'%s' is not a valid connection config", t)
				}
			}

			f := connection.GetManagerFactory(cc.Ref)
			if f == nil {
				return nil, fmt.Errorf("no connection factory registered for '%s'", cc.Ref)
			}

			cm, err := f.NewManager(cc.Settings)
			if err != nil {
				return nil, err
			}

			return cm, nil
		}
	case connection.Manager:
		return t, nil
	default:
		// try to create config from map[string]interface{}
		return nil, fmt.Errorf("unable to create connection from '%#v'", val)
	}
}