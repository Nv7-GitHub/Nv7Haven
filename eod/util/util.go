package util

func Map[T, V any](v []T, m func(a T) V) []V {
	out := make([]V, len(v))
	for i, val := range v {
		out[i] = m(val)
	}
	return out
}

type ordered interface {
	~int | ~float64 | ~float32 | ~int32 | ~int64
}

func Min[T ordered](a, b T) T {
	if a < b {
		return a
	}
	return b
}
