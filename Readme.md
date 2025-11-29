# Result Pattern in Go — A Rust-Style Result[T] for Safer Error Handling

## Why This Project Exists (Cloudflare 2025 Incident)

On **November 18th, 2025**, Cloudflare suffered a major global outage caused by
a single Rust `unwrap()` in production code.

The chain was simple but catastrophic:

- an unexpected error occurred,
- `unwrap()` triggered a panic,
- the process aborted instantly,
- the panic cascaded across services,
- large parts of Cloudflare's infrastructure went down.

Official postmortem:
https://blog.cloudflare.com/18-november-2025-outage/

This incident made me stop and think:

> *How can a single `unwrap()` bring down a global system?  
> What does Rust's Result/Option pattern enforce, and how does that differ from Go's error handling?*

This repository was created as an **educational exercise** to explore:

- why `unwrap()` is intentionally dangerous in Rust,
- what the Result pattern really is,
- how safe vs unsafe error paths behave,
- and how the same ideas can be expressed in **Go**.

The purpose is not to replace Go's error model but to illustrate the philosophy
of explicit error handling, using Go as a medium for learning.

---
This repository demonstrates a clean and idiomatic implementation of the **Result Pattern** in Go, inspired by Rust’s:

- `Ok(value)`
- `Err(error)`
- `IsOk()`
- `IsErr()`
- `Unwrap()`
- `Expect(msg)`
- `UnwrapOr(default)`
- `UnwrapOrErr(newErr)`

It shows how to write safer and more expressive code without abandoning Go’s simplicity.

---

#  Why the Result Pattern?

In Go, errors are commonly handled like this:

```go
value, err := f()
if err != nil {
    return ..., err
}
```

It works well — but it can become:

repetitive
noisy
easy to misuse (_ = err)
hard to scale into abstractions
error-prone in pipelines
Languages like Rust, Swift, Scala, Haskell, Elm, Zig, Gleam use a structured pattern:

`Result = Ok(value) OR Err(error)`

The caller must explicitly handle both cases.
This repo brings that expressiveness and clarity to Go.

Example Usage (main.go)

```go 
package main

import (
	"fmt"
	"os"

	"github.com/tzzed/go-result/result"
)

var ErrCannotOpenFile = fmt.Errorf("cannot open file")

// OpenFile Rust-like wrapper
func OpenFile(path string) result.Result[*os.File] {
	f, err := os.Open(path)
	if err != nil {
		return result.Err[*os.File](
			fmt.Errorf("%w: %s", ErrCannotOpenFile, err.Error()),
		)
	}
	return result.Ok(f)
}

func main() {

	// ------------------------------------------------------------
	// 1. Existing file: OK branch
	// ------------------------------------------------------------
	res := OpenFile("file.txt")
	if res.IsOk() {
		fmt.Println("OK: file exists")
	}

	file := res.Unwrap() // returns *os.File or panic
	fmt.Println("Opened:", file.Name())
	defer file.Close()

	// ------------------------------------------------------------
	// 2. Non-existing file: Error branch
	// ------------------------------------------------------------
	res = OpenFile("unknown.txt")
	if res.IsErr() {
		_, _ = fmt.Fprintln(os.Stderr, "ERROR:", res.Error())
	}

	// UnwrapOrErr → safe alternative to Unwrap
	_, err := res.UnwrapOrErr(
		fmt.Errorf("fatal: cannot read unknown.txt"),
	)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "UnwrapOrErr:", err.Error())
	}

	// ------------------------------------------------------------
	// 3. Dangerous branch: Unwrap / Expect
	// ------------------------------------------------------------

	// This will panic if the file does not exist.
	// OpenFile("unknown.txt").Unwrap()

	// This will panic with a custom message.
	// OpenFile("unknown.txt").Expect("file does not exist")
}
```

Why Unwrap() Is Dangerous — Cloudflare 2025 Incident

Rust’s unwrap() is powerful — but also dangerous when misused.
In a major 2025 Cloudflare outage, a single line like this:

What happened? An unexpected error occurred `unwrap()` triggered a panic the process aborted instantly
a distributed failure cascaded through the system global infrastructure went down. This caused a global service crash.
![Cloudflare 18 November 2025 Outage graph](https://cf-assets.www.cloudflare.com/zkvhlag99gkb/640fjk9dawDk7f0wJ8Jm5S/668bcf1f574ae9e896671d9eee50da1b/BLOG-3079_7.png)

The lesson is universal:
Unchecked unwraps are production landmines.
Use them only when failure is probably impossible.

In this Go project:

`OpenFile("unknown.txt").Unwrap()`

it behaves the same way → panic. This is intentional.
It mirrors Rust’s behavior to teach good error discipline.

Safe Alternatives
#### UnwrapOr(default)
`port := ReadConfigPort().UnwrapOr(8080)`

#### UnwrapOrErr(newErr)
```go 
file, err := OpenFile("config.json").UnwrapOrErr(
fmt.Errorf("cannot load config"),
)
```

#### IsOk() / IsErr()
```go
if res.IsErr() {
	return res.Error()
}
```

### API Summary (Behavior Diagram)

```markdown
                            ┌──────────────┐
                            │ Result[T]    │
                            └───────┬──────┘
        ┌───────────────────────────┼────────────────────────┐
        │                           │                        │
        Ok(T)                   IsOk() → true               Unwrap() → T or panic!
        Err(error)              IsErr() → true              UnwrapOr(T) → T 
                                                            UnwrapOrErr(err) → (T, error)
                                                            Expect(msg) → panic(msg)
        
```