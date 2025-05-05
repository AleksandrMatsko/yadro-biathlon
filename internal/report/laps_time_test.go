package report

import (
	"testing"
	"time"

	"github.com/AleksandrMatsko/yadro-biathlon/internal/event"
	"github.com/stretchr/testify/assert"
)

func Test_lapsTime(t *testing.T) {
	const (
		lapsCount uint32 = 3
		lapLen    uint32 = 1000
	)

	intervals := []time.Duration{
		time.Minute*10 + time.Second*32 + time.Millisecond*54,
		time.Minute*11 + time.Second*2 + time.Millisecond*432,
		time.Minute*10 + time.Second*54 + time.Millisecond*24,
	}

	startTime := time.Date(0, time.January, 1, 10, 30, 0, 0, time.UTC)
	lapEnds := []time.Time{
		startTime.Add(intervals[0]),
		startTime.Add(intervals[0] + intervals[1]),
		startTime.Add(intervals[0] + intervals[1] + intervals[2]),
	}

	t.Run("with all laps completed", func(t *testing.T) {
		t.Parallel()

		lt := newLapsTime(lapsCount-1, lapLen)

		lt.NotifyWithEvent(event.Event{
			ID:    event.StartTimeAssignment,
			Extra: lapEnds[0].Format(event.TimeFormat),
		})
		for _, lapEnd := range lapEnds[1:] {
			lt.NotifyWithEvent(event.Event{
				Time: lapEnd,
				ID:   event.CompetitorEndedMainLap,
			})
		}

		expected := make([]mainLapInfo, 0, len(intervals[1:]))
		for _, interval := range intervals[1:] {
			expected = append(expected, mainLapInfo{
				Interval: interval,
				Speed:    float64(lapLen) / interval.Seconds(),
			})
		}

		got := lt.GetLapTimesAndSpeed()

		assert.Equal(t, expected, got)
	})

	t.Run("with not all laps completed", func(t *testing.T) {
		t.Parallel()

		completedLaps := 1

		lt := newLapsTime(lapsCount, lapLen)

		lt.NotifyWithEvent(event.Event{
			ID:    event.StartTimeAssignment,
			Extra: startTime.Format(event.TimeFormat),
		})
		for i := range completedLaps {
			lt.NotifyWithEvent(event.Event{
				Time: lapEnds[i],
				ID:   event.CompetitorEndedMainLap,
			})
		}

		expected := make([]mainLapInfo, len(intervals))
		for i := range completedLaps {
			expected[i] = mainLapInfo{
				Interval: intervals[i],
				Speed:    float64(lapLen) / intervals[i].Seconds(),
			}
		}

		got := lt.GetLapTimesAndSpeed()

		assert.Equal(t, expected, got)
	})

	t.Run("with more laps completed", func(t *testing.T) {
		t.Parallel()

		lt := newLapsTime(lapsCount-1, lapLen)

		lt.NotifyWithEvent(event.Event{
			ID:    event.StartTimeAssignment,
			Extra: startTime.Format(event.TimeFormat),
		})
		for _, lapEnd := range lapEnds {
			lt.NotifyWithEvent(event.Event{
				Time: lapEnd,
				ID:   event.CompetitorEndedMainLap,
			})
		}

		expected := make([]mainLapInfo, lapsCount-1)
		for i := range expected {
			expected[i] = mainLapInfo{
				Interval: intervals[i],
				Speed:    float64(lapLen) / intervals[i].Seconds(),
			}
		}

		got := lt.GetLapTimesAndSpeed()

		assert.Equal(t, expected, got)
	})
}
