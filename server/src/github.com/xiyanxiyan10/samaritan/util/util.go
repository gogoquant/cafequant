package util

import (
	"unicode"
)

// toUnderScoreCase myName => my_name
func ToUnderScoreCase(s string) string {
	runes := []rune(s)
	length := len(runes)
	out := []rune{}
	for i := 0; i < length; i++ {
		if i > 0 && unicode.IsUpper(runes[i]) && ((i+1 < length && unicode.IsLower(runes[i+1])) || unicode.IsLower(runes[i-1])) {
			out = append(out, '_')
		}
		out = append(out, unicode.ToLower(runes[i]))
	}
	return string(out)
}

// DeepCopy
func DeepCopy(value interface{}) interface{} {
	if valueMap, ok := value.(map[string]int); ok {
		newMap := make(map[string]int)
		for k, v := range valueMap {
			newMap[k] = v
		}
		return newMap
	}
	return value
}
