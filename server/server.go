package server

import (
	"fmt"
	"net"
	"github.com/pkg/errors"
	"log"
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

type Config struct {
	TcpPort int `yaml:"tcp_port"`
	UdpPort int `yaml:"udp_port"`
}

func Run(config Config) error {
	if config.UdpPort != 0 {
		if err := startUdpServer(config.UdpPort); err != nil {
			return err
		}
	}

	// Block forever
	select {}
	return nil
}
