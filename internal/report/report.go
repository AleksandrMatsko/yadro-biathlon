package report

import (
	"fmt"
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
