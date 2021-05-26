package sliceHelper

func Contains(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

func IntContains(slice []int, val int) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

func Remove(slice []string, val string) []string {
	var r []string
	for _, item := range slice {
		if item != val {
			r = append(r, item)
		}
	}
	return r
}
