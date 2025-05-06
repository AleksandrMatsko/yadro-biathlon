package competition

import (
	"testing"
	"time"

	mock_observer "github.com/AleksandrMatsko/yadro-biathlon/internal/competition/mocks"
	"github.com/AleksandrMatsko/yadro-biathlon/internal/event"
	"github.com/AleksandrMatsko/yadro-biathlon/internal/event/parser"
	"go.uber.org/mock/gomock"
)

func Test_observeStartReferee(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	startAssignEvent := event.Event{
		ID:    event.StartTimeAssignment,
		Extra: "10:05:30.000",
	}

	t.Run("with no disqualification", func(t *testing.T) {
		t.Parallel()

		rootObserver := mock_observer.NewMockObserver(mockCtrl)

		actualStartTime, _ := parser.ParseTime(startAssignEvent.Extra)
		actualStartTime = actualStartTime.Add(time.Second)

		maxStartDelta := time.Minute
		competitorStartedEvent := event.Event{
			Time: actualStartTime,
			ID:   event.CompetitorStarted,
		}

		referee := newObserverStartReferee(rootObserver, maxStartDelta)

		referee.NotifyWithEvent(startAssignEvent)
		referee.NotifyWithEvent(competitorStartedEvent)

		maxStartDelta = time.Second

		referee = newObserverStartReferee(rootObserver, maxStartDelta)

		referee.NotifyWithEvent(startAssignEvent)
		referee.NotifyWithEvent(competitorStartedEvent)
	})

	t.Run("with disqualification", func(t *testing.T) {
		t.Parallel()

		rootObserver := mock_observer.NewMockObserver(mockCtrl)

		actualStartTime, _ := parser.ParseTime(startAssignEvent.Extra)
		actualStartTime = actualStartTime.Add(time.Minute)

		competitorStartedEvent := event.Event{
			Time:         actualStartTime,
			ID:           event.CompetitorStarted,
			CompetitorID: "vasya",
		}
		maxStartDelta := time.Second

		referee := newObserverStartReferee(rootObserver, maxStartDelta)

		referee.NotifyWithEvent(startAssignEvent)

		rootObserver.EXPECT().NotifyWithEvent(event.Event{
			Time:         competitorStartedEvent.Time,
			ID:           event.CompetitorDisqualified,
			CompetitorID: competitorStartedEvent.CompetitorID,
		}).Times(1)

		referee.NotifyWithEvent(competitorStartedEvent)
	})

	t.Run("event of competitor start received before start scheduling", func(t *testing.T) {
		t.Parallel()

		rootObserver := mock_observer.NewMockObserver(mockCtrl)
		maxStartDelta := time.Second

		referee := newObserverStartReferee(rootObserver, maxStartDelta)

		referee.NotifyWithEvent(startAssignEvent)
	})
}

func Test_observeFinishReferee(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	eventTime, _ := parser.ParseTime("10:30:25.321")

	mainLapEndEvent := event.Event{
		Time:         eventTime,
		ID:           event.CompetitorEndedMainLap,
		CompetitorID: "kolya",
	}

	t.Run("with ok events order", func(t *testing.T) {
		t.Parallel()

		rootObserver := mock_observer.NewMockObserver(mockCtrl)
		var lapsCount uint32 = 3

		referee := newObserveFinishReferee(rootObserver, lapsCount)

		referee.NotifyWithEvent(event.Event{ID: event.CompetitorStarted})

		rootObserver.EXPECT().NotifyWithEvent(event.Event{
			Time:         mainLapEndEvent.Time,
			ID:           event.CompetitorFinished,
			CompetitorID: mainLapEndEvent.CompetitorID,
		}).Times(1)

		referee.NotifyWithEvent(mainLapEndEvent)
		referee.NotifyWithEvent(mainLapEndEvent)
		referee.NotifyWithEvent(mainLapEndEvent)
	})

	t.Run("when competitor is disqualified", func(t *testing.T) {
		t.Parallel()

		rootObserver := mock_observer.NewMockObserver(mockCtrl)

		var lapsCount uint32 = 3

		referee := newObserveFinishReferee(rootObserver, lapsCount)

		referee.NotifyWithEvent(event.Event{ID: event.CompetitorStarted})
		referee.NotifyWithEvent(event.Event{ID: event.CompetitorDisqualified})

		referee.NotifyWithEvent(mainLapEndEvent)
		referee.NotifyWithEvent(mainLapEndEvent)
		referee.NotifyWithEvent(mainLapEndEvent)
	})

	t.Run("when competitor can't continue", func(t *testing.T) {
		t.Parallel()

		rootObserver := mock_observer.NewMockObserver(mockCtrl)

		var lapsCount uint32 = 3

		referee := newObserveFinishReferee(rootObserver, lapsCount)

		referee.NotifyWithEvent(event.Event{ID: event.CompetitorStarted})
		referee.NotifyWithEvent(event.Event{ID: event.CompetitorCannotContinue})

		referee.NotifyWithEvent(mainLapEndEvent)
		referee.NotifyWithEvent(mainLapEndEvent)
		referee.NotifyWithEvent(mainLapEndEvent)
	})
}
