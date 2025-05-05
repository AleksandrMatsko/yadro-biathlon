package competition

import (
	"fmt"

	"github.com/AleksandrMatsko/yadro-biathlon/internal/config"
	"github.com/AleksandrMatsko/yadro-biathlon/internal/event"
)

type Biathlon struct {
	rules    rules
	observer Observer
}

func NewBiathlon(conf config.BiathlonCompetition, observer Observer) (*Biathlon, error) {
	competitionRules, err := fromConfig(conf)
	if err != nil {
		return nil, fmt.Errorf("bad config %w", err)
	}

	composed := NewComposed().
		AddObservers(NewLogger(), observer)

	referees := NewReferees(competitionRules, composed)

	return &Biathlon{
		rules:    competitionRules,
		observer: composed.AddObservers(referees),
	}, nil
}

func (b *Biathlon) HandleEvent(e event.Event) {
	b.observer.NotifyWithEvent(e)
}
