package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"runtime"
	"time"

	"github.com/e1m0re/grdn/internal/storage"
)

func UpdateMetrics(data *storage.MetricsState) {
	data.Counters[storage.PollCount]++

	r := rand.New(rand.NewSource(data.Counters[storage.PollCount]))
	data.Gauges[storage.RandomValue] = r.Float64()

	// memory metrics
	var rtm runtime.MemStats
	runtime.ReadMemStats(&rtm)
	data.Gauges[storage.Alloc] = storage.GaugeDateType(rtm.Alloc)
	data.Gauges[storage.BuckHashSys] = storage.GaugeDateType(rtm.BuckHashSys)
	data.Gauges[storage.Frees] = storage.GaugeDateType(rtm.Frees)
	data.Gauges[storage.GCCPUFraction] = rtm.GCCPUFraction
	data.Gauges[storage.GCSys] = storage.GaugeDateType(rtm.GCSys)
	data.Gauges[storage.HeapAlloc] = storage.GaugeDateType(rtm.HeapAlloc)
	data.Gauges[storage.HeapIdle] = storage.GaugeDateType(rtm.HeapIdle)
	data.Gauges[storage.HeapInuse] = storage.GaugeDateType(rtm.HeapInuse)
	data.Gauges[storage.HeapObjects] = storage.GaugeDateType(rtm.HeapObjects)
	data.Gauges[storage.HeapReleased] = storage.GaugeDateType(rtm.HeapReleased)
	data.Gauges[storage.HeapSys] = storage.GaugeDateType(rtm.HeapSys)
	data.Gauges[storage.LastGC] = storage.GaugeDateType(rtm.LastGC)
	data.Gauges[storage.Lookups] = storage.GaugeDateType(rtm.Lookups)
	data.Gauges[storage.MCacheInuse] = storage.GaugeDateType(rtm.MCacheInuse)
	data.Gauges[storage.MCacheSys] = storage.GaugeDateType(rtm.MCacheSys)
	data.Gauges[storage.MSpanInuse] = storage.GaugeDateType(rtm.MSpanInuse)
	data.Gauges[storage.MSpanSys] = storage.GaugeDateType(rtm.MSpanSys)
	data.Gauges[storage.Mallocs] = storage.GaugeDateType(rtm.Mallocs)
	data.Gauges[storage.NextGC] = storage.GaugeDateType(rtm.NextGC)
	data.Gauges[storage.NumForcedGC] = storage.GaugeDateType(rtm.NumForcedGC)
	data.Gauges[storage.NumGC] = storage.GaugeDateType(rtm.NumGC)
	data.Gauges[storage.OtherSys] = storage.GaugeDateType(rtm.OtherSys)
	data.Gauges[storage.StackInuse] = storage.GaugeDateType(rtm.StackInuse)
	data.Gauges[storage.StackSys] = storage.GaugeDateType(rtm.StackSys)
	data.Gauges[storage.PauseTotalNs] = storage.GaugeDateType(rtm.PauseTotalNs)
	data.Gauges[storage.Sys] = storage.GaugeDateType(rtm.Sys)
	data.Gauges[storage.TotalAlloc] = storage.GaugeDateType(rtm.TotalAlloc)

}

func sendMetric(mType string, mName string, mValue string) (err error) {
	var u = url.URL{
		Scheme: "http",
		Host:   agentOptions.flagRunAddr,
		Path:   fmt.Sprintf("/update/%s/%s/%s", mType, mName, mValue),
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

func SendData(data *storage.MetricsState) {
	for key, value := range data.Gauges {
		err := sendMetric(storage.GaugeType, key, fmt.Sprintf("%v", value))
		if err != nil {
			fmt.Printf("%v\r\n", err)
		}
	}
	for key, value := range data.Counters {
		err := sendMetric(storage.CounterType, key, fmt.Sprintf("%v", value))
		if err != nil {
			fmt.Printf("%v\r\n", err)
		}
	}
}

func main() {
	parseFlags()

	pollInterval := time.Duration(agentOptions.pollInterval) * time.Second
	reportInterval := time.Duration(agentOptions.reportInterval) * time.Second
	state := storage.NewMetricsState()

	for {
		<-time.After(pollInterval)
		UpdateMetrics(state)
		<-time.After(reportInterval)
		SendData(state)
	}
}
