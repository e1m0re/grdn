// Package a implements tests package a with direct call os.Exit
package a

import "os"

func main() {
	defer func() {
		os.Exit(0)
	}()

	counter := 1
	if counter > 0 {
		os.Exit(0)
	}

	bar := func() {
		os.Exit(0)
	}
	bar()

	os.Exit(0)
}
