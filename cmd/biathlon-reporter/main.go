package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/AleksandrMatsko/yadro-biathlon/internal/config"
)

var (
	configFileNameFlag = flag.String("config", "", "Path to configuration file")
	printConfigFlag    = flag.Bool("print-config", false, "Print current config to stdout")
)

func main() {
	flag.Parse()

	if configFileNameFlag == nil || *configFileNameFlag == "" {
		fmt.Fprint(os.Stderr, "Config file required, provide it with --config option\n")
		os.Exit(1)
	}

	conf := config.BiathlonCompetition{}
	err := config.Read(*configFileNameFlag, &conf)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}

	if *printConfigFlag {
		config.Print(conf)
	}
}
