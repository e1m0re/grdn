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

func PrintWelcome() {
	fmt.Println(prepareWelcome())
}
