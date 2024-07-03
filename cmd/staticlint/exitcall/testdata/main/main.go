// Package main implements tests package main with direct call os.Exit
package main

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
