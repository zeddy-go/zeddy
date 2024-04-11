package tool

import "reflect"

func Or[T any](value T, orValue T) T {
	if reflect.ValueOf(value).IsZero() {
		return orValue
	}

	return value
}
