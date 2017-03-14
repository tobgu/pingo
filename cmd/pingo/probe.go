package main

import (
	"flag"
	"github.com/pkg/errors"
	"github.com/tobgu/pingo/probe"
)

func runProbe(args []string) error {
	flagset := flag.NewFlagSet("probe", flag.ExitOnError)
	var (
		configFile = flagset.String("config", "config.yaml", "Configuration file")
	)

	flagset.Usage = usageFor(flagset, "pingo probe [flags]")
	if err := flagset.Parse(args); err != nil {
		return errors.Wrap(err, "Failed parsing command line options")
	}

	var config probe.Config
	if err := readYaml(*configFile, &config); err != nil {
		return err
	}

	// This will block forever unless an error occurs
	return probe.Run(config)
}
