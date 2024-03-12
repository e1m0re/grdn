package main

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/e1m0re/grdn/internal/monitor"
)

func doRequest(uriPath string) (err error) {
	var u = url.URL{
		Scheme: "http",
		Host:   agentOptions.flagRunAddr,
		Path:   uriPath,
	}

	var buf []byte
	responseBody := bytes.NewBuffer(buf)

	resp, err := http.Post(u.String(), "text/plan", responseBody)

	if err != nil {
		return
	}

	defer resp.Body.Close()

	return
}

func SendData(monitor1 *monitor.MetricsMonitor) {
	list := monitor1.GetData()
	for _, row := range list {
		err := doRequest(fmt.Sprintf("/update/%s/%s/%s", row.MType, row.MName, row.MValue))
		if err != nil {
			fmt.Printf("%v\r\n", err)
		}
	}
}

func main() {
	parseFlags()

	pollInterval := time.Duration(agentOptions.pollInterval) * time.Second
	reportInterval := time.Duration(agentOptions.reportInterval) * time.Second
	monitor1 := monitor.NewMetricsMonitor()

	for {
		<-time.After(pollInterval)
		monitor1.UpdateData()
		<-time.After(reportInterval)
		SendData(monitor1)
	}
}
