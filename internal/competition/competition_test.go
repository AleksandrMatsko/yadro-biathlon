package competition

import (
	"testing"
	"time"

	mock_observer "github.com/AleksandrMatsko/yadro-biathlon/internal/competition/mocks"
	"github.com/AleksandrMatsko/yadro-biathlon/internal/config"
	"github.com/AleksandrMatsko/yadro-biathlon/internal/event"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestHandleEvent(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	conf := config.BiathlonCompetition{
		Laps:       3,
		StartDelta: "00:01:30",
	}

	observer := mock_observer.NewMockObserver(mockCtrl)

	biathlon, err := NewBiathlon(conf, observer)
	assert.Nil(t, err)

	givenEvent := event.Event{
		Time:         time.Date(0, time.January, 1, 1, 1, 1, 0, time.UTC),
		ID:           event.CompetitorRegistration,
		CompetitorID: "vasya",
	}

	observer.EXPECT().NotifyWithEvent(givenEvent).Times(1)

	biathlon.HandleEvent(givenEvent)
}
