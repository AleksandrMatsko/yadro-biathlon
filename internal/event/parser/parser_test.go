package parser

import (
	"fmt"
	"iter"
	"strings"
	"testing"
	"time"

	"github.com/AleksandrMatsko/yadro-biathlon/internal/event"
	"github.com/stretchr/testify/assert"
)

func TestLines(t *testing.T) {
	expectedLines := []string{
		"abra",
		"cadabra",
		"hello",
		"world",
	}

	givenString := strings.Join(expectedLines, "\n")
	reader := strings.NewReader(givenString)

	lines, retErrFunc := Lines(reader)

	next, stop := iter.Pull[string](lines)
	defer stop()

	for _, expectedLine := range expectedLines {
		gotLine, ok := next()

		assert.True(t, ok)
		assert.Equal(t, expectedLine, gotLine)
	}

	_, ok := next()
	assert.False(t, ok)
	assert.Nil(t, retErrFunc())
}

func TestParsedLines(t *testing.T) {
	t.Run("with valid events", func(t *testing.T) {
		t.Parallel()

		validEvents := strings.Join([]string{
			"[09:05:59.867] 1 1",
			"[09:15:00.841] 2 1 09:30:00.000",
			"[09:29:45.734] 3 2",
			"[09:49:31.659] 5 abra 1",
			"[09:49:31.700] 5 abra hello world",
		}, "\n")

		location := time.Now().UTC().Location()

		parsedValidEvents := []event.Event{
			{
				Time:         time.Date(0, time.January, 1, 9, 5, 59, 867_000_000, location),
				ID:           event.CompetitorRegistration,
				CompetitorID: "1",
			},
			{
				Time:         time.Date(0, time.January, 1, 9, 15, 0, 841_000_000, location),
				ID:           event.StartTimeAssignment,
				CompetitorID: "1",
				Extra:        "09:30:00.000",
			},
			{
				Time:         time.Date(0, time.January, 1, 9, 29, 45, 734_000_000, location),
				ID:           event.CompetitorOnStartine,
				CompetitorID: "2",
			},
			{
				Time:         time.Date(0, time.January, 1, 9, 49, 31, 659_000_000, location),
				ID:           event.CompetitorOnFiringRange,
				CompetitorID: "abra",
				Extra:        "1",
			},
			{
				Time:         time.Date(0, time.January, 1, 9, 49, 31, 700_000_000, location),
				ID:           event.CompetitorOnFiringRange,
				CompetitorID: "abra",
				Extra:        "hello world",
			},
		}

		reader := strings.NewReader(validEvents)

		lines, retErrFunc := Lines(reader)

		next, stop := iter.Pull2(ParsedLines(lines))
		defer stop()

		for _, expectedEvent := range parsedValidEvents {
			gotEvent, err, ok := next()

			assert.True(t, ok)
			assert.Nil(t, err)
			assert.Equal(t, expectedEvent, gotEvent)
		}

		_, _, ok := next()
		assert.False(t, ok)

		assert.Nil(t, retErrFunc())
	})

	t.Run("with invalid events", func(t *testing.T) {
		badEvents := strings.Join([]string{
			"hello world",
			"hello world program",
			"[09:05:59.867] 1",
			"[09:05:59.867] 0 2",
		}, "\n")

		reader := strings.NewReader(badEvents)

		lines, retErrFunc := Lines(reader)
		for gotEvent, err := range ParsedLines(lines) {
			assert.Equal(t, event.Event{}, gotEvent)
			assert.NotNil(t, err)
		}

		assert.Nil(t, retErrFunc())
	})
}

func TestParseSingleLine(t *testing.T) {
	t.Run("with less parameters than need", func(t *testing.T) {
		t.Parallel()

		type testcase struct {
			name       string
			line       string
			paramCount int
		}

		cases := []testcase{
			{
				name:       "with 0 params",
				line:       "",
				paramCount: 0,
			},
			{
				name:       "with 1 param",
				line:       "[09:05:59.867]",
				paramCount: 1,
			},
			{
				name:       "with 2 params",
				line:       "[09:05:59.867] 1",
				paramCount: 2,
			},
		}

		for i, singleCase := range cases {
			t.Run(fmt.Sprintf("case %d: %s", i+1, singleCase.name), func(t *testing.T) {
				gotEvent, err := ParseSingleLine(singleCase.line)

				assert.Equal(t, event.Event{}, gotEvent)
				assert.Equal(t, err, fmt.Errorf("not enough arguments, expected at least 3, got %d", singleCase.paramCount))
			})
		}
	})

	t.Run("with bad timestamp", func(t *testing.T) {
		t.Parallel()

		gotEvent, err := ParseSingleLine("09.05.59,867")

		assert.Equal(t, event.Event{}, gotEvent)
		assert.NotNil(t, err)
	})

	t.Run("with bad event id", func(t *testing.T) {
		t.Parallel()

		t.Run("with unparseable value", func(t *testing.T) {
			gotEvent, err := ParseSingleLine("[09:05:59.867] abra 1")

			assert.Equal(t, event.Event{}, gotEvent)
			assert.NotNil(t, err)
		})

		t.Run("with invalid value", func(t *testing.T) {
			gotEvent, err := ParseSingleLine("[09:05:59.867] 0 1")

			assert.Equal(t, event.Event{}, gotEvent)
			assert.Equal(t, err, fmt.Errorf("failed to parse event id: %w", fmt.Errorf("invalid event id: 0")))
		})
	})
}
