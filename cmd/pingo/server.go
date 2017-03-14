package main

import (
	"flag"
	"github.com/pkg/errors"
	"github.com/tobgu/pingo/server"
)

func runServer(args []string) error {
	flagset := flag.NewFlagSet("server", flag.ExitOnError)
	var (
		configFile = flagset.String("config", "config.yaml", "Configuration file")
	)

	flagset.Usage = usageFor(flagset, "pingo server [flags]")
	if err := flagset.Parse(args); err != nil {
		return errors.Wrap(err, "Failed parsing command line options")
	}

	var config server.Config
	if err := readYaml(*configFile, &config); err != nil {
		return err
	}

	return server.Run(config)
}
