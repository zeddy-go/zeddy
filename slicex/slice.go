package slicex

func First[T any](arr ...T) (result T) {
	if len(arr) == 0 {
		return
	}

	return arr[0]
}
