package server

import (
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"log"
	"net"
	"time"
)

func echoUdp(conn *net.UDPConn) {
	defer conn.Close()
	buf := make([]byte, 2048)
	for {
		n, addr, err := conn.ReadFromUDP(buf)
		if err != nil {
			log.Println("Error reading UDP datagram: ", err)
		} else {
			conn.WriteTo(buf[:n], addr)
		}
	}

}

func startUdpServer(port int) error {
	log.Println("Starting UDP server on port", port)
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", port))
	if err != nil {
		return errors.Wrap(err, "Error binding UDP port")
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return errors.Wrap(err, "Error listenig to UDP port")
	}

	go echoUdp(conn)
	return nil
}

func runTcpServer(port, timeout int) error {
	log.Println("Starting TCP server on port", port)
	server, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return errors.Wrap(err, "Error listening to TCP port")
	}

	for {
		client, err := server.Accept()
		if err != nil {
			log.Println("Error accepting TCP connection: ", err)
		}

		if timeout > 0 {
			client.SetDeadline(time.Now().Add(time.Duration(timeout) * time.Second))
		}
		go echoTcp(client)
	}
}

func echoTcp(client net.Conn) {
	defer client.Close()

	// Since the client will first submit all data and then read it back
	// we first read all data then submit the answer instead of streaming
	// it since that risk filling up TCP buffers in the client for large
	// payloads. We also want to measure the true roundtrip time without
	// pipelining of data transmission back and forth.
	buf, err := ioutil.ReadAll(client)
	if err != nil {
		log.Println("Error reading TCP data", len(buf), err)
		return
	}

	n, err := client.Write(buf)
	if err != nil {
		log.Println("Error writing TCP data", n, err)
		return
	}

}

type Config struct {
	TcpPort    int `yaml:"tcp_port"`
	TcpTimeout int `yaml:"tcp_timeout"`
	UdpPort    int `yaml:"udp_port"`
}

func Run(config Config) error {
	if config.UdpPort != 0 {
		if err := startUdpServer(config.UdpPort); err != nil {
			return err
		}
	}

	if config.TcpPort != 0 {
		if err := runTcpServer(config.TcpPort, config.TcpTimeout); err != nil {
			return err
		}
	}

	// Block forever
	select {}
	return nil
}
