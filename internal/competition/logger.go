package competition

import (
	"fmt"

	"github.com/AleksandrMatsko/yadro-biathlon/internal/event"
)

type logger struct{}

func newLogger() *logger {
	return &logger{}
}

func (l *logger) NotifyWithEvent(e event.Event) {
	fmt.Println(e)
}
