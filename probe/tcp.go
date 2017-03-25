package probe

import (
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

type tcpProbe struct {
	inputBytes  []byte
	outputBytes []byte
	config      Config
	host        Host
	statistics  *Statistics
}

func newTcp(config Config, host Host, statistics *Statistics) *tcpProbe {
	return &tcpProbe{
		inputBytes:  randomBytes(config.TcpSize),
		outputBytes: make([]byte, config.TcpSize),
		host: host,
		statistics: statistics,
		config: config,
	}
}

func (p *tcpProbe) addMetric(kind string, value float64) {
	p.statistics.Add(p.host.Name, "tcp", kind, value)
}

func (p *tcpProbe) execute() {
	start := time.Now()
	con, err := net.DialTimeout("tcp",
		fmt.Sprintf("%s:%d", p.host.Address, p.host.TcpPort),
		time.Duration(p.config.ConnectionTimeout)*time.Second)
	if err != nil {
		p.addMetric(classifyError(err), 1.0)
		return
	}

	defer con.Close()

	t0 := time.Now()
	err = con.SetDeadline(time.Now().Add(time.Duration(p.config.ReadTimeout) * time.Second))
	if err != nil {
		p.addMetric(classifyError(err), 1.0)
		return
	}

	_, err = con.Write(p.inputBytes)
	if err != nil {
		log.Println("Write error after", time.Now().Sub(t0))
		p.addMetric(classifyError(err), 1.0)
		return
	}

	if tcpcon, ok := con.(*net.TCPConn); ok {
		tcpcon.CloseWrite()
	} else {
		log.Println("Connection was not of type TCP, this was unexpected...")
		return
	}

	_, err = io.ReadFull(con, p.outputBytes)
	if err != nil {
		log.Println("Read error after", time.Now().Sub(t0))
		p.addMetric(classifyError(err), 1.0)
		return
	}

	if !byteArrayEquals(p.outputBytes, p.inputBytes) {
		log.Println("Content error")
		p.addMetric("content_error", 1.0)
		return
	}

	duration := time.Now().Sub(start).Seconds()
	p.addMetric("ping_success", duration)
}
