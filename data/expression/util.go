package expression

import (
	"encoding/json"
	"strconv"
	"strings"
)

// GetLiteral return the literal value if it is a literal
func GetLiteral(strVal string) (interface{}, bool) {

	newStrVal := strings.TrimSpace(strVal)

	var val interface{}
	var err error

	//check if is integer
	val, err = strconv.Atoi(newStrVal)
	if err == nil {
		return val, true
	}

	//check if is float
	val, err = strconv.ParseFloat(newStrVal, 64)
	if err == nil {
		return val, true
	}

	//check if is string
	s, isStr := isQuotedString(newStrVal)
	if isStr {
		return s, true
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
}

func isQuotedString(newStrVal string) (string, bool) {
	//String must surround with quotes and only one pair
	if newStrVal[0] == '"' {
		newStrVal = newStrVal[1:]
		if len(newStrVal) == strings.Index(newStrVal, `"`)+1 {
			return newStrVal[:len(newStrVal)-1], true
		}
		return "", false
	} else if newStrVal[0] == '\'' {
		newStrVal = newStrVal[1:]
		if len(newStrVal) == strings.Index(newStrVal, `'`)+1 {
			return newStrVal[:len(newStrVal)-1], true
		}
		return "", false

	} else if newStrVal[0] == '`' {
		newStrVal = newStrVal[1:]
		if len(newStrVal) == strings.Index(newStrVal, "`")+1 {
			return newStrVal[:len(newStrVal)-1], true
		}
		return "", false
	}
	return "", false
}
