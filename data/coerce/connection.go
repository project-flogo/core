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

			cm, err := connection.NewManager(cc)
			if err != nil {
				return nil, err
			}

			return cm, nil
		}
	case connection.Manager:
		return t, nil
	case map[string]interface{}:
		cfg, err := connection.ToConfig(t)
		if err != nil {
			return nil, err
		}

		cm, err := connection.NewManager(cfg)
		if err != nil {
			return nil, err
		}

		return cm, nil
	default:
		// try to create config from map[string]interface{}
		return nil, fmt.Errorf("unable to create connection from '%#v'", val)
	}
}

