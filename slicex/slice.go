package slicex

import "errors"

func Find[T any](f func(item T) bool, arr ...T) (result T, ok bool) {
	if len(arr) == 0 {
		return
	}

	for _, item := range arr {
		if f(item) {
			ok = true
			result = item
			return
		}
	}

	return
}

func Filter[T any](f func(item T) bool, arr ...T) (list []T) {
	if len(arr) == 0 {
		return
	}

	list = make([]T, 0, len(arr))
	for _, item := range arr {
		if !f(item) {
			list = append(list, item)
		}
	}

	return
}

func Last[T any](arr ...T) (result T) {
	if len(arr) == 0 {
		return
	}

	return arr[len(arr)-1]
}

func Pop[T any](arr ...T) (newArr []T, result T) {
	if len(arr) == 0 {
		return
	}
	return arr[:len(arr)-1], arr[len(arr)-1]
}

func First[T any](arr ...T) (result T) {
	if len(arr) == 0 {
		return
	}

	return arr[0]
}

func Contains[T comparable](finds any, data []T) (found bool) {
	if data == nil {
		return false
	}
	switch x := finds.(type) {
	case []T:
		for _, find := range x {
			for _, item := range data {
				if item == find {
					found = true
					return
				}
			}
		}
	case T:
		for _, item := range data {
			if item == x {
				found = true
				return
			}
		}
	default:
		panic(errors.New("type is mismatched"))
	}

	return
}
