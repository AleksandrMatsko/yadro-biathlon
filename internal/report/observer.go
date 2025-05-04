package report

import (
	"github.com/AleksandrMatsko/yadro-biathlon/internal/event"
)

type Observer interface {
	NotifyWithEvent(event.Event)
}
