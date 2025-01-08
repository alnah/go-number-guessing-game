package timer_test

import (
	"testing"
	"time"

	"github.com/go-number-guessing-game/internal/timer"
	"github.com/stretchr/testify/assert"
)

func TestUnitGameTimer(t *testing.T) {
	t.Run("start time with the current local time", func(t *testing.T) {
		timer := timer.NewGameTimer()

		timer.Start()

		want := time.Now().Second()
		got := timer.StartTime.Second()

		assert.Equal(t, want, got)
	})

	t.Run("end time with the current local time", func(t *testing.T) {
		timer := timer.NewGameTimer()

		timer.Start()
		timer.End()

		want := time.Now().Second()
		got := timer.EndTime.Second()

		assert.Equal(t, want, got)
	})

	t.Run("end time return the elapsed time", func(t *testing.T) {
		timer := timer.GameTimer{Timer: &StubTimer{}}

		timer.Start()
		timer.End()

		want := 10 * time.Second
		got := timer.End()

		assert.Equal(t, want, got)
	})
}

type StubTimer struct {
	calls int
}

func (s *StubTimer) Now() time.Time {
	if s.calls == 0 {
		s.calls++
		return time.Date(2001, 1, 1, 1, 1, 0, 0, time.UTC)
	}
	return time.Date(2001, 1, 1, 1, 1, 10, 0, time.UTC)
}

func TestUnitNewGameTimer(t *testing.T) {
	t.Run("return a game timer", func(t *testing.T) {
		want := timer.GameTimer{Timer: &timer.DefaultTimer{}}
		got := timer.NewGameTimer()

		assert.Equal(t, want, got)
	})
}
