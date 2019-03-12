package support

var aliases = make(map[string]map[string]string)

func RegisterAlias(contribType string, alias, ref string) {

	aliasToRefMap, exists := aliases[contribType]
	if !exists {
		aliasToRefMap = make(map[string]string)
		aliases[contribType] = aliasToRefMap
	}

	aliasToRefMap[alias] = ref
}

func GetAliasRef(contribType string, alias string) (string, bool) {

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
