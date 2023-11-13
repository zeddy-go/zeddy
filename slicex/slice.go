package slicex

import "errors"

func First[T any](arr ...T) (result T) {
	if len(arr) == 0 {
		return
	}

	return arr[0]
}

func Contains[T comparable](finds any, data []T) (found bool) {
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
