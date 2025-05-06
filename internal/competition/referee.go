package competition

import (
	"time"

	"github.com/AleksandrMatsko/yadro-biathlon/internal/event"
	"github.com/AleksandrMatsko/yadro-biathlon/internal/event/parser"
)

// referees is responsible for generating outgoing events of different types.
type referees struct {
	rules              rules
	root               Observer
	competitorReferees map[string]Observer
}

func newReferees(rules rules, rootObserver Observer) *referees {
	return &referees{
		rules:              rules,
		root:               rootObserver,
		competitorReferees: make(map[string]Observer),
	}
}

func (r *referees) NotifyWithEvent(e event.Event) {
	obs, ok := r.competitorReferees[e.CompetitorID]
	if !ok && e.ID != event.CompetitorRegistration {
		return
	}

	if e.ID == event.CompetitorRegistration {
		obs = NewComposedObserver().
			AddObservers(
				newObserverStartReferee(r.root, r.rules.MaxStartDelta),
				newObserveFinishReferee(r.root, r.rules.Laps))
		r.competitorReferees[e.CompetitorID] = obs
	}

	obs.NotifyWithEvent(e)
}

// observeStartReferee is responsible for checking competitor disqualification
// because of starting too late.
type observeStartReferee struct {
	root              Observer
	maxStartDelta     time.Duration
	started           bool
	assignedStartTime time.Time
}

func newObserverStartReferee(rootObserver Observer, maxStartDelta time.Duration) *observeStartReferee {
	return &observeStartReferee{
		root:          rootObserver,
		maxStartDelta: maxStartDelta,
	}
}

func (r *observeStartReferee) NotifyWithEvent(e event.Event) {
	if r.started {
		return
	}

	if e.ID == event.StartTimeAssignment {
		r.assignedStartTime, _ = parser.ParseTime(e.Extra)
		return
	}

	if e.ID == event.CompetitorStarted {
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

type competitorState string

const (
	notStarted     competitorState = "NotStarted"
	running        competitorState = "Running"
	disqualified   competitorState = "Disqualified"
	cannotContinue competitorState = "CannotContinue"
	finished       competitorState = "Finished"
)

// observeFinishReferee is responsible for checking competitor race finish.
type observeFinishReferee struct {
	root            Observer
	lapsCount       uint32
	lapsCompleted   uint32
	competitorState competitorState
}

func newObserveFinishReferee(rootObserver Observer, laps uint32) *observeFinishReferee {
	return &observeFinishReferee{
		root:            rootObserver,
		lapsCount:       laps,
		lapsCompleted:   0,
		competitorState: notStarted,
	}
}

func (r *observeFinishReferee) NotifyWithEvent(e event.Event) {
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
