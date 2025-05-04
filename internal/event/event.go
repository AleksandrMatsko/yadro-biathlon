// event contains types and functions need to work with incoming and outgoing events.
package event

import (
	"fmt"
	"time"
)

// EventID is used to determine the special event to happen.
type EventID uint8

const (
	CompetitorRegistration     EventID = 1
	StartTimeAssignment        EventID = 2
	CompetitorOnStartine       EventID = 3
	CompetitorStarted          EventID = 4
	CompetitorOnFiringRange    EventID = 5
	TargetHit                  EventID = 6
	CompetitorLeftFiringRange  EventID = 7
	CompetitorEnterPenaltyLaps EventID = 8
	CompetitorLeftPenaltyLaps  EventID = 9
	CompetitorEndedMainLap     EventID = 10
	CompetitorCannotContinue   EventID = 11
	CompetitorDisqualified     EventID = 32
	CompetitorFinished         EventID = 33
)

// Event respesents incoming or outgoing event happened with competitor.
type Event struct {
	// Time of event happened.
	Time time.Time
	// ID of Event.
	ID EventID
	// CompetitorID to which Event relates.
	CompetitorID string
	// Extra arguments in single string.
	Extra string
}

// ValidIncomingEventID checks if the given value is a valid event id.
func ValidIncomingEventID(candidate uint8) bool {
	return candidate >= uint8(CompetitorRegistration) && candidate <= uint8(CompetitorCannotContinue)
}

// TimeFormat use to parse.
const TimeFormat = "15:05:04.000"

// FormatTime according to TimeFormat.
func (e Event) FormatTime() string {
	return fmt.Sprintf("[%s]", e.Time.Format(TimeFormat))
}

// String method to implement Stringer interface.
func (e Event) String() string {
	formattedTime := e.FormatTime() + " "

	eventMsg := "unknown event"
	switch e.ID {
	case CompetitorRegistration:
		eventMsg = fmt.Sprintf("The competitor(%s) registered", e.CompetitorID)
	case StartTimeAssignment:
		eventMsg = fmt.Sprintf("The start time for the competitor(%s) was set by a draw to %s", e.CompetitorID, e.Extra)
	case CompetitorOnStartine:
		eventMsg = fmt.Sprintf("The competitor(%s) is on the start line", e.CompetitorID)
	case CompetitorStarted:
		eventMsg = fmt.Sprintf("The competitor(%s) has started", e.CompetitorID)
	case CompetitorOnFiringRange:
		eventMsg = fmt.Sprintf("The competitor(%s) is on the firing range(%s)", e.CompetitorID, e.Extra)
	case TargetHit:
		eventMsg = fmt.Sprintf("The target(%s) has been hit by competitor(%s)", e.Extra, e.CompetitorID)
	case CompetitorLeftFiringRange:
		eventMsg = fmt.Sprintf("The competitor(%s) left the firing range", e.CompetitorID)
	case CompetitorEnterPenaltyLaps:
		eventMsg = fmt.Sprintf("The competitor(%s) entered the penalty laps", e.CompetitorID)
	case CompetitorLeftPenaltyLaps:
		eventMsg = fmt.Sprintf("The competitor(%s) left the penalty laps", e.CompetitorID)
	case CompetitorEndedMainLap:
		eventMsg = fmt.Sprintf("The competitor(%s) ended the main lap", e.CompetitorID)
	case CompetitorCannotContinue:
		eventMsg = fmt.Sprintf("The competitor(%s) can`t continue: %s", e.CompetitorID, e.Extra)
	case CompetitorDisqualified:
		eventMsg = fmt.Sprintf("The competitor(%s) has been disqualified", e.CompetitorID)
	case CompetitorFinished:
		eventMsg = fmt.Sprintf("The competitor(%s) has finished", e.CompetitorID)
	}

	return formattedTime + eventMsg
}
