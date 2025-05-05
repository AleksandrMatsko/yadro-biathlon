package report

import (
	"testing"
	"time"

	"github.com/AleksandrMatsko/yadro-biathlon/internal/event"
	"github.com/stretchr/testify/assert"
)

func Test_totalTime(t *testing.T) {
	t.Run("when competitor finished", func(t *testing.T) {
		t.Parallel()

		tt := newTotalTime()

		tt.NotifyWithEvent(event.Event{
			ID:    event.StartTimeAssignment,
			Extra: "10:03:30.000",
		})
		tt.NotifyWithEvent(event.Event{ID: event.CompetitorStarted})
		tt.NotifyWithEvent(event.Event{
			Time: time.Date(0, time.January, 1, 10, 40, 52, 342_000_000, time.UTC),
			ID:   event.CompetitorFinished,
		})

		duration, finalState := tt.GetTotalTime()
		assert.Equal(t, finished, finalState)
		assert.Equal(t, time.Minute*37+time.Second*22+time.Millisecond*342, duration)
	})

	t.Run("when competitor is disqualified", func(t *testing.T) {
		t.Parallel()

		t.Run("disqualified after competitor start", func(t *testing.T) {
			tt := newTotalTime()

			tt.NotifyWithEvent(event.Event{
				ID:    event.StartTimeAssignment,
				Extra: "10:03:30.000",
			})
			tt.NotifyWithEvent(event.Event{ID: event.CompetitorStarted})
			tt.NotifyWithEvent(event.Event{ID: event.CompetitorDisqualified})
			tt.NotifyWithEvent(event.Event{
				Time: time.Date(0, time.January, 1, 10, 40, 52, 342_000_000, time.UTC),
				ID:   event.CompetitorFinished,
			})

			duration, finalState := tt.GetTotalTime()
			assert.Equal(t, notStarted, finalState)
			assert.Equal(t, time.Duration(0), duration)
		})

		t.Run("disqualified before competitor start", func(t *testing.T) {
			tt := newTotalTime()

			tt.NotifyWithEvent(event.Event{
				ID:    event.StartTimeAssignment,
				Extra: "10:03:30.000",
			})
			tt.NotifyWithEvent(event.Event{ID: event.CompetitorDisqualified})
			tt.NotifyWithEvent(event.Event{ID: event.CompetitorStarted})
			tt.NotifyWithEvent(event.Event{
				Time: time.Date(0, time.January, 1, 10, 40, 52, 342_000_000, time.UTC),
				ID:   event.CompetitorFinished,
			})

			duration, finalState := tt.GetTotalTime()
			assert.Equal(t, notStarted, finalState)
			assert.Equal(t, time.Duration(0), duration)
		})
	})

	t.Run("when competitor can't continue", func(t *testing.T) {
		t.Parallel()

		tt := newTotalTime()

		tt.NotifyWithEvent(event.Event{
			ID:    event.StartTimeAssignment,
			Extra: "10:03:30.000",
		})
		tt.NotifyWithEvent(event.Event{ID: event.CompetitorStarted})
		tt.NotifyWithEvent(event.Event{ID: event.CompetitorCannotContinue})
		tt.NotifyWithEvent(event.Event{
			Time: time.Date(0, time.January, 1, 10, 40, 52, 342_000_000, time.UTC),
			ID:   event.CompetitorFinished,
		})

		duration, finalState := tt.GetTotalTime()
		assert.Equal(t, notFinished, finalState)
		assert.Equal(t, time.Duration(0), duration)
	})

	t.Run("when result is called too early", func(t *testing.T) {
		t.Parallel()

		t.Run("competitor has started", func(t *testing.T) {
			tt := newTotalTime()

			tt.NotifyWithEvent(event.Event{
				ID:    event.StartTimeAssignment,
				Extra: "10:03:30.000",
			})
			tt.NotifyWithEvent(event.Event{ID: event.CompetitorStarted})

			duration, finalState := tt.GetTotalTime()
			assert.Equal(t, notFinished, finalState)
			assert.Equal(t, time.Duration(0), duration)
		})

		t.Run("competitor has not started", func(t *testing.T) {
			tt := newTotalTime()

			tt.NotifyWithEvent(event.Event{
				ID:    event.StartTimeAssignment,
				Extra: "10:03:30.000",
			})

			duration, finalState := tt.GetTotalTime()
			assert.Equal(t, notStarted, finalState)
			assert.Equal(t, time.Duration(0), duration)
		})
	})
}
