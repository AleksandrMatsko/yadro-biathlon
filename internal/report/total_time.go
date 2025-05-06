package report

import (
	"time"

	"github.com/AleksandrMatsko/yadro-biathlon/internal/event"
	"github.com/AleksandrMatsko/yadro-biathlon/internal/event/parser"
)

type totalTimeReporterCompetitorState string

const (
	initial     totalTimeReporterCompetitorState = "Initial"
	running     totalTimeReporterCompetitorState = "Running"
	notStarted  totalTimeReporterCompetitorState = "NotStarted"
	notFinished totalTimeReporterCompetitorState = "NotFinished"
	finished    totalTimeReporterCompetitorState = "Finished"
)

// totalTimeReporter calculates total time of single competitor.
type totalTimeReporter struct {
	state totalTimeReporterCompetitorState
	start time.Time
	end   time.Time
}

func newTotalTimeReporter() *totalTimeReporter {
	return &totalTimeReporter{
		state: initial,
	}
}

func (tt *totalTimeReporter) NotifyWithEvent(e event.Event) {
	if tt.state == initial && e.ID == event.StartTimeAssignment {
		tt.start, _ = parser.ParseTime(e.Extra)
		return
	}

	if (tt.state == initial || tt.state == running) && e.ID == event.CompetitorDisqualified {
		tt.state = notStarted
		return
	}

	if tt.state == initial && e.ID == event.CompetitorStarted {
		tt.state = running
		return
	}

	if tt.state == running && e.ID == event.CompetitorFinished {
		tt.end = e.Time
		tt.state = finished
		return
	}

	if tt.state == running && e.ID == event.CompetitorCannotContinue {
		tt.state = notFinished
		return
	}
}

// GetTotalTime returns:
//   - Calculated total time (time interval between scheduled start for competitor and time of completing last lap).
//   - Final state of the competitor (on of: notStarted, notFinished, finished).
func (tt *totalTimeReporter) GetTotalTime() (time.Duration, totalTimeReporterCompetitorState) {
	if tt.state == finished {
		return tt.end.Sub(tt.start), finished
	}

	if tt.state != notStarted && tt.state != notFinished {
		if tt.state == running {
			tt.state = notFinished
		} else {
			tt.state = notStarted
		}
	}

	return 0, tt.state
}
