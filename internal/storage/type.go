package storage

type MetricsType = string

const (
	GuageType   = MetricsType("guage")
	CounterType = MetricsType("counter")
)

func IsValidMetricsType(value string) bool {
	switch value {
	case GuageType, CounterType:
		return true
	default:
		return false
	}
}

type GuageDateType = float64
type GuageName = string

const (
	Alloc         = GuageName("Alloc")
	BuckHashSys   = GuageName("BuckHashSys")
	Frees         = GuageName("Frees")
	GCCPUFraction = GuageName("GCCPUFraction")
	GCSys         = GuageName("GCSys")
	HeapAlloc     = GuageName("HeapAlloc")
	HeapIdle      = GuageName("HeapIdle")
	HeapInuse     = GuageName("HeapInuse")
	HeapObjects   = GuageName("HeapObjects")
	HeapReleased  = GuageName("HeapReleased")
	HeapSys       = GuageName("HeapSys")
	LastGC        = GuageName("LastGC")
	Lookups       = GuageName("Lookups")
	MCacheInuse   = GuageName("MCacheInuse")
	MCacheSys     = GuageName("MCacheSys")
	MSpanInuse    = GuageName("MSpanInuse")
	MSpanSys      = GuageName("MSpanSys")
	Mallocs       = GuageName("Mallocs")
	NextGC        = GuageName("NextGC")
	NumForcedGC   = GuageName("NumForcedGC")
	NumGC         = GuageName("NumGC")
	OtherSys      = GuageName("OtherSys")
	PauseTotalNs  = GuageName("PauseTotalNs")
	StackInuse    = GuageName("StackInuse")
	StackSys      = GuageName("StackSys")
	Sys           = GuageName("Sys")
	TotalAlloc    = GuageName("TotalAlloc")
	RandomValue   = GuageName("RandomValue")
)

func IsValidGuageName(value string) bool {
	switch value {
	case Alloc, BuckHashSys, Frees, GCCPUFraction, GCSys, HeapAlloc, HeapIdle, HeapInuse, HeapObjects, HeapReleased,
		HeapSys, LastGC, Lookups, MCacheInuse, MCacheSys, MSpanInuse, MSpanSys, Mallocs, NextGC, NumForcedGC, NumGC,
		OtherSys, PauseTotalNs, StackInuse, StackSys, Sys, TotalAlloc, RandomValue:
		return true
	default:
		return false
	}
}

type CounterDateType = int64
type CounterName = string

const (
	PollCount = CounterName("PollCount")
)

func IsValidCounterName(value string) bool {
	switch value {
	case PollCount:
		return true
	}

	return false
}

type MemStorage struct {
	Guages   map[GuageName]GuageDateType
	Counters map[CounterName]CounterDateType
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		Guages:   make(map[GuageName]GuageDateType),
		Counters: make(map[CounterName]CounterDateType),
	}
}

func (s *MemStorage) UpdateGuageMetric(name GuageName, value GuageDateType) {
	s.Guages[name] = value
}

func (s *MemStorage) UpdateCounterMetric(name CounterName, value CounterDateType) {
	s.Counters[name] += value
}

type MetricsState struct {
	Guages   map[GuageName]GuageDateType
	Counters map[CounterName]CounterDateType
}

func NewMetricsState() *MetricsState {
	return &MetricsState{
		Guages:   make(map[GuageName]GuageDateType),
		Counters: make(map[CounterName]CounterDateType),
	}
}
