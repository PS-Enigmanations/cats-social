package util

type Result[T interface{}] struct {
	Result T
	Error  error
}
