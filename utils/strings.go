package utils

func ArrayContains[T comparable](array []T, search T) bool {
	for _, item := range array {
		if search == item {
			return true
		}
	}
	return false
}
