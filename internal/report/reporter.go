package report

import (
	"github.com/AleksandrMatsko/yadro-biathlon/internal/config"
	"github.com/AleksandrMatsko/yadro-biathlon/internal/event"
)

// Reporter is used to create Report with aggregated data based on occurred events.
type Reporter struct {
	conf      config.BiathlonCompetition
	reporters map[string]*competitorReporter
}

// NewReporter creates new Reporter.
func NewReporter(conf config.BiathlonCompetition) *Reporter {
	return &Reporter{
		conf:      conf,
		reporters: make(map[string]*competitorReporter),
	}
}

// NotifyWithEvent implements competition.Observer interface, to observe incoming events.
func (r *Reporter) NotifyWithEvent(incomingEvent event.Event) {
	reporter, ok := r.reporters[incomingEvent.CompetitorID]
	if !ok && incomingEvent.ID != event.CompetitorRegistration {
		return
	}

	if incomingEvent.ID == event.CompetitorRegistration {
		reporter = newCompetitorReporter(r.conf)
		r.reporters[incomingEvent.CompetitorID] = reporter
	}

	reporter.NotifyWithEvent(incomingEvent)
}

// MakeReport creates report from previously observed events.
func (r *Reporter) MakeReport() Report {
	records := make([]reportRecord, 0, len(r.reporters))
	for competitorID, reporter := range r.reporters {
		record := reporter.createRecord()
		record.competitorID = competitorID

		records = append(records, record)
	}

	return Report(records)
}

type competitorReporter struct {
	totalTime *totalTimeReporter
	lapsTime  *lapsTimeReporter
	shooting  *shootingReporter
}

func newCompetitorReporter(conf config.BiathlonCompetition) *competitorReporter {
	return &competitorReporter{
		totalTime: newTotalTime(),
		lapsTime:  newLapsTime(conf.Laps, conf.LapLen),
		shooting:  newShootingReporter(conf.FiringLines, conf.PenaltyLen),
	}
}

func (cr *competitorReporter) NotifyWithEvent(e event.Event) {
	cr.totalTime.NotifyWithEvent(e)
	cr.lapsTime.NotifyWithEvent(e)
	cr.shooting.NotifyWithEvent(e)
}

func (cr *competitorReporter) createRecord() reportRecord {
	record := reportRecord{}

	record.totalTime, record.finalState = cr.totalTime.GetTotalTime()
	record.mainLapsInfo = cr.lapsTime.GetLapTimesAndSpeed()
	record.shootingInfo = cr.shooting.GetInfo()

	return record
}
