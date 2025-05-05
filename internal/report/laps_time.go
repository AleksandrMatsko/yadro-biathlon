package report

import (
	"fmt"
	"time"

	"github.com/AleksandrMatsko/yadro-biathlon/internal/event"
	"github.com/AleksandrMatsko/yadro-biathlon/internal/event/parser"
)

type lapsTimeReporter struct {
	lapsCompleted uint32
	lapStart      time.Time
	lapLen        uint32
	lapTimes      []time.Duration
	stop          bool
}

func newLapsTimeReporter(laps uint32, lapLen uint32) *lapsTimeReporter {
	return &lapsTimeReporter{
		lapsCompleted: 0,
		lapLen:        lapLen,
		lapTimes:      make([]time.Duration, laps),
	}
}

func (lt *lapsTimeReporter) NotifyWithEvent(e event.Event) {
	if lt.lapsCompleted == uint32(len(lt.lapTimes)) {
		return
	}

	if lt.lapsCompleted == 0 && e.ID == event.StartTimeAssignment {
		lt.lapStart, _ = parser.ParseTime(e.Extra)
		return
	}

	if e.ID == event.CompetitorEndedMainLap {
		lapTime := e.Time.Sub(lt.lapStart)
		lt.lapTimes[lt.lapsCompleted] = lapTime
		lt.lapsCompleted += 1

		lt.lapStart = e.Time

		return
	}

	if e.ID == event.CompetitorDisqualified || e.ID == event.CompetitorCannotContinue {
		lt.stop = true
	}
}

type mainLapInfo struct {
	Interval time.Duration
	Speed    float64
}

func (info mainLapInfo) String() string {
	if info.Interval == 0 && info.Speed == 0 {
		return "{,}"
	}

	return fmt.Sprintf("{%s, %.3f}", formatDuration(info.Interval), info.Speed)
}

func (lt *lapsTimeReporter) GetLapTimesAndSpeed() []mainLapInfo {
	result := make([]mainLapInfo, len(lt.lapTimes))
	for i, interval := range lt.lapTimes {
		if interval == 0 {
			continue
		}

		result[i].Interval = interval
		result[i].Speed = float64(lt.lapLen) / interval.Seconds()
	}

	return result
}
