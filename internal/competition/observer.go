package competition

import (
	"github.com/AleksandrMatsko/yadro-biathlon/internal/event"
)

// Observer is the interface for observing and handling incoming and outgoing events.
type Observer interface {
	// NotifyWithEvent should update Observer's state according to the given event.
	NotifyWithEvent(event.Event)
}

// Composed - is an Observer, that can contain several observers.
type Composed struct {
	observers []Observer
}

// NewComposed - creates new Composed Observer.
func NewComposed() *Composed {
	return &Composed{
		observers: make([]Observer, 0),
	}
}

// NotifyWithEvent call this method for each Observer.
func (c *Composed) NotifyWithEvent(e event.Event) {
	for _, observer := range c.observers {
		observer.NotifyWithEvent(e)
	}
}

// AddObservers add passed observers to Composed Observer,
// so they also will be notified when event is got.
func (c *Composed) AddObservers(observer ...Observer) *Composed {
	c.observers = append(c.observers, observer...)
	return c
}
