package probe

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/http"
	"syscall"
	"time"
)

type HostName string

type Host struct {
	Name    HostName `yaml:"name"`
	Address string   `yaml:"address"`
	UdpPort int      `yaml:"udp_port"`
	TcpPort int      `yaml:"tcp_port"`
	Icmp    bool     `yaml:"icmp"`
}

type Config struct {
	StatisticsPort            int    `yaml:"statistics_port"`
	StatisticsRetentionPeriod int    `yaml:"statistics_retention_period"`
	PingInterval              int    `yaml:"ping_interval"`
	ConnectionTimeout         int    `yaml:"connection_timeout"`
	ReadTimeout               int    `yaml:"read_timeout"`
	TcpSize                   int    `yaml:"tcp_size"`
	UdpSize                   int    `yaml:"udp_size"`
	IcmpSize                  int    `yaml:"icmp_port"`
	Hosts                     []Host `yaml:"hosts"`
}

type probe struct {
	config     Config
	host       Host
	statistics *Statistics
}

func startProbes(config Config, host Host, statistics *Statistics) {
	if host.TcpPort != 0 {
		p := newTcp(config, host, statistics)
		executeAtInterval(p.execute, config.PingInterval)
	}

	if host.UdpPort != 0 {
		p := newUdp(config, host, statistics)
		executeAtInterval(p.execute, config.PingInterval)
	}

	// TODO: ICMP, HTTP
}

func executeAtInterval(f func(), interval int) {
	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				f()
			}
		}
	}()
}

func byteArrayEquals(a []byte, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

func classifyError(err error) string {
	errorStr := ""

	switch t := err.(type) {
	case *net.OpError:
		errorStr = fmt.Sprintf("%s_error", t.Op)
	case syscall.Errno:
		if t == syscall.ECONNREFUSED {
			errorStr = "connection_refused_error"
		} else {
			errorStr = "unknown_syscall_error"
		}
	default:
		errorStr = "unknown_error"
	}

	if netError, ok := err.(net.Error); ok && netError.Timeout() {
		errorStr = fmt.Sprintf("%s_timeout", errorStr)
	}

	log.Println("Network error,", errorStr, ",", err)
	return errorStr
}

func randomBytes(size int) []byte {
	b := make([]byte, size)
	rand.Read(b)
	return b
}

func Run(config Config) error {
	stats := NewStatistics(config.StatisticsRetentionPeriod)

	log.Println("Starting probes", config)
	for _, host := range config.Hosts {
		startProbes(config, host, stats)
	}

	log.Println("Serving statistics API at port", config.StatisticsPort)
	apiHandler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(stats.Dump())
	}

	http.HandleFunc("/statistics", apiHandler)
	return http.ListenAndServe(fmt.Sprintf(":%d", config.StatisticsPort), nil)
}
