package numberHelper

import (
	"fmt"
	"strconv"
)

func Int(value string) int {
	number, err := strconv.Atoi(value)
	if err != nil {
		fmt.Println(err)
		return 0
	}
	return number
}
