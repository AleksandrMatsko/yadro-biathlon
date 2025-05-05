package competition

import (
	"fmt"
	"time"

	"github.com/AleksandrMatsko/yadro-biathlon/internal/config"
)

type rules struct {
	Laps          uint32
	MaxStartDelta time.Duration
}

func fromConfig(conf config.BiathlonCompetition) (rules, error) {
	parsedTime, err := time.Parse(time.TimeOnly, conf.StartDelta)
	if err != nil {
		return rules{}, fmt.Errorf("failed to parse startDelta: %w", err)
	}

	return rules{
		Laps:          conf.Laps,
		MaxStartDelta: parsedTime.Sub(time.Date(0, time.January, 1, 0, 0, 0, 0, parsedTime.Location())),
	}, nil
}
