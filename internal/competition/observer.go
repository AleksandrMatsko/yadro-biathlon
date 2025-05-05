package competition

import (
	"github.com/AleksandrMatsko/yadro-biathlon/internal/event"
)

type Observer interface {
	NotifyWithEvent(event.Event)
}

type Composed struct {
	observers []Observer
}

func NewComposed() *Composed {
	return &Composed{}
}

func (c *Composed) NotifyWithEvent(e event.Event) {
	for _, observer := range c.observers {
		observer.NotifyWithEvent(e)
	}
}

func (c *Composed) AddObserver(observer Observer) *Composed {
	c.observers = append(c.observers, observer)
	return c
}
