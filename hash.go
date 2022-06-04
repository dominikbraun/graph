package graph

type Hash[K comparable, T any] func(T) K

func StringHash(v string) string {
	return v
}

func IntHash(v int) int {
	return v
}
