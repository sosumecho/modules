package utils

func InArray[T comparable](arr []T, item T) bool {
	for _, v := range arr {
		if v == item {
			return true
		}
	}
	return false
}

func ToArray[T comparable, U any](item map[T]U) []U {
	var rs = make([]U, 0, len(item))
	for _, i := range item {
		rs = append(rs, i)
	}
	return rs
}
