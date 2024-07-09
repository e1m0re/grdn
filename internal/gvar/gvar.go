package gvar

import "fmt"

var (
	BuildVersion = "N/A"
	BuildDate    = "N/A"
	BuildCommit  = "N/A"
)

func prepareWelcome() string {
	return fmt.Sprintf(
		"Build version: %s\nBuild date: %s\nBuild commit: %s",
		BuildVersion,
		BuildDate,
		BuildCommit,
	)
}

// PrintWelcome prints to STDOUT welcome message.
func PrintWelcome() {
	fmt.Println(prepareWelcome())
}
