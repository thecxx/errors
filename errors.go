// Package errors is a thin facade over the standard library's errors and fmt packages.
// It centralizes error construction and wrapping while preserving unwrap semantics for
// errors.Is and errors.As.
package errors

import (
	"errors"
	"fmt"
)

// New returns an error that formats as the given text. It delegates to the standard
// library's errors.New.
func New(message string) error {
	return errors.New(message)
}

// Newf returns fmt.Errorf(format, args...). Use Wrap when adding a fixed prefix around
// an existing error with consistent %w chaining.
func Newf(format string, args ...any) error {
	return fmt.Errorf(format, args...)
}

// Wrap returns an error that prefixes err with message. The result supports errors.Unwrap
// and unwraps to err, so Is and As traverse the chain like the standard library.
func Wrap(err error, message string) error {
	return fmt.Errorf("%s: %w", message, err)
}

// Unwrap returns the result of errors.Unwrap(err) from the standard library.
func Unwrap(err error) error {
	return errors.Unwrap(err)
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
