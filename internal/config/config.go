// config contains all types, functions and methods need to work with config.
package config

import (
	"encoding/json"
	"fmt"
	"os"
)

// BiathlonCompetition represents config for biathlon competition.
type BiathlonCompetition struct {
	// Laps - amount of laps for main distance.
	Laps uint32 `json:"laps"`
	// LapLen - length of each main lap.
	LapLen uint32 `json:"lapLen"`
	// PenaltyLen - length of each penalty lap.
	PenaltyLen uint32 `json:"penaltyLen"`
	// FiringLines - number of firing lines in race.
	FiringLines uint32 `json:"firingLines"`
	// Start - planned start time for the first competitor.
	Start string `json:"start"`
	// StartDelta - planned interval between starts
	StartDelta string `json:"startDelta"`
}

// Read given configFileName and fill config.
func Read(configFileName string, config interface{}) error {
	bytes, err := os.ReadFile(configFileName)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bytes, config)
	if err != nil {
		return fmt.Errorf("failed to parse config file: %s, err: %w", configFileName, err)
	}

	return nil
}

// Print given config to stdout.
func Print(config interface{}) {
	bytes, _ := json.Marshal(config)
	fmt.Println(string(bytes))
}
