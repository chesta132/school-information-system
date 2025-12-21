package slicelib

func Map[T any, R any](slice []T, f func(idx int, val T) R) []R {
	result := make([]R, len(slice))
	for i, v := range slice {
		result[i] = f(i, v)
	}
	return result
}

func Filter[T any](slice []T, f func(idx int, val T) bool) []T {
	if len(slice) == 0 {
		return slice
	}

	result := make([]T, 0, len(slice))
	for i, v := range slice {
		if f(i, v) {
			result = append(result, v)
		}
	}
	return result
}

func Unique[T comparable](slice []T) []T {
	if len(slice) <= 1 {
		return slice
	}

	seen := make(map[T]struct{}, len(slice))
	result := make([]T, 0, len(slice))

	for _, element := range slice {
		if _, ok := seen[element]; !ok {
			seen[element] = struct{}{}
			result = append(result, element)
		}
	}
	return result
}
