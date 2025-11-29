package main

import (
	"fmt"
	"os"
)

var errCannotOpenFile = fmt.Errorf("cannot open file")

// OpenFile attempts to open a file at the given path and returns a Result containing the *os.File or an error.
func OpenFile(path string) Result[*os.File] {
	f, err := os.Open(path)
	if err != nil {
		return Err[*os.File](fmt.Errorf("%w: %s", errCannotOpenFile, err.Error()))
	}
	return Ok(f)
}

func main() {

	// ------------------------------------------------------------
	// 1. Existing file: OK branch
	// ------------------------------------------------------------
	res := OpenFile("file.txt")
	if res.IsOk() {
		fmt.Println("OK: file exists")
	}

	// Unwrap -> returns *os.File or panic. ok in this case
	file := res.Unwrap()
	fmt.Println("Opened:", file.Name())
	defer file.Close()

	// ------------------------------------------------------------
	// 2. Non-existing file: Error branch
	// ------------------------------------------------------------
	res = OpenFile("unknown.txt")
	if res.IsErr() {
		fmt.Fprintln(os.Stderr, "ERROR:", res.Error())
	}

	// UnwrapOrErr -> safe alternative to Unwrap
	// returns (T, error)
	_, err := res.UnwrapOrErr(
		fmt.Errorf("fatal: cannot read unknown.txt"),
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, "UnwrapOrErr:", err.Error())
	}

	// ------------------------------------------------------------
	// 3. Dangerous branch: Unwrap / Expect
	// ------------------------------------------------------------

	// This will panic if the file does not exist.
	// OpenFile("unknown.txt").Unwrap()

	// This will panic with a custom message.
	// OpenFile("unknown.txt").Expect("file does not exist")
}
