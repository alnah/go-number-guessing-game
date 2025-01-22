// Package timer provides a game timer to track elapsed time during game
// sessions. It also offers a wrapper around the Go built-in time package's Now
// method for testing.
package timer

import "time"

// Timer defines an interface for getting the current time, useful for testing.
type Timer interface {
	Now() time.Time
}

// DefaultTimer implements the Timer interface to return the current local time.
type DefaultTimer struct{}

// Now returns the current local time.
func (t *DefaultTimer) Now() time.Time {
	return time.Now()
}

// GameTimer manages the start, end, and elapsed time for a game session,
// allowing for user feedback and score storage.
type GameTimer struct {
	Timer       Timer
	StartTime   *time.Time
	EndTime     *time.Time
	ElapsedTime *time.Duration
}

// Start records the current time as the start time.
func (g *GameTimer) Start() {
	now := g.Now()
	g.StartTime = &now
}

// End records the current time as the end time and calculates the elapsed
// duration.
func (g *GameTimer) End() time.Duration {
	now := g.Now()
	g.EndTime = &now

	return g.EndTime.Sub(*g.StartTime).Truncate(time.Second)
}

// Now retrieves the current time using the Timer interface for testing.
func (g *GameTimer) Now() time.Time {
	return g.Timer.Now()
}

// NewGameTimer creates a GameTimer instance with a DefaultTimer for tracking
// game time.
func NewGameTimer() GameTimer {
	return GameTimer{Timer: &DefaultTimer{}}
}
