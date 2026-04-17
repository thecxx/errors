// Package errors is a thin facade over the standard library's errors and fmt packages.
// It centralizes error construction and wrapping while preserving unwrap semantics for
// errors.Is and errors.As.
package errors

import (
	"errors"
	"fmt"

	"github.com/thecxx/runpoint"
)

// New returns an error that formats as the given text. It delegates to the standard
// library's errors.New.
func New(message string) error {
	return newWrappedError(errors.New(message), runpoint.PC(1), nil)
}

// Newf returns fmt.Errorf(format, args...). Use Wrap when adding a fixed prefix around
// an existing error with consistent %w chaining.
func Newf(format string, args ...any) error {
	return newWrappedError(fmt.Errorf(format, args...), runpoint.PC(1), nil)
}

// Join joins the errors with the standard library's errors.Join.
func Join(errs ...error) error {
	joined := errors.Join(errs...)
	if joined == nil {
		return nil
	}
	return newJoinError(joined, runpoint.PC(1), nil)
}

// Wrap returns an error that prefixes err with message. The result supports errors.Unwrap
// and unwraps to err, so Is and As traverse the chain like the standard library.
func Wrap(err error, message string) error {
	return newWrappedError(fmt.Errorf("%s: %w", message, err), runpoint.PC(1), nil)
}

// WithValue returns a new error with the given key and value, and the program counter of the caller.
func WithValue(err error, key string, value any) error {
	return newWrappedError(err, runpoint.PC(1), map[string]any{key: value})
}

// WithValues returns a new error with the given values, and the program counter of the caller.
func WithValues(err error, values map[string]any) error {
	return newWrappedError(err, runpoint.PC(1), values)
}

// Unwrap returns the result of errors.Unwrap(err) from the standard library.
func Unwrap(err error) error {
	return errors.Unwrap(err)
}

// PC returns the program counter of the error.
// If the error is not a wrapped error, it returns nil.
func PC(err error) *runpoint.PCounter {
	x, ok := err.(interface {
		// PC returns the program counter of the error.
		PC() *runpoint.PCounter
	})
	if !ok {
		return nil
	}
	return x.PC()
}

// Value returns the value of the field from the error.
// If the error is not a wrapped error, it returns nil.
func Value(err error, key string) any {
	x, ok := err.(interface {
		// Value returns the value of the field from the error.
		Value(string) any
	})
	if !ok {
		return nil
	}
	return x.Value(key)
}

// Is reports whether any error in err's chain matches target, delegating to the
// standard library's errors.Is.
func Is(err error, target error) bool {
	return errors.Is(err, target)
}

// As finds the first error in err's chain that matches target, delegating to the
// standard library's errors.As.
func As(err error, target any) bool {
	return errors.As(err, target)
}
