package main

type Result[T any] struct {
	value T
	err   error
}

// Ok creates a Result representing a successful operation with the provided value of type T.
func Ok[T any](v T) Result[T] {
	return Result[T]{value: v, err: nil}
}

// Err creates a Result containing a zero value of type T and the provided error. It represents a failure state.
func Err[T any](e error) Result[T] {
	var zero T
	return Result[T]{value: zero, err: e}
}

// Inspectors ---------------------------------------------------------------

func (r Result[T]) IsOk() bool {
	return r.err == nil
}

func (r Result[T]) IsErr() bool {
	return r.err != nil
}

func (r Result[T]) Error() error {
	return r.err
}

// Unwrap returns the value if the Result is Ok, otherwise it panics with the error message.
func (r Result[T]) Unwrap() T {
	if r.err != nil {
		panic("called Unwrap on Err: " + r.err.Error())
	}
	return r.value
}

// Expect retrieves the value if the Result is Ok; otherwise, it panics with the provided message and the error.
func (r Result[T]) Expect(msg string) T {
	if r.err != nil {
		panic(msg + ": " + r.err.Error())
	}
	return r.value
}

// UnwrapOr returns the contained value if the Result is Ok, otherwise it returns the provided default value.
func (r Result[T]) UnwrapOr(def T) T {
	if r.err != nil {
		return def
	}
	return r.value
}

// UnwrapOrErr returns the value if the Result is Ok; otherwise, it returns a zero value and the provided new error.
func (r Result[T]) UnwrapOrErr(newErr error) (T, error) {
	if r.err != nil {
		var zero T
		return zero, newErr
	}
	return r.value, nil
}
