package main

import "flag"

type Options struct {
	flagRunAddr    string
	reportInterval uint
	pollInterval   uint
}

var agentOptions = new(Options)

func parseFlags() {
	flag.StringVar(&agentOptions.flagRunAddr, "a", "localhost:8080", "address and port to run server")
	flag.UintVar(&agentOptions.reportInterval, "r", 10, "frequency of sending metrics to the server")
	flag.UintVar(&agentOptions.pollInterval, "p", 2, "frequency of polling metrics from the package")
	flag.Parse()
}
