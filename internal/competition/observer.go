package competition

import (
	"github.com/AleksandrMatsko/yadro-biathlon/internal/event"
)

// Observer is the interface for observing and handling incoming and outgoing events.
type Observer interface {
	// NotifyWithEvent should update Observer's state according to the given event.
	NotifyWithEvent(event.Event)
}

// ComposedObserver - is an Observer, that can contain several observers.
type ComposedObserver struct {
	observers []Observer
}

// NewComposedObserver - creates new ComposedObserver.
func NewComposedObserver() *ComposedObserver {
	return &ComposedObserver{
		observers: make([]Observer, 0),
	}
}

// NotifyWithEvent call this method for each Observer.
func (c *ComposedObserver) NotifyWithEvent(e event.Event) {
	for _, observer := range c.observers {
		observer.NotifyWithEvent(e)
	}
}

// AddObservers add passed observers to ComposedObserver,
// so they also will be notified when event is got.
// ComposedObserver will not add nil observer.
func (c *ComposedObserver) AddObservers(observers ...Observer) *ComposedObserver {
	for i := range observers {
		if observers[i] != nil {
			c.observers = append(c.observers, observers[i])
		}
	}

	return c
}
