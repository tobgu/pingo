package probe

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
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

func startProbes(config Config, host Host, statistics *Statistics) {
	if host.TcpPort != 0 {
		startTcpProbe(config, host, statistics)
	}

	if host.UdpPort != 0 {
		startUdpProbe(config, host, statistics)
	}

	// TODO: ICMP
}

func startUdpProbe(config Config, host Host, statistics *Statistics) {
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

func startTcpProbe(config Config, host Host, statistics *Statistics) {
	inputBytes := make([]byte, config.TcpSize)
	rand.Read(inputBytes)

	doProbe := func() {
		// TODO: Categorize error cases and store the correct metrics
		// TODO: Proper timeouts on the connections
		start := time.Now()
		con, err := net.Dial("tcp", fmt.Sprintf("%s:%d", host.Address, host.TcpPort))
		if err != nil {
			log.Println("Error dialing TCP", err)
			return
		}

		inputBuffer := bytes.NewBuffer(inputBytes)
		_, err = io.Copy(con, inputBuffer)
		if err != nil {
			log.Println("Error writing TCP buffer", err)
			return
		}

		if tcpcon, ok := con.(*net.TCPConn); ok {
			tcpcon.CloseWrite()
		} else {
			log.Println("Unexpected TCP connection type")
			return
		}

		outputBuffer := bytes.Buffer{}
		_, err = io.Copy(&outputBuffer, con)
		if err != nil {
			log.Println("Error reading TCP", err)
			return
		}

		if !byteArrayEquals(outputBuffer.Bytes(), inputBytes) {
			log.Println("Input and output differs")
			return
		}

		err = con.Close()
		if err != nil {
			log.Println("Error closing Connection", err)
			return
		}

		duration := time.Now().Sub(start).Seconds()
		statistics.Add(host.Name, "tcp", "ping_success", duration)
	}

	executeAtInterval(doProbe, config.PingInterval)
}

func Run(config Config) error {
	s := NewStatistics(config.StatisticsRetentionPeriod)
	for _, h := range config.Hosts {
		startProbes(config, h, s)
	}

	fmt.Println("Running probes", config)

	// TODO: Statistics HTTP API
	select {}
	return nil
}
