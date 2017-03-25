package probe

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

type udpProbe struct {
	probe
	inputBytes []byte
}

func newUdp(config Config, host Host, statistics *Statistics) *udpProbe {
	return &udpProbe{
		probe: probe{config: config,
			host:       host,
			statistics: statistics},
		inputBytes: randomBytes(config.UdpSize),
	}
}

func (p *udpProbe) addMetric(kind string, value float64) {
	p.statistics.Add(p.host.Name, "udp", kind, value)
}

func (p *udpProbe) execute() {
	start := time.Now()
	con, err := net.DialTimeout("udp",
		fmt.Sprintf("%s:%d", p.host.Address, p.host.UdpPort),
		time.Duration(p.config.ConnectionTimeout)*time.Second)
	if err != nil {
		p.addMetric(classifyError(err), 1.0)
		return
	}

	defer con.Close()

	err = con.SetDeadline(time.Now().Add(time.Duration(p.config.ReadTimeout) * time.Second))
	if err != nil {
		p.addMetric(classifyError(err), 1.0)
		return
	}

	inputBuffer := bytes.NewBuffer(p.inputBytes)
	_, err = io.Copy(con, inputBuffer)
	if err != nil {
		p.addMetric(classifyError(err), 1.0)
		return
	}

	udpCon, ok := con.(*net.UDPConn)
	if !ok {
		log.Println("Connection was not of type UDP, this was unexpected...")
		return
	}

	outputBytes := make([]byte, len(p.inputBytes))
	n, _, err := udpCon.ReadFromUDP(outputBytes)
	if err != nil {
		p.addMetric(classifyError(err), 1.0)
		return
	}

	if !byteArrayEquals(outputBytes[:n], p.inputBytes) {
		p.addMetric("content_error", 1.0)
		return
	}

	duration := time.Now().Sub(start).Seconds()
	p.addMetric("ping_success", duration)
}
