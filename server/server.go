package server

import "fmt"

type Config struct {
	TcpPort int `yaml:"tcp_port"`
	UdpPort int `yaml:"udp_port"`
}

func Run(config Config) error {
	fmt.Println("Running server", config)
	return nil
}
