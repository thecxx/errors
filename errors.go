package errors

import (
	"path"
	"reflect"
	"runtime"
)

func New(text string, references ...ErrorReference) error {
	file, line := fileline(1)
	return newWrappedError(nil, text, file, line, references...)
}

func Wrap(prev error, text string, references ...ErrorReference) error {
	file, line := fileline(1)
	return newWrappedError(prev, text, file, line, references...)
}

func Unwrap(err error) error {
	w, ok := isWrappedError(err)
	if !ok {
		return nil
	}
	return w.prev
}

func FileLine(err error) (file string, line int) {
	e, ok := isWrappedError(err)
	if !ok {
		return
	}
	return e.file, e.line
}

func References(err error) map[string]interface{} {
	e, ok := isWrappedError(err)
	if !ok {
		return make(map[string]interface{})
	}
	return e.refs
}

func Primary(err error) error {
	return primary(err)
}

func primary(err error) error {
	u, ok := isWrappedError(err)
	if ok && u.prev != nil {
		return primary(u.prev)
	}
	return err
}

func Stack(err error) []error {
	if err == nil {
		return make([]error, 0)
	}
	var (
		e = err
		s []error
	)
	for e != nil {
		s = append(s, e)
		// Prev error
		u, ok := isWrappedError(e)
		if ok {
			e = u.prev
		} else {
			e = nil
		}
	}
	return s
}

func isWrappedError(err error) (e *wrappedError, ok bool) {
	e, ok = err.(*wrappedError)
	return
}

func fileline(skip int) (file string, line int) {
	var ok bool

	_, file, line, ok = runtime.Caller(skip + 1)
	if !ok {
		return "", 0
	}
	return path.Base(file), line
}

type ErrorReference func(*wrappedError)

func Ref(key string, value interface{}) ErrorReference {
	return func(w *wrappedError) {
		w.refs[key] = value
	}
}

type wrappedError struct {
	prev error
	file string
	line int
	text string
	refs map[string]interface{}
}

// newWrappedError
func newWrappedError(prev error, text string, file string, line int, references ...ErrorReference) *wrappedError {
	w := &wrappedError{
		prev: prev,
		file: file,
		line: line,
		text: text,
		refs: make(map[string]interface{}),
	}
	if n := len(references); n > 0 {
		for i := 0; i < len(references); i++ {
			references[i](w)
		}
	}
	return w
}

func (w *wrappedError) Error() string {
	return w.text
}

// Contain returns a boolean result
// if the target error is included in the chain errors.
//
// See errors.Is
func Contain(err error, target error) bool {
	if target == nil {
		return err == target
	}

	cb := reflect.TypeOf(target).Comparable()
	for {
		if cb && err == target {
			return true
		}
		w, ok := isWrappedError(err)
		if !ok {
			return false
		}
		err = w.prev
	}
}
