package competition

import (
	"testing"

	mock_observer "github.com/AleksandrMatsko/yadro-biathlon/internal/competition/mocks"
	"github.com/AleksandrMatsko/yadro-biathlon/internal/event"
	"go.uber.org/mock/gomock"
)

func Test_ComposedObserver(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockObservers := []*mock_observer.MockObserver{
		mock_observer.NewMockObserver(mockCtrl),
		mock_observer.NewMockObserver(mockCtrl),
		mock_observer.NewMockObserver(mockCtrl),
		mock_observer.NewMockObserver(mockCtrl),
	}

	composed := NewComposedObserver()

	composed.NotifyWithEvent(event.Event{})

	composed.AddObservers(mockObservers[0])

	mockObservers[0].EXPECT().NotifyWithEvent(event.Event{}).Times(1)

	composed.NotifyWithEvent(event.Event{})

	composed.AddObservers(mockObservers[1], mockObservers[2], mockObservers[3])

	for i := range mockObservers {
		mockObservers[i].EXPECT().NotifyWithEvent(event.Event{}).Times(1)
	}

	composed.NotifyWithEvent(event.Event{})

	composed.AddObservers(nil)

	for i := range mockObservers {
		mockObservers[i].EXPECT().NotifyWithEvent(event.Event{}).Times(1)
	}

	composed.NotifyWithEvent(event.Event{})
}
