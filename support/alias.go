package support

import "fmt"

var aliases = make(map[string]map[string]string)

func RegisterAlias(contribType, alias, ref string) error {

	aliasToRefMap, exists := aliases[contribType]
	if !exists {
		aliasToRefMap = make(map[string]string)
		aliases[contribType] = aliasToRefMap
	}

	if _, exists := aliasToRefMap[alias]; exists {
		return fmt.Errorf("alias '%s' for %s already registered", alias, contribType)
	}

	aliasToRefMap[alias] = ref
	return nil
}

func GetAliasRef(contribType, alias string) (string, bool) {

	if alias == "" {
		return "", false
	}

	if alias[0] == '#' {
		alias = alias[1:]
	}

	aliasToRefMap, exists := aliases[contribType]
	if !exists {
		return "", false
	}

	ref, exists := aliasToRefMap[alias]
	if !exists {
		return "", false
	}

	return ref, true
}
