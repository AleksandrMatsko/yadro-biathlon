package event

import (
	"fmt"
	"io"
	"os"
)

// Logger implements competition.Observer.
type Logger struct {
	writer io.Writer
}

// NewLogger creates new logger.
func NewLogger(w io.Writer) *Logger {
	if w == nil {
		w = os.Stdout
	}

	return &Logger{
		writer: w,
	}
}

// NotifyWithEvent logs received event.
func (l *Logger) NotifyWithEvent(e Event) {
	fmt.Fprintln(l.writer, e)
}
