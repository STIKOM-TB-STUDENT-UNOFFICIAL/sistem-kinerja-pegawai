package utils

func Condition[T any](expect bool, a *T, b T) T {
	if expect {
		return *a
	}
	return b
}
