package main

import (
	"flag"
	"os"
	"strconv"
	"time"

	"github.com/e1m0re/grdn/internal/apiclient"
	"github.com/e1m0re/grdn/internal/monitor"
)

func main() {
	var flagRunAddr string

	flag.StringVar(&flagRunAddr, "a", "localhost:8080", "address and port to run server")

	envRunAddr := os.Getenv("ADDRESS")
	if envRunAddr != "" {
		flagRunAddr = envRunAddr
	}

	var reportInterval uint

	flag.UintVar(&reportInterval, "r", 10, "frequency of sending metrics to the server")

	envReportInterval := os.Getenv("REPORT_INTERVAL")
	if envReportInterval != "" {
		envValue, err := strconv.Atoi(envReportInterval)
		if err == nil {
			reportInterval = uint(envValue)
		}
	}

	var pollInterval uint

	flag.UintVar(&pollInterval, "p", 2, "frequency of polling metrics from the package")

	envPollInterval := os.Getenv("POLL_INTERVAL")
	if envPollInterval != "" {
		envValue, err := strconv.Atoi(envPollInterval)
		if err == nil {
			pollInterval = uint(envValue)
		}
	}

	flag.Parse()

	api := apiclient.NewAPI("http://" + flagRunAddr)
	monitor1 := monitor.NewMetricsMonitor()

	for {
		<-time.After(time.Duration(pollInterval) * time.Second)
		monitor1.UpdateData()
		<-time.After(time.Duration(reportInterval) * time.Second)
		monitor1.SendDataToServers(api)
	}
}
