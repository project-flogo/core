package resolve

import (
	"fmt"
	"strings"

	"github.com/project-flogo/core/data"
)

// Resolver is for resolving a value for a specific environment or construct, ex. OS environment
type Resolver interface {
	// GetResolverInfo returns ResolverInfo which contains information about the resolver
	GetResolverInfo() *ResolverInfo

	// Resolve resolves to a value using the specified item and valueName, note: the scope might not be used
	Resolve(scope data.Scope, itemName, valueName string) (interface{}, error)
}

// CompositeResolver is a resolver that is typically composed of other simple 'Resolvers',
// the specified directive is dispatched to the appropriate embedded resolvers
type CompositeResolver interface {
	// GetResolution creates a "delayed" resolution object, who's value isn't fully resolved
	// util GetValue is called
	GetResolution(resolveDirective string) (Resolution, error)

	// Resolve resolves to a value using the specified directive
	Resolve(resolveDirective string, scope data.Scope) (value interface{}, err error)
}

// Resolution structure that is allows for delayed resolving of values, the value can then be fully
// resolved calling GetValue for a particular scope
type Resolution interface {
	// IsStatic indicates that resolution can be done statically without a scope
	IsStatic() bool

	// GetValue resolves and returns the value using the specified scope
	GetValue(scope data.Scope) (interface{}, error)
}

// NewResolverInfo creates a ResolverInfo object
func NewResolverInfo(isStatic, usesItemFormat bool) *ResolverInfo {
	return &ResolverInfo{isStatic: isStatic, usesItemFormat: usesItemFormat}
}

func NewImplicitResolverInfo(isStatic, isImplicit bool) *ResolverInfo {
	return &ResolverInfo{isStatic: isStatic, isImplicit: isImplicit}
}

func CreateResolverInfo(options ...Option) *ResolverInfo {
	rf := &ResolverInfo{}
	for _, option := range options {
		option(rf)
	}
	return rf
}

type Option func(*ResolverInfo)

func OptUseItemFormat(r *ResolverInfo) {
	r.usesItemFormat = true
}
func OptImplicit(r *ResolverInfo) {
	r.isImplicit = true
}
func OptStatic(r *ResolverInfo) {
	r.isStatic = true
}

// ResolverInfo structure that contains information about the resolver
type ResolverInfo struct {
	usesItemFormat bool
	isStatic       bool
	isImplicit     bool
}

// IsStatic determines if the resolver's values are static and can be resolved immediately without a scope
func (i *ResolverInfo) IsStatic() bool {
	return i.isStatic
}

// UsesItemFormat determines if the resolver uses the item format (ex. $test[itemName])
func (i *ResolverInfo) UsesItemFormat() bool {
	return i.usesItemFormat
}

// IsImplicit determines if the resolver try to uses the item format and no item format then
func (i *ResolverInfo) IsImplicit() bool {
	return i.isImplicit
}

// GetResolverInfo gets the resolver name and position to start parsing the ResolutionDetails from
func GetResolverInfo(toResolve string) (string, int) {

	if toResolve[0] == '.' {
		return ".", 1
	}

	for i, char := range toResolve {
		if char == '.' || char == '[' {
			return toResolve[0:i], i
		}
	}

	return toResolve, len(toResolve)
}

// ResolveDirectiveDetails is the Resolve Directive broken into components to assist in resolving the value
type ResolveDirectiveDetails struct {
	ItemName  string
	ValueName string
	Path      string
}

// GetResolveDirectiveDetails breaks Resolution Directive into components
func GetResolveDirectiveDetails(directive string, hasItems, isImplicit bool) (*ResolveDirectiveDetails, error) {

	//todo optimize
	details := &ResolveDirectiveDetails{}

	if len(directive) == 0 {

		if hasItems {
			return nil, fmt.Errorf("invalid resolve directive: '%s' needs to start with [item]", directive)
		}

		return details, nil
	}

	start := 0
	strLen := len(directive)
	hasNamedValue := true

	//isImplicit will try to support both item or without item
	if isImplicit || hasItems {
		//uses the "item format" (ex. foo[bar].valueName; where 'bar' is the item)
		if directive[0] != '[' {
			if hasItems {
				return nil, fmt.Errorf("invalid resolve directive: '%s' needs to start with [item]", directive)
			}

			if isImplicit {
				hasItems = false
				hasNamedValue = true
			}
		} else {
			hasItems = true
			hasNamedValue = false
		}

		if hasItems {
			start = 1

			for i := 1; i < strLen; i++ {
				if directive[i] == ']' {
					details.ItemName = removeQuotes(directive[start:i])
					start = i + 1

					//if we started with an item, it must either end or the next segment should start with '.' or '['
					if start < strLen {
						if directive[start] != '.' && directive[start] != '[' {
							return nil, fmt.Errorf("invalid resolve directive: '%s'", directive)
						}

						if directive[start] == '.' {
							hasNamedValue = true
							start++
						}
					}

					break
				}
			}
		}

	}
	var i int

	if hasNamedValue {
		if start == 0 && directive[0] == '.' {
			start = 1
		}

		for i = start; i < strLen; i++ {
			if directive[i] == '.' || directive[i] == '[' {
				details.ValueName = directive[start:i]
				start = i
				break
			}
		}
	}

	if i == strLen {
		// we have gotten to the end of the string, so the last part of the string is the ValueName
		details.ValueName = directive[start:]
	} else if start < strLen {
		// we have a remaining component, should be the 'path'
		details.Path = directive[start:]
	}

	return details, nil
}

var ends = map[byte]byte{'(': ')', '"': '"', '\'': '\'', '[': ']', '`': '`'}

func isLetter(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
}

// IsResolveExpr determines if the provide expression string is a resolver expression, ie. "$env[MY_VAR]
func IsResolveExpr(exprStr string) bool {

	strLen := len(exprStr)

	if exprStr[0] == '$' && strLen > 2 {

		if exprStr[1] != '.' && !isLetter(exprStr[2]) {
			return false
		}

		for i := 1; i < strLen; i++ {
			switch c := exprStr[i]; c {
			case ' ':
				return false
			case '"', '\'', '`':
				end := ends[c]
				i++
				for i < len(exprStr) {
					if exprStr[i] == end {
						break
					}
					i++
				}
			case '[':
				end := ends[c]
				i++
				j := i
				for i < len(exprStr) {
					if exprStr[i] == end {
						//Checking whether value insde [] is expr
						if hasExprChar(exprStr[j:i]) {
							return false
						}
						break
					}
					i++
				}
			case '.':
				if i+1 >= strLen || !isLetter(exprStr[i+1]) {
					return false
				}

			default:
				if isExprChar(c) {
					return false
				}
			}

		}
	} else {
		return false
	}

	return true
}

func hasExprChar(str string) bool {
	//condition expression, tenray expression, array indexer expression
	if len(str) > 0 {
		//String
		if str[0] == '"' || str[0] == '\'' || str[0] == '`' {
			return false
		} else {
			strLen := len(str)
			for i := 0; i < strLen; i++ {
				if isExprChar(str[i]) {
					return true
				}
			}
		}
	}
	return false
}

func isExprChar(ch byte) bool {
	//condition expression, tenray expression, array indexer expression
	switch ch {
	case '(', '=', '>', '<', '*', '/', '!', '&', '%', '+', '-', '|', '?', ':', '$':
		return true
	}
	return false
}

func removeQuotes(str string) string {
	if strings.HasPrefix(str, `"`) && strings.HasSuffix(str, `"`) {
		return str[1 : len(str)-1]
	}
	if strings.HasPrefix(str, "'") && strings.HasSuffix(str, "'") {
		return str[1 : len(str)-1]
	}
	if strings.HasPrefix(str, "`") && strings.HasSuffix(str, "`") {
		return str[1 : len(str)-1]
	}
	return str
}
