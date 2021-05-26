package numberHelper

import "strconv"

func Float(value string) float64 {
	number, _ := strconv.ParseFloat(value, 64)
	return number
}
