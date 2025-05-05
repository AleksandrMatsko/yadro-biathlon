package competitor

import (
	"fmt"
	"time"

	"github.com/AleksandrMatsko/yadro-biathlon/internal/event"
	"github.com/AleksandrMatsko/yadro-biathlon/internal/event/parser"
)

type stateType uint8

const (
	registered stateType = iota
	startScheduled
	readyToStart
	runningMainLap
	shooting
	runningPenaltyLap
	notStarted
	notFinised
	finished
	dataDamaged
)

func (st stateType) String() string {
	switch st {
	case registered:
		return "registered"
	case startScheduled:
		return "startScheduled"
	case readyToStart:
		return "readyToStart"
	case runningMainLap:
		return "runningMainLap"
	case shooting:
		return "shooting"
	case runningPenaltyLap:
		return "runningPenaltyLap"
	case notStarted:
		return "notStarted"
	case notFinised:
		return "notFinished"
	case finished:
		return "finished"
	case dataDamaged:
		return "dataDamaged"
	default:
		return fmt.Sprintf("unknown state: %d", st)
	}
}

type Competitor struct {
	competitorID             string
	state                    stateType
	completedLaps            uint32
	completedShootingsPerLap uint32
	totalLaps                uint32
	firingLines              uint32
	assignedStartTime        time.Time
	actualStartTime          time.Time
}

func NewCompetitor(competitorID string, laps, firingLines uint32) *Competitor {
	return &Competitor{
		competitorID:             competitorID,
		state:                    registered,
		completedLaps:            0,
		completedShootingsPerLap: 0,
		totalLaps:                laps,
		firingLines:              firingLines,
	}
}

func (c *Competitor) ChangeState(incomingEvent event.Event) error {
	switch c.state {
	case registered:
		return c.fromRegistered(incomingEvent)
	case readyToStart:
		return c.fromStartScheduled(incomingEvent)
	default:
		return nil
	}
}

func (c *Competitor) fromRegistered(incomingEvent event.Event) error {
	if incomingEvent.ID != event.StartTimeAssignment {
		return fmt.Errorf("bad competitor state")
	}

	parsedTime, err := parser.ParseTime(incomingEvent.Extra)
	if err != nil {
		c.state = dataDamaged
		return err
	}

	c.assignedStartTime = parsedTime
	c.state = startScheduled

	return nil
}

func (c *Competitor) fromStartScheduled(incomingEvent event.Event) error {
	if incomingEvent.ID == event.CompetitorOnStartine {
		c.state = readyToStart
	}

	return nil
}
