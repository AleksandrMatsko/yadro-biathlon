// competition is a package that contains biathlon competition logic,
// need to perform it.
package competition

import (
	"fmt"

	"github.com/AleksandrMatsko/yadro-biathlon/internal/config"
	"github.com/AleksandrMatsko/yadro-biathlon/internal/event"
)

// Biathlon represents biathlon competition, that handles incoming events.
type Biathlon struct {
	rules    rules
	observer Observer
}

// NewBiathlon creates new Biathlon with given config and observer.
// Use Observer if you need custom events handling.
// Biathlon will log recieved events to stdout.
func NewBiathlon(conf config.BiathlonCompetition, observer Observer) (*Biathlon, error) {
	competitionRules, err := fromConfig(conf)
	if err != nil {
		return nil, fmt.Errorf("bad config %w", err)
	}

	biathlonReferees := newReferees(competitionRules, observer)

	return &Biathlon{
		rules:    competitionRules,
		observer: NewComposedObserver().AddObservers(observer, biathlonReferees),
	}, nil
}

// HandleEvent handles given event and notify observer with the same event.
func (b *Biathlon) HandleEvent(e event.Event) {
	b.observer.NotifyWithEvent(e)
}
