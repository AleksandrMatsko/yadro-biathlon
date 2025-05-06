package competition

import (
	"testing"
	"time"

	"github.com/AleksandrMatsko/yadro-biathlon/internal/config"
	"github.com/stretchr/testify/assert"
)

func Test_rules_fromConfig(t *testing.T) {
	t.Run("with bad start delta", func(t *testing.T) {
		conf := config.BiathlonCompetition{
			Laps:       2,
			StartDelta: "hello",
		}

		gotRules, err := fromConfig(conf)

		assert.Equal(t, rules{}, gotRules)
		assert.NotNil(t, err)
	})

	t.Run("with ok delta", func(t *testing.T) {
		conf := config.BiathlonCompetition{
			Laps:       3,
			StartDelta: "00:01:30",
		}

		gotRules, err := fromConfig(conf)

		assert.Nil(t, err)
		assert.Equal(t,
			rules{
				Laps:          conf.Laps,
				MaxStartDelta: time.Minute + 30*time.Second,
			},
			gotRules)
	})
}
