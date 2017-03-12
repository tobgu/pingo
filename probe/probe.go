package probe

import (
	"fmt"
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

func startProbe(config Config, host Host, statistics *Statistics) error {
	return nil
}

func Run(config Config) error {
	s := NewStatistics(config.StatisticsRetentionPeriod)
	for _, h := range config.Hosts {
		startProbe(config, h, s)
	}

	fmt.Println("Running probe", config)
	return nil
}
