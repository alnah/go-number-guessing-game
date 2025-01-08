package timer

import "time"

type Timer interface {
	Now() time.Time
}

type DefaultTimer struct{}

func (t *DefaultTimer) Now() time.Time {
	return time.Now()
}

type GameTimer struct {
	Timer       Timer
	StartTime   *time.Time
	EndTime     *time.Time
	ElapsedTime *time.Duration
}

func (g *GameTimer) Start() {
	now := g.Now()
	g.StartTime = &now
}

func (g *GameTimer) End() time.Duration {
	now := g.Now()
	g.EndTime = &now

	return g.EndTime.Sub(*g.StartTime)
}

func (g *GameTimer) Now() time.Time {
	return g.Timer.Now()
}

func NewGameTimer() GameTimer {
	return GameTimer{Timer: &DefaultTimer{}}
}
