package report

import (
	"github.com/AleksandrMatsko/yadro-biathlon/internal/event"
)

type Reporter struct {
	reporters map[string]*competitorReporter
}

func (r *Reporter) HandleEvent(incomingEvent event.Event) {
	reporter, ok := r.reporters[incomingEvent.CompetitorID]
	if !ok {
		if incomingEvent.ID == event.CompetitorRegistration {
			r.reporters[incomingEvent.CompetitorID] = newCompetitorReporter()
		}

		return
	}

	reporter.NotifyWithEvent(incomingEvent)
}

type competitorReporter struct {
	totalTime *totalTime
}

func newCompetitorReporter() *competitorReporter {
	return &competitorReporter{}
}

func (cr *competitorReporter) NotifyWithEvent(e event.Event) {
	cr.totalTime.NotifyWithEvent(e)
}

func (cr *competitorReporter) createRecord() reportRecord {
	record := reportRecord{}

	record.totalTime, record.finalState = cr.totalTime.GetTotalTime()

	return record
}
