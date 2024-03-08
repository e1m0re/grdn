package main

import (
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
	data.Guages[storage.RandomValue] = r.Float64()

	// memory metrics
	var rtm runtime.MemStats
	runtime.ReadMemStats(&rtm)
	data.Guages[storage.Alloc] = storage.GuageDateType(rtm.Alloc)
	data.Guages[storage.BuckHashSys] = storage.GuageDateType(rtm.BuckHashSys)
	data.Guages[storage.Frees] = storage.GuageDateType(rtm.Frees)
	data.Guages[storage.GCCPUFraction] = rtm.GCCPUFraction
	data.Guages[storage.GCSys] = storage.GuageDateType(rtm.GCSys)
	data.Guages[storage.HeapAlloc] = storage.GuageDateType(rtm.HeapAlloc)
	data.Guages[storage.HeapIdle] = storage.GuageDateType(rtm.HeapIdle)
	data.Guages[storage.HeapInuse] = storage.GuageDateType(rtm.HeapInuse)
	data.Guages[storage.HeapObjects] = storage.GuageDateType(rtm.HeapObjects)
	data.Guages[storage.HeapReleased] = storage.GuageDateType(rtm.HeapReleased)
	data.Guages[storage.HeapSys] = storage.GuageDateType(rtm.HeapSys)
	data.Guages[storage.LastGC] = storage.GuageDateType(rtm.LastGC)
	data.Guages[storage.Lookups] = storage.GuageDateType(rtm.Lookups)
	data.Guages[storage.MCacheInuse] = storage.GuageDateType(rtm.MCacheInuse)
	data.Guages[storage.MCacheSys] = storage.GuageDateType(rtm.MCacheSys)
	data.Guages[storage.MSpanInuse] = storage.GuageDateType(rtm.MSpanInuse)
	data.Guages[storage.MSpanSys] = storage.GuageDateType(rtm.MSpanSys)
	data.Guages[storage.Mallocs] = storage.GuageDateType(rtm.Mallocs)
	data.Guages[storage.NextGC] = storage.GuageDateType(rtm.NextGC)
	data.Guages[storage.NumForcedGC] = storage.GuageDateType(rtm.NumForcedGC)
	data.Guages[storage.NumGC] = storage.GuageDateType(rtm.NumGC)
	data.Guages[storage.OtherSys] = storage.GuageDateType(rtm.OtherSys)
	data.Guages[storage.StackInuse] = storage.GuageDateType(rtm.StackInuse)
	data.Guages[storage.StackSys] = storage.GuageDateType(rtm.StackSys)
	data.Guages[storage.PauseTotalNs] = storage.GuageDateType(rtm.PauseTotalNs)
	data.Guages[storage.Sys] = storage.GuageDateType(rtm.Sys)
	data.Guages[storage.TotalAlloc] = storage.GuageDateType(rtm.TotalAlloc)

}

func sendMetric(mType string, mName string, mValue string) (err error) {
	var u = url.URL{
		Scheme: "http",
		Host:   "localhost:8080",
		Path:   fmt.Sprintf("/update/%s/%s/%s", mType, mName, mValue),
	}
	_, err = http.Post(u.String(), "text/plan", nil)
	return
}

func SendData(data *storage.MetricsState) {
	for key, value := range data.Guages {
		err := sendMetric(storage.GuageType, key, fmt.Sprintf("%v", value))
		if err != nil {
			fmt.Printf("%s\r\n", err)
		}
	}
	for key, value := range data.Counters {
		err := sendMetric(storage.CounterType, key, fmt.Sprintf("%v", value))
		fmt.Printf("%s\r\n", err)
	}
}

func main() {
	pollInterval := time.Duration(2) * time.Second
	reportInterval := time.Duration(8) * time.Second
	state := storage.NewMetricsState()

	for {
		<-time.After(pollInterval)
		UpdateMetrics(state)
		<-time.After(reportInterval)
		SendData(state)
	}
}
