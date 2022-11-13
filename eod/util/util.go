package util

func Map[T, V any](v []T, m func(a T) V) []V {
	out := make([]V, len(v))
	for i, val := range v {
		out[i] = m(val)
	}
	return out
}
