package report

import (
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/AleksandrMatsko/yadro-biathlon/internal/event"
)

type Report []reportRecord

func (report Report) String() string {
	builder := strings.Builder{}

	for _, record := range report {
		_, _ = builder.WriteString(record.String())
		_, _ = builder.WriteString("\n")
	}

	return builder.String()
}

func (report Report) Sort() {
	slices.SortFunc(report, func(first, second reportRecord) int {
		if first.finalState != second.finalState {
			if first.finalState == notStarted {
				return -1
			}

			if second.finalState == notStarted {
				return 1
			}

			if first.finalState == notFinished {
				return 1
			}

			if second.finalState == notFinished {
				return -1
			}
		}

		if first.finalState == second.finalState && first.finalState == finished {
			return int(first.totalTime - second.totalTime)
		}

		return strings.Compare(first.competitorID, second.competitorID)
	})
}

type reportRecord struct {
	totalTime    time.Duration
	finalState   totalTimeState
	competitorID string
}

func (rr reportRecord) String() string {
	totalTimeValue := string(rr.finalState)
	if rr.finalState == finished {
		timeVal := time.Date(0, time.January, 1, 0, 0, 0, 0, time.UTC).Add(rr.totalTime)
		totalTimeValue = timeVal.Format(event.TimeFormat)
	}

	return fmt.Sprintf("[%s] %s", totalTimeValue, rr.competitorID)
}
