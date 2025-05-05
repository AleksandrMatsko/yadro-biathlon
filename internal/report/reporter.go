package report

import (
	"github.com/AleksandrMatsko/yadro-biathlon/internal/config"
	"github.com/AleksandrMatsko/yadro-biathlon/internal/event"
)

type Reporter struct {
	conf      config.BiathlonCompetition
	reporters map[string]*competitorReporter
}

func NewReporter(conf config.BiathlonCompetition) *Reporter {
	return &Reporter{
		conf:      conf,
		reporters: make(map[string]*competitorReporter),
	}
}

func (r *Reporter) NotifyWithEvent(incomingEvent event.Event) {
	reporter, ok := r.reporters[incomingEvent.CompetitorID]
	if !ok && incomingEvent.ID != event.CompetitorRegistration {
		return
	}

	if incomingEvent.ID == event.CompetitorRegistration {
		reporter = newCompetitorReporter()
		r.reporters[incomingEvent.CompetitorID] = reporter
	}

	reporter.NotifyWithEvent(incomingEvent)
}

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
	totalTime *totalTime
}

func newCompetitorReporter() *competitorReporter {
	return &competitorReporter{
		totalTime: newTotalTime(),
	}
}

func (cr *competitorReporter) NotifyWithEvent(e event.Event) {
	cr.totalTime.NotifyWithEvent(e)
}

func (cr *competitorReporter) createRecord() reportRecord {
	record := reportRecord{}

	record.totalTime, record.finalState = cr.totalTime.GetTotalTime()

	return record
}
