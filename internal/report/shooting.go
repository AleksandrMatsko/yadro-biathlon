package report

import (
	"fmt"
	"time"

	"github.com/AleksandrMatsko/yadro-biathlon/internal/event"
)

type shootingReporterCompetitorState uint8

const (
	runningMainLap shootingReporterCompetitorState = iota
	shooting
	runningPenaltyLaps
	ended
)

// shootingReporter is responsible for calculating several values:
//   - total number of hit targets;
//   - total time spent on penalty laps;
//   - average speed on penalty laps.
//
// Penalty laps calculations is here because, penalty laps count depends on
// the amount of not hit targets on the firing range.
type shootingReporter struct {
	hitTargetsOnCurrentFireRange map[string]struct{}
	firingLinesCount             uint32
	completedShootings           uint32
	totalNumberOfHitTarges       uint32

	totalPenaltyLapCount             uint32
	penaltyLapToPerformAfterShooting uint8
	penaltyLapLen                    uint32

	timeSpentOnPenaltyLaps time.Duration
	enterPenaltyLap        time.Time

	state shootingReporterCompetitorState
}

func newShootingReporter(firingLinesCount uint32, penaltyLapLen uint32) *shootingReporter {
	return &shootingReporter{
		hitTargetsOnCurrentFireRange: make(map[string]struct{}, len(event.AvailableTargets)),
		firingLinesCount:             firingLinesCount,
		completedShootings:           0,
		totalNumberOfHitTarges:       0,
		totalPenaltyLapCount:         0,
		penaltyLapLen:                penaltyLapLen,
		state:                        runningMainLap,
	}
}

func (s *shootingReporter) NotifyWithEvent(e event.Event) {
	if s.state == ended {
		return
	}

	if e.ID == event.CompetitorCannotContinue || e.ID == event.CompetitorDisqualified || e.ID == event.CompetitorFinished {
		s.state = ended
		return
	}

	switch s.state {
	case runningMainLap:
		s.onRunningMainLap(e)
	case shooting:
		s.onShooting(e)
	case runningPenaltyLaps:
		s.onRunningPenaltyLaps(e)
	}
}

func (s *shootingReporter) onRunningMainLap(e event.Event) {
	if e.ID == event.CompetitorOnFiringRange && s.completedShootings < s.firingLinesCount {
		clear(s.hitTargetsOnCurrentFireRange)
		s.state = shooting
		s.penaltyLapToPerformAfterShooting = uint8(len(event.AvailableTargets))
		return
	}

	if e.ID == event.CompetitorEnterPenaltyLaps {
		s.enterPenaltyLap = e.Time
		s.state = runningPenaltyLaps
		return
	}
}

func (s *shootingReporter) onShooting(e event.Event) {
	if e.ID == event.TargetHit {
		if _, ok := s.hitTargetsOnCurrentFireRange[e.Extra]; !ok {
			s.hitTargetsOnCurrentFireRange[e.Extra] = struct{}{}
			s.totalNumberOfHitTarges += 1
			s.penaltyLapToPerformAfterShooting -= 1
		}

		return
	}

	if e.ID == event.CompetitorLeftFiringRange {
		s.completedShootings += 1
		s.state = runningMainLap
		return
	}
}

func (s *shootingReporter) onRunningPenaltyLaps(e event.Event) {
	if e.ID == event.CompetitorLeftPenaltyLaps {
		s.totalPenaltyLapCount += uint32(s.penaltyLapToPerformAfterShooting)
		s.penaltyLapToPerformAfterShooting = 0
		s.timeSpentOnPenaltyLaps += e.Time.Sub(s.enterPenaltyLap)

		s.state = runningMainLap
		return
	}
}

type shootingInfo struct {
	TotalTargets              uint32
	TotalHitTargets           uint32
	TimeSpentOnPenaltyLaps    time.Duration
	AverageSpeedOnPenaltyLaps float64
}

func (info shootingInfo) String() string {
	return fmt.Sprintf("{%s, %.3f} %d/%d",
		formatDuration(info.TimeSpentOnPenaltyLaps),
		info.AverageSpeedOnPenaltyLaps,
		info.TotalHitTargets,
		info.TotalTargets,
	)
}

func (s *shootingReporter) GetInfo() shootingInfo {
	res := shootingInfo{}

	res.TotalTargets = uint32(len(event.AvailableTargets) * int(s.firingLinesCount))
	res.TotalHitTargets = s.totalNumberOfHitTarges
	res.TimeSpentOnPenaltyLaps = s.timeSpentOnPenaltyLaps
	if s.timeSpentOnPenaltyLaps != 0 {
		res.AverageSpeedOnPenaltyLaps = float64(s.totalPenaltyLapCount*s.penaltyLapLen) / s.timeSpentOnPenaltyLaps.Seconds()
	}

	return res
}
