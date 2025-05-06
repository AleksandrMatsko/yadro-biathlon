package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/AleksandrMatsko/yadro-biathlon/internal/competition"
	"github.com/AleksandrMatsko/yadro-biathlon/internal/config"
	"github.com/AleksandrMatsko/yadro-biathlon/internal/event"
	"github.com/AleksandrMatsko/yadro-biathlon/internal/event/parser"
	"github.com/AleksandrMatsko/yadro-biathlon/internal/report"
)

var (
	configFilePathFlag     = flag.String("config", "", "Path to configuration file")
	printConfigFlag        = flag.Bool("print-config", false, "Print current config to stdout")
	incomingEventsFilePath = flag.String("events", "", "Path to events file")
	reportFilePathFlag     = flag.String("report", "report.txt", "Path of file to save report")
)

var (
	errNoConfigFile = errors.New("no config file, provide it with --config option")
	errNoEventsFile = errors.New("no events file, provide it with --events option")
)

func main() {
	flag.Parse()

	if err := makeReport(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}

func makeReport() error {
	if *configFilePathFlag == "" {
		return errNoConfigFile
	}

	conf := config.BiathlonCompetition{}
	err := config.Read(*configFilePathFlag, &conf)
	if err != nil {
		return fmt.Errorf("read config: %w", err)
	}

	if *printConfigFlag {
		config.Print(conf)
	}

	if *incomingEventsFilePath == "" {
		return errNoEventsFile
	}

	reporter := report.NewReporter(conf)

	biathlon, err := competition.NewBiathlon(
		conf,
		competition.NewComposedObserver().AddObservers(
			event.NewLogger(os.Stdout),
			reporter,
		),
	)
	if err != nil {
		return fmt.Errorf("failed to create biathlon competition: %w", err)
	}

	file, err := os.Open(*incomingEventsFilePath)
	if err != nil {
		return fmt.Errorf("open events file: %w", err)
	}
	defer file.Close()

	lines, retErrFunc := parser.Lines(file)
	for event, err := range parser.ParsedLines(lines) {
		if err != nil {
			return fmt.Errorf("parsing file: %w", err)
		}

		biathlon.HandleEvent(event)
	}

	err = retErrFunc()
	if err != nil {
		return fmt.Errorf("reading file: %w", err)
	}

	report := reporter.MakeReport()
	report.Sort()

	reportFile, err := os.OpenFile(*reportFilePathFlag, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0o666)
	if err != nil {
		return fmt.Errorf("open file for report: '%s': %w", *reportFilePathFlag, err)
	}
	defer reportFile.Close()

	fmt.Fprint(reportFile, report)

	return nil
}
