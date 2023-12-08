package mapx

func CloneSimpleMapSlice[K comparable, V any](src map[K][]V) (dest map[K][]V) {
	if src == nil {
		return nil
	}

	// Find total number of values.
	nv := 0
	for _, vv := range src {
		nv += len(vv)
	}
	sv := make([]V, nv) // shared backing array for headers' values
	dest = make(map[K][]V, len(src))
	for k, vv := range src {
		if vv == nil {
			// Preserve nil values. ReverseProxy distinguishes
			// between nil and zero-length header values.
			dest[k] = nil
			continue
		}
		n := copy(sv, vv)
		dest[k] = sv[:n:n]
		sv = sv[n:]
	}
	return dest
}
