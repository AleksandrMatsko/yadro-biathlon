package competition

import (
	"time"

	"github.com/AleksandrMatsko/yadro-biathlon/internal/event"
	"github.com/AleksandrMatsko/yadro-biathlon/internal/event/parser"
)

type Referees struct {
	rules              rules
	root               Observer
	competitorReferees map[string]Observer
}

func NewReferees(rules rules, rootObserver Observer) *Referees {
	return &Referees{
		rules:              rules,
		root:               rootObserver,
		competitorReferees: make(map[string]Observer),
	}
}

func (r *Referees) NotifyWithEvent(e event.Event) {
	obs, ok := r.competitorReferees[e.CompetitorID]
	if !ok && e.ID != event.CompetitorRegistration {
		return
	}

	if e.ID == event.CompetitorRegistration {
		obs = NewComposed().
			AddObserver(newWatchStartReferee(r.root, r.rules.MaxStartDelta)).
			AddObserver(newWatchFinishReferee(r.root, r.rules.Laps))
		r.competitorReferees[e.CompetitorID] = obs
	}

	obs.NotifyWithEvent(e)
}

type watchStartReferee struct {
	root              Observer
	maxStartDelta     time.Duration
	started           bool
	assignedStartTime time.Time
}

func newWatchStartReferee(rootObserver Observer, maxStartDelta time.Duration) *watchStartReferee {
	return &watchStartReferee{
		root:          rootObserver,
		maxStartDelta: maxStartDelta,
	}
}

func (r *watchStartReferee) NotifyWithEvent(e event.Event) {
	if r.started {
		return
	}

	if !r.started && e.ID == event.StartTimeAssignment {
		r.assignedStartTime, _ = parser.ParseTime(e.Extra)
		return
	}

	if !r.started && e.ID == event.CompetitorStarted {
		r.started = true

		actualDelta := e.Time.Sub(r.assignedStartTime)
		if actualDelta > r.maxStartDelta {
			r.root.NotifyWithEvent(event.Event{
				Time:         e.Time,
				ID:           event.CompetitorDisqualified,
				CompetitorID: e.CompetitorID,
			})
		}
	}
}

type watchFinishState string

const (
	notStarted     watchFinishState = "NotStarted"
	running        watchFinishState = "Running"
	disqualified   watchFinishState = "Disqualified"
	cannotContinue watchFinishState = "CannotContinue"
	finished       watchFinishState = "Finished"
)

type watchFinishReferee struct {
	root            Observer
	lapsCount       uint32
	lapsCompleted   uint32
	competitorState watchFinishState
}

func newWatchFinishReferee(rootObserver Observer, laps uint32) *watchFinishReferee {
	return &watchFinishReferee{
		root:            rootObserver,
		lapsCount:       laps,
		lapsCompleted:   0,
		competitorState: notStarted,
	}
}

func (r *watchFinishReferee) NotifyWithEvent(e event.Event) {
	if r.competitorState == notStarted && e.ID == event.CompetitorStarted {
		r.competitorState = running
		return
	}

	if r.competitorState == running && e.ID == event.CompetitorCannotContinue {
		r.competitorState = cannotContinue
		return
	}

	if r.competitorState == running && e.ID == event.CompetitorDisqualified {
		r.competitorState = disqualified
		return
	}

	if r.competitorState == running && e.ID == event.CompetitorEndedMainLap {
		r.lapsCompleted += 1

		if r.lapsCount == r.lapsCompleted {
			r.competitorState = finished
			r.root.NotifyWithEvent(event.Event{
				Time:         e.Time,
				ID:           event.CompetitorFinished,
				CompetitorID: e.CompetitorID,
			})
		}

		return
	}
}
