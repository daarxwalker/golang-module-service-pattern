package vectorHelper

import (
	"fmt"
	"strings"

	"example/core/helper/stringHelper"
)

func Format(value string) string {
	latinized, err := stringHelper.Latinize(value)
	if err != nil {
		fmt.Println(err)
		return strings.ToLower(value)
	}
	return strings.ToLower(latinized)
}
