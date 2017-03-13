package probe

import (
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
	retentionPeriod int
	stats           map[HostName]map[Protocol][]Measurement
}

func NewStatistics(retentionPeriod int) *Statistics {
	return &Statistics{stats: map[HostName]map[Protocol][]Measurement{}, retentionPeriod: retentionPeriod}
}

func (s *Statistics) Add(host HostName, protocol Protocol, kind string, value float64) {
	s.lock.Lock()
	defer s.lock.Unlock()
	// TODO: Enforce retention period
	if _, ok := s.stats[host]; !ok {
		s.stats[host] = make(map[Protocol][]Measurement)
	}

	s.stats[host][protocol] = append(s.stats[host][protocol],
		Measurement{Kind: kind, TimeStamp: time.Now().UnixNano(), Value: value})
}

func (s *Statistics) Dump() map[HostName]map[Protocol][]Measurement {
	s.lock.Lock()
	defer s.lock.Unlock()
	result := s.stats
	s.stats = map[HostName]map[Protocol][]Measurement{}
	return result
}
