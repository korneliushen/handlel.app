package lib

func isIn[T comparable](e T, arr []T) bool {
	for i := range arr {
		if e == arr[i] {
			return true
		}
	}
	return false
}
