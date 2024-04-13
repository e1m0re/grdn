package main

import (
	"time"

	"github.com/e1m0re/grdn/internal/apiclient"
	"github.com/e1m0re/grdn/internal/monitor"
)

func main() {
	options := config()

	apiClient := apiclient.NewAPI("http://" + options.serverAddr)
	monitor1 := monitor.NewMetricsMonitor(apiClient)

	for {
		<-time.After(time.Duration(options.pollInterval) * time.Second)
		monitor1.UpdateData()
		<-time.After(time.Duration(options.reportInterval) * time.Second)
		monitor1.SendDataToServer()
	}
}
