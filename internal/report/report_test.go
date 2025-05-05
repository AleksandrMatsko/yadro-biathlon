package report

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestReport_String(t *testing.T) {
	givenReport := Report([]reportRecord{
		{
			totalTime:    time.Minute*1 + time.Second*2 + time.Millisecond*345,
			finalState:   finished,
			competitorID: "1",
			mainLapsInfo: []mainLapInfo{
				{
					Interval: time.Minute*29 + time.Second*14 + time.Millisecond*7,
					Speed:    2.04,
				},
				{
					Interval: time.Minute*31 + time.Second*2 + time.Millisecond*56,
					Speed:    1.92,
				},
			},
			shootingInfo: shootingInfo{
				TotalTargets:              10,
				TotalHitTargets:           8,
				TimeSpentOnPenaltyLaps:    time.Minute*2 + time.Second*32,
				AverageSpeedOnPenaltyLaps: 1.87,
			},
		},
		{
			finalState:   notStarted,
			competitorID: "2",
			mainLapsInfo: make([]mainLapInfo, 2),
			shootingInfo: shootingInfo{
				TotalTargets: 10,
			},
		},
		{
			finalState:   notFinished,
			competitorID: "3",
			mainLapsInfo: []mainLapInfo{
				{
					Interval: time.Minute*28 + time.Second*53 + time.Millisecond*497,
					Speed:    2.04,
				},
				{},
			},
			shootingInfo: shootingInfo{
				TotalTargets:              10,
				TotalHitTargets:           4,
				TimeSpentOnPenaltyLaps:    time.Minute*1 + time.Second*10,
				AverageSpeedOnPenaltyLaps: 2.798,
			},
		},
	})

	expectedString := strings.Join([]string{
		"[00:01:02.345] 1 [{00:29:14.007, 2.040}, {00:31:02.056, 1.920}] {00:02:32.000, 1.870} 8/10",
		"[NotStarted] 2 [{,}, {,}] {00:00:00.000, 0.000} 0/10",
		"[NotFinished] 3 [{00:28:53.497, 2.040}, {,}] {00:01:10.000, 2.798} 4/10",
	}, "\n") + "\n"

	assert.Equal(t, expectedString, givenReport.String())
}
