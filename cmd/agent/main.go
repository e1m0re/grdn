package main

import (
	"time"

	"github.com/e1m0re/grdn/internal/apiclient"
	"github.com/e1m0re/grdn/internal/monitor"
)

func main() {
	parseFlags()

	pollInterval := time.Duration(agentOptions.pollInterval) * time.Second
	reportInterval := time.Duration(agentOptions.reportInterval) * time.Second
	api := apiclient.NewAPI("http://" + agentOptions.flagRunAddr)
	monitor1 := monitor.NewMetricsMonitor()

	for {
		<-time.After(pollInterval)
		monitor1.UpdateData()
		<-time.After(reportInterval)
		monitor1.SendDataToServers(api)
	}
}
