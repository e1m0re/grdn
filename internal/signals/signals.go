package signals

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

func HandleOSSignals(cancel context.CancelFunc) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)

	<-c
	cancel()
}
