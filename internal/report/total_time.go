package report

import (
	"time"

	"github.com/AleksandrMatsko/yadro-biathlon/internal/event"
	"github.com/AleksandrMatsko/yadro-biathlon/internal/event/parser"
)

type totalTimeState string

const (
	initial     totalTimeState = "Initial"
	running     totalTimeState = "Running"
	notStarted  totalTimeState = "NotStarted"
	notFinished totalTimeState = "NotFinished"
	finished    totalTimeState = "Finished"
)

type totalTime struct {
	state totalTimeState
	start time.Time
	end   time.Time
}

func newTotalTime() *totalTime {
	return &totalTime{
		state: initial,
	}
}

func (tt *totalTime) NotifyWithEvent(e event.Event) {
	if tt.state == initial && e.ID == event.StartTimeAssignment {
		tt.start, _ = parser.ParseTime(e.Extra)
		return
	}

	if tt.state == initial && e.ID == event.CompetitorDisqualified {
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
}

func (tt *totalTime) GetTotalTime() (time.Duration, totalTimeState) {
	if tt.state == finished {
		return tt.end.Sub(tt.start), finished
	}

	return 0, tt.state
}
