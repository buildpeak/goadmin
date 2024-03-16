package slices

func Map[T, U any](src []T, fnc func(T) U) []U {
	if len(src) == 0 {
		return []U{}
	}

	r := make([]U, len(src))
	for i, v := range src {
		r[i] = fnc(v)
	}

	return r
}
