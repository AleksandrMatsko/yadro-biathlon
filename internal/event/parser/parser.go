// parser for incoming events.
package parser

import (
	"bufio"
	"fmt"
	"io"
	"iter"
	"strconv"
	"strings"
	"time"

	"github.com/AleksandrMatsko/yadro-biathlon/internal/event"
)

// Lines returns sequence of lines from the reader and function
// that returns any occurred error. After iterration over sequence is completed
// call function to check error.
func Lines(r io.Reader) (iter.Seq[string], func() error) {
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanLines)

	return func(yield func(string) bool) {
			for scanner.Scan() {
				if !yield(scanner.Text()) {
					return
				}
			}
		},
		func() error {
			return scanner.Err()
		}
}

// ParsedLines returns sequence of (event, error) got after parsing sequence of lines.
func ParsedLines(lines iter.Seq[string]) iter.Seq2[event.Event, error] {
	return func(yield func(event.Event, error) bool) {
		for line := range lines {
			if !yield(ParseSingleLine(line)) {
				return
			}
		}
	}
}

const maxNumParts = 4

// ParseSingleLine expect line to have such format:
//
//	[HH:MM:SS.sss] eventID competitorID extra
//
// If there are less than 3 arguments or if values are invalid,
// the error is returned.
func ParseSingleLine(line string) (event.Event, error) {
	split := strings.SplitN(line, " ", maxNumParts)

	if len(split) < maxNumParts-1 {
		argsCount := len(split)
		if line == "" {
			argsCount = 0
		}

		return event.Event{}, fmt.Errorf("not enough arguments, expected at least 3, got %d", argsCount)
	}

	parsedTime, err := ParseTime(strings.Trim(split[0], "[]"))
	if err != nil {
		return event.Event{}, fmt.Errorf("failed to parse event timestamp: %w", err)
	}

	eventID, err := ParseEventID(split[1])
	if err != nil {
		return event.Event{}, fmt.Errorf("failed to parse event id: %w", err)
	}

	extra := ""
	if len(split) == maxNumParts {
		extra = split[3]
	}

	return event.Event{
		Time:         parsedTime,
		ID:           eventID,
		CompetitorID: split[2],
		Extra:        extra,
	}, nil
}

// ParseTime from given string into time.Time according to event.TimeFormat.
func ParseTime(s string) (time.Time, error) {
	return time.Parse(event.TimeFormat, s)
}

// ParseEventID from given string into event ID.
func ParseEventID(s string) (event.EventID, error) {
	val, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse value: %w", err)
	}

	converted := uint8(val)
	if !event.ValidIncomingEventID(converted) {
		return 0, fmt.Errorf("invalid event id: %s", s)
	}

	return event.EventID(converted), nil
}
