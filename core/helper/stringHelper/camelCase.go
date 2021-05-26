package stringHelper

import (
	"regexp"
	"strings"
)

var link = regexp.MustCompile("(^[A-Za-z])|_([A-Za-z])")

func CamelCase(str string) string {
	camelStr := link.ReplaceAllStringFunc(str, func(s string) string {
		return strings.ToUpper(strings.Replace(s, "_", "", -1))
	})
	return strings.ToLower(string(camelStr[0])) + camelStr[1:]
}
