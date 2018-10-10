package expression

import (
	"encoding/json"
	"strconv"
	"strings"
)

// GetLiteral return the literal value if it is a literal
func GetLiteral(strVal string) (interface{}, bool) {

	if strVal[0] == '=' {

		newStrVal := strings.TrimSpace(strVal[1:])

		var val interface{}
		var err error

		//check if is integer
		_, err = strconv.Atoi(newStrVal)
		if err == nil {
			return val, true
		}

		//check if is float
		val, err = strconv.ParseFloat(newStrVal, 64)
		if err == nil {
			return val, true
		}

		//check if is string
		if newStrVal[0] == '`' && newStrVal[len(newStrVal)-1] == '`' {
			return newStrVal[1 : len(newStrVal)-1], true
		}
		if newStrVal[0] == '\'' && newStrVal[len(newStrVal)-1] == '\'' {
			return newStrVal[1 : len(newStrVal)-1], true
		}
		if newStrVal[0] == '"' && newStrVal[len(newStrVal)-1] == '"' {
			return newStrVal[1 : len(newStrVal)-1], true
		}

		//check for booleans
		if strings.EqualFold(newStrVal, "false") {
			return false, true
		}

		if strings.EqualFold(newStrVal, "true") {
			return true, true
		}

		//check if is object or array
		var js interface{}
		err = json.Unmarshal([]byte(newStrVal), &js)
		if err == nil {
			return js, true
		}

		return nil, false
	} else {
		return strVal, true
	}
}
