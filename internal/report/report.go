package report

import (
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/AleksandrMatsko/yadro-biathlon/internal/event"
)

// Report is list of aggregated data by competitor.
type Report []reportRecord

// String formats Report and return it as string.
// Each line in the returned string is related to unique competitor.
// Examples of formatter lines (if there were 2 laps in the race and 2 firing ranges):
//
//	[00:44:51.123] 1 [{00:22:20.100, 2.2}, {00:22:31.023, 1.9}] {00:01:25.467, 0.496} 8/10
//	[NotStarted] 2 [{,}, {,}] {00:00:00.000, 0.000} 0/10
//	[NotFinished] 3 [{00:29:03.872, 2.093}, {,}] {00:01:44.296, 0.481} 4/10
func (report Report) String() string {
	builder := strings.Builder{}

	for _, record := range report {
		_, _ = builder.WriteString(record.String())
		_, _ = builder.WriteString("\n")
	}

	return builder.String()
}

// Sort Report records by competitors' total time.
// Sorted records have the following order:
//   - [NotStarted] competitors sorted by competitorID.
//   - competirors sorted by total time.
//   - [NotFinished] competitors sorted by competitorID.
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
			res := int(first.totalTime - second.totalTime)
			if res != 0 {
				return res
			}
		}

		return strings.Compare(first.competitorID, second.competitorID)
	})
}

type reportRecord struct {
	totalTime    time.Duration
	finalState   totalTimeState
	competitorID string
	mainLapsInfo []mainLapInfo
	shootingInfo shootingInfo
}

func (rr reportRecord) String() string {
	totalTimeValue := string(rr.finalState)
	if rr.finalState == finished {
		totalTimeValue = formatDuration(rr.totalTime)
	}

	mainLapInfoStrings := make([]string, 0, len(rr.mainLapsInfo))
	for _, info := range rr.mainLapsInfo {
		mainLapInfoStrings = append(mainLapInfoStrings, info.String())
	}

	return fmt.Sprintf("[%s] %s [%s] %s",
		totalTimeValue,
		rr.competitorID,
		strings.Join(mainLapInfoStrings, ", "),
		rr.shootingInfo,
	)
}

func formatDuration(d time.Duration) string {
	timeVal := time.Date(0, time.January, 1, 0, 0, 0, 0, time.UTC).Add(d)
	return timeVal.Format(event.TimeFormat)
}
