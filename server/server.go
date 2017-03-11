package server

import "fmt"

type Config struct {
	TcpPort int `yaml:"tcp_port"`
	UdpPort int `yaml:"udp_port"`
}

func RunServer() {
	fmt.Println("Running server")
}
