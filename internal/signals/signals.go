package signals

import (
	"os"
	"os/signal"
	"syscall"
)

func HandleOSSignals(f func()) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)

	<-c
	f()
}
