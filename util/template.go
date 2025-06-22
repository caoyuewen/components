package util

func Keys[K comparable, V any](m map[K]V) []K {
	var res []K
	for k := range m {
		res = append(res, k)
	}
	return res
}

func IncludeWithFunc[T any, S any](src []T, searchValue S, compare func(element T, searchValue S) bool) bool {
	for _, element := range src {
		if compare(element, searchValue) {
			return true
		}
	}
	return false
}
