# errors

A thin facade over the Go standard library’s `errors` and `fmt` packages. It centralizes how you construct and wrap errors while keeping normal unwrap behavior so `errors.Is` and `errors.As` (and this package’s `Is` / `As`) work on the full chain.

## Install

```bash
go get github.com/thecxx/errors
```

The import path is `github.com/thecxx/errors`, and the **package name is `errors`**, the same as the standard library. If one file needs both, alias the **standard** library (not this module):

```go
import (
	stderrors "errors"

	"github.com/thecxx/errors"
)
```

Use `errors.New` / `errors.Wrap` / … for this package, and `stderrors.Is` only when you must call the standard API by name. Often `errors.Is` from this module is enough because it forwards to the standard library.

## API

| Function | Delegates to | Typical use |
|----------|----------------|-------------|
| `New(message string) error` | `errors.New` | Simple sentinel or message-only errors. |
| `Newf(format string, args ...any) error` | `fmt.Errorf(format, args...)` | Formatted messages; use `%w` in the format when you need a wrapped cause. |
| `Wrap(err error, message string) error` | `fmt.Errorf("%s: %w", message, err)` | Add a fixed human-readable prefix in front of an existing error; always wraps with `%w`. |
| `Unwrap(err error) error` | `errors.Unwrap` | Walk the chain manually. |
| `Is(err, target error) bool` | `errors.Is` | Compare against sentinels or types in the chain. |
| `As(err error, target any) bool` | `errors.As` | Extract a typed error from the chain. |

## `Wrap` vs `Newf`

- **`Wrap(err, "doing X")`** — consistent prefix style: `"doing X: " + err`. The wrapped error is always the root of the `%w` clause, so unwrap semantics are predictable.
- **`Newf("user %d: %w", id, err)`** — use when the message is format-driven or the wrapped error sits inside a longer format string.

Both preserve unwrap when you use `%w` (`Newf`) or `Wrap` (which always uses `%w` internally).

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

## Why use this package

It does not add new error types or stack traces; it is only a single import point and naming consistency for teams that want wrappers like `Wrap` without repeating `fmt.Errorf("%s: %w", ...)` everywhere. If you only need the standard library, `errors` + `fmt` is enough—this package is optional sugar and documentation of conventions.
