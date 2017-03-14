package probe

import (
	"log"
	"sync"
	"time"
)

type Protocol string

type Measurement struct {
	Kind      string  `json:"kind"`
	TimeStamp int64   `json:"timestamp"`
	Value     float64 `json:"value"`
}

type Statistics struct {
	lock            sync.Mutex
	retentionPeriod time.Duration
	stats           map[HostName]map[Protocol][]Measurement
}

func NewStatistics(retentionPeriod int) *Statistics {
	return &Statistics{
		stats:           map[HostName]map[Protocol][]Measurement{},
		retentionPeriod: time.Duration(retentionPeriod) * time.Second}
}

func (s *Statistics) Add(host HostName, protocol Protocol, kind string, value float64) {
	s.lock.Lock()
	defer s.lock.Unlock()
	if _, ok := s.stats[host]; !ok {
		s.stats[host] = make(map[Protocol][]Measurement)
	}

	stats := s.stats[host][protocol]
	if len(stats) > 0 {
		statAge := time.Now().Sub(time.Unix(0, stats[0].TimeStamp))
		if statAge > s.retentionPeriod {
			// Throw away half of the stats to free up memory
			log.Println("Truncating too old statistics for", protocol, kind)
			newLen := len(stats) / 2
			newStats := make([]Measurement, newLen, newLen+1)
			copy(newStats, stats[newLen:])
			stats = newStats
		}
	}

	s.stats[host][protocol] = append(stats, Measurement{Kind: kind, TimeStamp: time.Now().UnixNano(), Value: value})
}

func (s *Statistics) Dump() map[HostName]map[Protocol][]Measurement {
	s.lock.Lock()
	defer s.lock.Unlock()
	result := s.stats
	s.stats = map[HostName]map[Protocol][]Measurement{}
	return result
}
