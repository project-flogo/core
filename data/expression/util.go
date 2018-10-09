package expression

import (
	"encoding/json"
	"strconv"
	"strings"
)

// GetLiteral return the literal value if it is a literal
func GetLiteral(strVal string) (interface{}, bool) {

	var val interface{}
	var err error

	//check if is integer
	val, err = strconv.Atoi(strVal)
	if err == nil {
		return val, true
	}

	//check if is float
	val, err = strconv.ParseFloat(strVal, 64)
	if err == nil {
		return val, true
	}

	//check if is string
	if strVal[0] == '`' && strVal[len(strVal)-1] == '`' {
		return strVal[1 : len(strVal)-1], true
	}
	if strVal[0] == '\'' && strVal[len(strVal)-1] == '\'' {
		return strVal[1 : len(strVal)-1], true
	}
	if strVal[0] == '"' && strVal[len(strVal)-1] == '"' {
		return strVal[1 : len(strVal)-1], true
	}

	//check for booleans
	if strings.EqualFold(strVal, "false") {
		return false, true
	}

	if strings.EqualFold(strVal, "true") {
		return true, true
	}

	//check if is object or array
	var js interface{}
	err = json.Unmarshal([]byte(strVal), &js)
	if err == nil {
		return js, true
	}

	return nil, false
}
