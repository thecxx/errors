package errors

import (
	"errors"

	"github.com/thecxx/runpoint"
)

type wrappedError struct {
	err  error
	pc   *runpoint.PCounter
	refs map[string]any
}

// newWrappedError creates a new wrapped error.
func newWrappedError(err error, pc *runpoint.PCounter, refs map[string]any) *wrappedError {
	return &wrappedError{err: err, refs: refs, pc: pc}
}

// Error returns the error message.
func (w *wrappedError) Error() string {
	return w.err.Error()
}

// Unwrap returns the wrapped error.
func (w *wrappedError) Unwrap() error {
	return errors.Unwrap(w.err)
}

// PC returns the program counter.
func (w *wrappedError) PC() *runpoint.PCounter {
	return w.pc
}

// Value returns the value of the field.
func (w *wrappedError) Value(key string) any {
	value, ok := w.refs[key]
	if ok {
		return value
	}
	x, ok := w.err.(interface {
		// Value returns the value of the field from the error.
		Value(string) any
	})
	if !ok {
		return nil
	}
	return x.Value(key)
}

type joinError struct {
	*wrappedError
}

// newWrappedError creates a new wrapped error.
func newJoinError(err error, pc *runpoint.PCounter, refs map[string]any) *joinError {
	return &joinError{wrappedError: newWrappedError(err, pc, refs)}
}

// Unwrap returns the wrapped errors.
// If the error is not a wrapped error, it returns nil.
func (j *joinError) Unwrap() []error {
	u, ok := j.err.(interface {
		// Unwrap returns the wrapped errors.
		Unwrap() []error
	})
	if !ok {
		return nil
	}
	return u.Unwrap()
}
