package probe

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
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

func startProbes(config Config, host Host, statistics *Statistics) {
	if host.TcpPort != 0 {
		startTcpProbe(config, host, statistics)
	}

	if host.UdpPort != 0 {
		startUdpProbe(config, host, statistics)
	}

	// TODO: ICMP
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

func startTcpProbe(config Config, host Host, statistics *Statistics) {
	inputBytes := randomBytes(config.TcpSize)

	addMetric := func(kind string, value float64) {
		statistics.Add(host.Name, "tcp", kind, value)
	}

	doProbe := func() {
		start := time.Now()
		con, err := net.DialTimeout("tcp",
			fmt.Sprintf("%s:%d", host.Address, host.TcpPort),
			time.Duration(config.ConnectionTimeout)*time.Second)
		if err != nil {
			addMetric(classifyError(err), 1.0)
			return
		}

		defer con.Close()

		err = con.SetDeadline(time.Now().Add(time.Duration(config.ReadTimeout) * time.Second))
		if err != nil {
			addMetric(classifyError(err), 1.0)
			return
		}

		inputBuffer := bytes.NewBuffer(inputBytes)
		_, err = io.Copy(con, inputBuffer)
		if err != nil {
			addMetric(classifyError(err), 1.0)
			return
		}

		if tcpcon, ok := con.(*net.TCPConn); ok {
			tcpcon.CloseWrite()
		} else {
			log.Println("Connection was not of type TCP, this was unexpected...")
			return
		}

		outputBuffer := bytes.Buffer{}
		_, err = io.Copy(&outputBuffer, con)
		if err != nil {
			addMetric(classifyError(err), 1.0)
			return
		}

		if !byteArrayEquals(outputBuffer.Bytes(), inputBytes) {
			addMetric("content_error", 1.0)
			return
		}

		duration := time.Now().Sub(start).Seconds()
		addMetric("ping_success", duration)
	}

	executeAtInterval(doProbe, config.PingInterval)
}

func startUdpProbe(config Config, host Host, statistics *Statistics) {
	inputBytes := randomBytes(config.UdpSize)

	addMetric := func(kind string, value float64) {
		statistics.Add(host.Name, "udp", kind, value)
	}

	doProbe := func() {
		start := time.Now()
		con, err := net.DialTimeout("udp",
			fmt.Sprintf("%s:%d", host.Address, host.UdpPort),
			time.Duration(config.ConnectionTimeout)*time.Second)
		if err != nil {
			addMetric(classifyError(err), 1.0)
			return
		}

		defer con.Close()

		err = con.SetDeadline(time.Now().Add(time.Duration(config.ReadTimeout) * time.Second))
		if err != nil {
			addMetric(classifyError(err), 1.0)
			return
		}

		inputBuffer := bytes.NewBuffer(inputBytes)
		_, err = io.Copy(con, inputBuffer)
		if err != nil {
			addMetric(classifyError(err), 1.0)
			return
		}

		udpCon, ok := con.(*net.UDPConn)
		if !ok {
			fmt.Println("Connection was not of type UDP, this was unexpected...")
			return
		}

		outputBytes := make([]byte, len(inputBytes))
		n, _, err := udpCon.ReadFromUDP(outputBytes)
		if err != nil {
			addMetric(classifyError(err), 1.0)
			return
		}

		if !byteArrayEquals(outputBytes[:n], inputBytes) {
			addMetric("content_error", 1.0)
			return
		}

		duration := time.Now().Sub(start).Seconds()
		addMetric("ping_success", duration)
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
