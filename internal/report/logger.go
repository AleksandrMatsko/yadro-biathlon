package report

import (
	"fmt"

	"github.com/AleksandrMatsko/yadro-biathlon/internal/event"
)

type Logger struct{}

func (*Logger) NotifyWithEvent(e event.Event) {
	fmt.Println(e)
}
