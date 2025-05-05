package competition

import (
	"fmt"

	"github.com/AleksandrMatsko/yadro-biathlon/internal/event"
)

type Logger struct{}

func NewLogger() *Logger {
	return &Logger{}
}

func (l *Logger) NotifyWithEvent(e event.Event) {
	fmt.Println(e)
}
