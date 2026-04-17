# errors

A thin facade over the Go standard library’s `errors` and `fmt` packages. It centralizes how you construct and wrap errors while keeping normal unwrap behavior so `errors.Is` and `errors.As` (and this package’s `Is` / `As`) work on the full chain.

Errors produced by this package’s constructors (`New`, `Newf`, `Wrap`, `Join`, `WithValue`, `WithValues`) are wrapped in a small internal type that records the **caller’s program counter** (via [`github.com/thecxx/runpoint`](https://github.com/thecxx/runpoint)) and optional **string-keyed metadata**. `Error` and `Unwrap` match the underlying error; wrapping does not break `%w` semantics for `Is` / `As`.

## Install

```bash
go get github.com/thecxx/errors
```

This module depends on `github.com/thecxx/runpoint` (Go 1.18+).

The import path is `github.com/thecxx/errors`, and the **package name is `errors`**, the same as the standard library. If one file needs both, alias the **standard** library (not this module):

```go
import (
	stderrors "errors"

	"github.com/thecxx/errors"
)
```

Use `errors.New` / `errors.Wrap` / … for this package, and `stderrors.Is` only when you must call the standard API by name. Often `errors.Is` from this module is enough because it forwards to the standard library.

## API

| Function | Behavior | Typical use |
|----------|----------|-------------|
| `New(message string) error` | `errors.New`, then an outer shell with PC | Simple sentinel or message-only errors |
| `Newf(format string, args ...any) error` | `fmt.Errorf(format, args...)`, then shell with PC | Formatted messages; use `%w` when you need a wrapped cause |
| `Wrap(err error, message string) error` | `fmt.Errorf("%s: %w", message, err)` plus shell | Fixed human-readable prefix; predictable `%w` chain |
| `Join(errs ...error) error` | `errors.Join`; if non-nil, wrapped with PC; supports `Unwrap() []error` | Combine several errors (Go 1.20+ join semantics) |
| `WithValue(err error, key string, value any) error` | Attach one key/value, keep `err`, record PC | Request IDs, user IDs, etc. |
| `WithValues(err error, values map[string]any) error` | Attach many keys at once | Bulk metadata |
| `Unwrap(err error) error` | `errors.Unwrap` | Walk the single-error chain one step |
| `PC(err error) *runpoint.PCounter` | Returns PC if `err` exposes it, else `nil` | See where the error was constructed/wrapped |
| `Value(err error, key string) any` | Reads `Value(string) any`; this package’s shell checks local map then delegates inward | Read attached metadata |
| `Is(err, target error) bool` | `errors.Is` | Match sentinels or types in the chain |
| `As(err error, target any) bool` | `errors.As` | Extract a typed error from the chain |

Errors created only with the standard library usually won’t carry this package’s outer shell; `Value` can still work for custom types that implement the `Value(string) any` convention.

## `Wrap` vs `Newf`

- **`Wrap(err, "doing X")`** — consistent prefix style: `"doing X: " + err`. The wrapped error is always the root of the `%w` clause, so unwrap semantics are predictable.
- **`Newf("user %d: %w", id, err)`** — use when the message is format-driven or the wrapped error sits inside a longer format string.

Both preserve unwrap when you use `%w` (`Newf`) or `Wrap` (which always uses `%w` internally). The extra outer shell is for PC/metadata only and does not change `Is` / `As` over the underlying chain.

## Example

```go
package main

import (
	"fmt"

	"github.com/thecxx/errors"
)

func main() {
	base := errors.New("not found")

	chained := errors.Wrap(base, "load config")
	fmt.Println(chained)
	// load config: not found

	fmt.Println(errors.Is(chained, base)) // true
}
```

Wrapping with `Newf`:

```go
err := doThing()
if err != nil {
	return errors.Newf("doThing(id=%d): %w", id, err)
}
```

Attach and read metadata (local map first, then delegate to inner errors that implement `Value`):

```go
err := errors.WithValue(errors.New("denied"), "user_id", 42)
if v := errors.Value(err, "user_id"); v != nil {
	fmt.Println(v)
}
```

Using `Join`:

```go
e := errors.Join(errOpen, errParse)
if e != nil {
	return errors.Wrap(e, "load")
}
```

## Why use this package

Beyond avoiding repeated `fmt.Errorf("%s: %w", ...)` at call sites, you get a **single import** for constructors, optional **runpoint** caller PCs, and lightweight **key/value context** (`WithValue` / `WithValues` + `Value`) so logs and telemetry can align with “who returned this error” without pulling in a heavy stack-capture library.

If you only need the standard library, `errors` + `fmt` is enough—this package is optional sugar and documented conventions. Custom error types can participate in `Value` by implementing `Value(string) any` (and may expose `PC() *runpoint.PCounter` if you want the same introspection hooks as this package’s shell).
