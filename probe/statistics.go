package probe

import (
	"sync"
	"time"
)

type Protocol string

type Measurement struct {
	kind      string
	timeStamp int64
	value     float64
}

type Statistics struct {
	lock            sync.Mutex
	retentionPeriod int
	stats           map[HostName]map[Protocol][]Measurement
}

func NewStatistics(retentionPeriod int) *Statistics {
	return &Statistics{retentionPeriod: retentionPeriod}
}

func (s *Statistics) Add(host HostName, protocol Protocol, kind string, value float64) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.stats[host][protocol] = append(s.stats[host][protocol],
		Measurement{kind: kind, timeStamp: time.Now().UnixNano(), value: value})
}
