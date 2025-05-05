package report

import (
	"testing"
	"time"

	"github.com/AleksandrMatsko/yadro-biathlon/internal/event"
	"github.com/stretchr/testify/assert"
)

func Test_shootingReporter(t *testing.T) {
	const (
		firingLines   uint32 = 2
		penaltyLapLen uint32 = 100
	)

	t.Run("when all targets are hit", func(t *testing.T) {
		t.Parallel()

		totalTargets := firingLines * uint32(len(event.AvailableTargets))

		reporter := newShootingReporter(firingLines, penaltyLapLen)

		for range firingLines {
			reporter.NotifyWithEvent(event.Event{ID: event.CompetitorOnFiringRange})
			for target := range event.AvailableTargets {
				reporter.NotifyWithEvent(event.Event{
					ID:    event.TargetHit,
					Extra: target,
				})
			}
			reporter.NotifyWithEvent(event.Event{ID: event.CompetitorLeftFiringRange})
		}

		info := reporter.GetInfo()
		assert.Equal(t, totalTargets, info.TotalTargets)
		assert.Equal(t, totalTargets, info.TotalHitTargets)
		assert.Equal(t, time.Duration(0), info.TimeSpentOnPenaltyLaps)
		assert.Equal(t, 0.0, info.AverageSpeedOnPenaltyLaps)
	})

	t.Run("when not all targets are hit", func(t *testing.T) {
		t.Parallel()

		totalTargets := firingLines * uint32(len(event.AvailableTargets))

		intervals := []time.Duration{
			time.Minute*2 + time.Second + time.Millisecond*535,
			time.Minute*1 + time.Second*55 + time.Millisecond*123,
		}
		timeSpentOnPenaltyLaps := intervals[0] + intervals[1]

		enterFirstPenalty := time.Date(0, time.January, 1, 10, 30, 0, 0, time.UTC)
		leaveFirstPenalty := enterFirstPenalty.Add(intervals[0])

		enterSecondPenalty := time.Date(0, time.January, 1, 11, 45, 0, 0, time.UTC)
		leaveSecondPenalty := enterSecondPenalty.Add(intervals[1])

		reporter := newShootingReporter(firingLines, penaltyLapLen)

		reporter.NotifyWithEvent(event.Event{ID: event.CompetitorOnFiringRange})
		reporter.NotifyWithEvent(event.Event{ID: event.TargetHit, Extra: "1"})
		reporter.NotifyWithEvent(event.Event{ID: event.TargetHit, Extra: "2"})
		reporter.NotifyWithEvent(event.Event{ID: event.TargetHit, Extra: "4"})
		reporter.NotifyWithEvent(event.Event{ID: event.TargetHit, Extra: "5"})
		reporter.NotifyWithEvent(event.Event{ID: event.CompetitorLeftFiringRange})

		reporter.NotifyWithEvent(event.Event{
			Time: enterFirstPenalty,
			ID:   event.CompetitorEnterPenaltyLaps,
		})
		reporter.NotifyWithEvent(event.Event{
			Time: leaveFirstPenalty,
			ID:   event.CompetitorLeftPenaltyLaps,
		})

		reporter.NotifyWithEvent(event.Event{ID: event.CompetitorOnFiringRange})
		reporter.NotifyWithEvent(event.Event{ID: event.TargetHit, Extra: "1"})
		reporter.NotifyWithEvent(event.Event{ID: event.TargetHit, Extra: "2"})
		reporter.NotifyWithEvent(event.Event{ID: event.TargetHit, Extra: "3"})
		reporter.NotifyWithEvent(event.Event{ID: event.CompetitorLeftFiringRange})

		reporter.NotifyWithEvent(event.Event{
			Time: enterSecondPenalty,
			ID:   event.CompetitorEnterPenaltyLaps,
		})
		reporter.NotifyWithEvent(event.Event{
			Time: leaveSecondPenalty,
			ID:   event.CompetitorLeftPenaltyLaps,
		})

		info := reporter.GetInfo()
		assert.Equal(t, totalTargets, info.TotalTargets)
		assert.Equal(t, totalTargets-3, info.TotalHitTargets)
		assert.Equal(t, timeSpentOnPenaltyLaps, info.TimeSpentOnPenaltyLaps)
		assert.Equal(t, float64(3*penaltyLapLen)/timeSpentOnPenaltyLaps.Seconds(), info.AverageSpeedOnPenaltyLaps)
	})

	t.Run("when competitor can't continue", func(t *testing.T) {
		t.Run("and was running main lap without penalty laps", func(t *testing.T) {
			t.Parallel()

			totalTargets := firingLines * uint32(len(event.AvailableTargets))

			reporter := newShootingReporter(firingLines, penaltyLapLen)

			reporter.NotifyWithEvent(event.Event{ID: event.CompetitorOnFiringRange})
			reporter.NotifyWithEvent(event.Event{ID: event.TargetHit, Extra: "1"})
			reporter.NotifyWithEvent(event.Event{ID: event.TargetHit, Extra: "4"})
			reporter.NotifyWithEvent(event.Event{ID: event.TargetHit, Extra: "5"})
			reporter.NotifyWithEvent(event.Event{ID: event.CompetitorLeftFiringRange})

			interval := time.Minute*2 + time.Second + time.Millisecond*535
			enterFirstPenalty := time.Date(0, time.January, 1, 10, 30, 0, 0, time.UTC)
			leaveFirstPenalty := enterFirstPenalty.Add(interval)

			reporter.NotifyWithEvent(event.Event{
				Time: enterFirstPenalty,
				ID:   event.CompetitorEnterPenaltyLaps,
			})
			reporter.NotifyWithEvent(event.Event{
				Time: leaveFirstPenalty,
				ID:   event.CompetitorLeftPenaltyLaps,
			})

			reporter.NotifyWithEvent(event.Event{ID: event.CompetitorCannotContinue})

			reporter.NotifyWithEvent(event.Event{ID: event.CompetitorOnFiringRange})
			for target := range event.AvailableTargets {
				reporter.NotifyWithEvent(event.Event{
					ID:    event.TargetHit,
					Extra: target,
				})
			}
			reporter.NotifyWithEvent(event.Event{ID: event.CompetitorLeftFiringRange})

			info := reporter.GetInfo()
			assert.Equal(t, totalTargets, info.TotalTargets)
			assert.Equal(t, uint32(3), info.TotalHitTargets)
			assert.Equal(t, interval, info.TimeSpentOnPenaltyLaps)
			assert.Equal(t, float64(2*penaltyLapLen)/interval.Seconds(), info.AverageSpeedOnPenaltyLaps)
		})

		t.Run("and was shooting", func(t *testing.T) {
			t.Parallel()

			totalTargets := firingLines * uint32(len(event.AvailableTargets))

			reporter := newShootingReporter(firingLines, penaltyLapLen)

			reporter.NotifyWithEvent(event.Event{ID: event.CompetitorOnFiringRange})
			reporter.NotifyWithEvent(event.Event{ID: event.TargetHit, Extra: "1"})
			reporter.NotifyWithEvent(event.Event{ID: event.TargetHit, Extra: "4"})
			reporter.NotifyWithEvent(event.Event{ID: event.TargetHit, Extra: "5"})
			reporter.NotifyWithEvent(event.Event{ID: event.CompetitorLeftFiringRange})

			interval := time.Minute*2 + time.Second + time.Millisecond*535
			enterFirstPenalty := time.Date(0, time.January, 1, 10, 30, 0, 0, time.UTC)
			leaveFirstPenalty := enterFirstPenalty.Add(interval)

			reporter.NotifyWithEvent(event.Event{
				Time: enterFirstPenalty,
				ID:   event.CompetitorEnterPenaltyLaps,
			})
			reporter.NotifyWithEvent(event.Event{
				Time: leaveFirstPenalty,
				ID:   event.CompetitorLeftPenaltyLaps,
			})

			reporter.NotifyWithEvent(event.Event{ID: event.CompetitorOnFiringRange})
			reporter.NotifyWithEvent(event.Event{ID: event.TargetHit, Extra: "1"})
			reporter.NotifyWithEvent(event.Event{ID: event.TargetHit, Extra: "3"})

			reporter.NotifyWithEvent(event.Event{ID: event.CompetitorCannotContinue})

			reporter.NotifyWithEvent(event.Event{ID: event.TargetHit, Extra: "4"})
			reporter.NotifyWithEvent(event.Event{ID: event.TargetHit, Extra: "5"})
			reporter.NotifyWithEvent(event.Event{ID: event.CompetitorLeftFiringRange})

			info := reporter.GetInfo()
			assert.Equal(t, totalTargets, info.TotalTargets)
			assert.Equal(t, uint32(5), info.TotalHitTargets)
			assert.Equal(t, interval, info.TimeSpentOnPenaltyLaps)
			assert.Equal(t, float64(2*penaltyLapLen)/interval.Seconds(), info.AverageSpeedOnPenaltyLaps)
		})

		t.Run("and was running main lap and have penalty laps to complete", func(t *testing.T) {
			t.Parallel()

			totalTargets := firingLines * uint32(len(event.AvailableTargets))

			reporter := newShootingReporter(firingLines, penaltyLapLen)

			reporter.NotifyWithEvent(event.Event{ID: event.CompetitorOnFiringRange})
			reporter.NotifyWithEvent(event.Event{ID: event.TargetHit, Extra: "1"})
			reporter.NotifyWithEvent(event.Event{ID: event.TargetHit, Extra: "4"})
			reporter.NotifyWithEvent(event.Event{ID: event.TargetHit, Extra: "5"})
			reporter.NotifyWithEvent(event.Event{ID: event.CompetitorLeftFiringRange})

			reporter.NotifyWithEvent(event.Event{ID: event.CompetitorCannotContinue})

			interval := time.Minute*2 + time.Second + time.Millisecond*535
			enterFirstPenalty := time.Date(0, time.January, 1, 10, 30, 0, 0, time.UTC)
			leaveFirstPenalty := enterFirstPenalty.Add(interval)

			reporter.NotifyWithEvent(event.Event{
				Time: enterFirstPenalty,
				ID:   event.CompetitorEnterPenaltyLaps,
			})
			reporter.NotifyWithEvent(event.Event{
				Time: leaveFirstPenalty,
				ID:   event.CompetitorLeftPenaltyLaps,
			})

			info := reporter.GetInfo()
			assert.Equal(t, totalTargets, info.TotalTargets)
			assert.Equal(t, uint32(3), info.TotalHitTargets)
			assert.Equal(t, time.Duration(0), info.TimeSpentOnPenaltyLaps)
			assert.Equal(t, 0.0, info.AverageSpeedOnPenaltyLaps)
		})

		t.Run("and was running penalty lap", func(t *testing.T) {
			t.Parallel()

			totalTargets := firingLines * uint32(len(event.AvailableTargets))

			reporter := newShootingReporter(firingLines, penaltyLapLen)

			reporter.NotifyWithEvent(event.Event{ID: event.CompetitorOnFiringRange})
			reporter.NotifyWithEvent(event.Event{ID: event.TargetHit, Extra: "1"})
			reporter.NotifyWithEvent(event.Event{ID: event.TargetHit, Extra: "4"})
			reporter.NotifyWithEvent(event.Event{ID: event.TargetHit, Extra: "5"})
			reporter.NotifyWithEvent(event.Event{ID: event.CompetitorLeftFiringRange})

			interval := time.Minute*2 + time.Second + time.Millisecond*535
			enterFirstPenalty := time.Date(0, time.January, 1, 10, 30, 0, 0, time.UTC)
			leaveFirstPenalty := enterFirstPenalty.Add(interval)

			reporter.NotifyWithEvent(event.Event{
				Time: enterFirstPenalty,
				ID:   event.CompetitorEnterPenaltyLaps,
			})

			reporter.NotifyWithEvent(event.Event{ID: event.CompetitorCannotContinue})

			reporter.NotifyWithEvent(event.Event{
				Time: leaveFirstPenalty,
				ID:   event.CompetitorLeftPenaltyLaps,
			})

			info := reporter.GetInfo()
			assert.Equal(t, totalTargets, info.TotalTargets)
			assert.Equal(t, uint32(3), info.TotalHitTargets)
			assert.Equal(t, time.Duration(0), info.TimeSpentOnPenaltyLaps)
			assert.Equal(t, 0.0, info.AverageSpeedOnPenaltyLaps)
		})
	})
}
