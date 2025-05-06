package report

import (
	"fmt"
	"time"

	"github.com/AleksandrMatsko/yadro-biathlon/internal/event"
	"github.com/AleksandrMatsko/yadro-biathlon/internal/event/parser"
)

// lapsTimeReporter is responsible for calculating time and speed for each main lap.
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
	if lt.stop {
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

		if lt.lapsCompleted == uint32(len(lt.lapTimes)) {
			lt.stop = true
		}

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

// GetLapTimesAndSpeed returns time spent and speed for each lap.
//
// Note that time for the first lap is time interval between scheduled start
// and end of the first lap. For other laps lap time is time interval between
// end of previous lap and end of current lap.
//
// Also note that lap time includes time spent on penalty laps.
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
