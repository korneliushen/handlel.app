package lib

// Legg til funksjoner som skal / vil kunne brukes ulike steder

// Sjekker om et element er i et array, returnerer true eller false
func IsIn[T comparable](e T, arr []T) bool {
	for i := range arr {
		if e == arr[i] {
			return true
		}
	}
	return false
}
