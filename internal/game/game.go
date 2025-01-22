package game

import (
	"fmt"
	"math"
	"math/rand/v2"
)

type TurnsLengthError struct {
	Turns       Turns
	MaxAttempts int
}

func (e *TurnsLengthError) Error() string {
	message := "Turns length (%d) must be less than max attempts (%d)"
	return fmt.Sprintf(message, len(e.Turns), e.MaxAttempts)
}

func NewTurnsLengthError(turns Turns, maxAttempts int) error {
	return &TurnsLengthError{Turns: turns, MaxAttempts: maxAttempts}
}

type LevelError struct{}

func (e *LevelError) Error() string {
	return `Level must be "Easy", "Medium" or "Hard".`
}

func NewLevelError() error {
	return &LevelError{}
}

type MaxAttemptsError struct{}

func (e *MaxAttemptsError) Error() string {
	return "Max attempts must be 10 (Easy), 5 (Medium) or 3 (Hard)."
}

func NewMaxAttemptsError() error {
	return &MaxAttemptsError{}
}

type RandomNumberFoundError struct{}

func (e *RandomNumberFoundError) Error() string {
	return "Random number already found."
}

func NewRandomNumberFoundError() error {
	return &RandomNumberFoundError{}
}

type EmptyTurnsError struct{}

func (e *EmptyTurnsError) Error() string {
	return "Turns must be an non-empty slice."
}

func NewEmptyTurnsError() error {
	return &EmptyTurnsError{}
}

func NewRandomNumber() int {
	return rand.IntN(100) + 1
}

type Turn struct {
	GuessNumber int
	Outcome     *int
	Difference  *int
}

type Turns []Turn

type GameState struct {
	Level        string
	MaxAttempts  int
	RandomNumber int
	Turns        Turns
}

func (gs *GameState) PlayTurn(turn Turn) error {
	if err := gs.validateLevelAndMaxAttempts(); err != nil {
		return err
	}

	if err := gs.validateRandomNumberNotFound(); err != nil {
		return err
	}

	gs.newOutcome(&turn)
	gs.newDifference(&turn)
	gs.compareNumbers(turn)
	gs.getDifference(turn)
	gs.appendTurn(turn)

	if err := gs.validateMaxLengthTurn(); err != nil {
		return err
	}

	return nil
}

func (gs *GameState) GetAttempts() int {
	return len(gs.Turns)
}

func (gs *GameState) GetLastTurn() (Turn, error) {
	if len(gs.Turns) == 0 {
		return Turn{}, NewEmptyTurnsError()
	}
	return gs.Turns[len(gs.Turns)-1], nil
}

func (gs *GameState) NoMoreAttempts() bool {
	return gs.GetAttempts() == gs.MaxAttempts
}

func (gs *GameState) validateLevelAndMaxAttempts() error {
	levelMaxAttempts := map[string]int{
		"Easy":   10,
		"Medium": 5,
		"Hard":   3,
	}

	expectedMaxAttempts, exists := levelMaxAttempts[gs.Level]
	if !exists {
		return NewLevelError()
	}

	if gs.MaxAttempts != expectedMaxAttempts {
		return NewMaxAttemptsError()
	}

	return nil
}

func (gs *GameState) validateRandomNumberNotFound() error {
	if len(gs.Turns) > 0 {
		lastTurn := gs.Turns[len(gs.Turns)-1]
		if *lastTurn.Outcome == 0 {
			return NewRandomNumberFoundError()
		}
	}
	return nil
}

func (gs *GameState) validateMaxLengthTurn() error {
	if len(gs.Turns) > gs.MaxAttempts {
		return NewTurnsLengthError(gs.Turns, gs.MaxAttempts)
	}
	return nil
}

func (*GameState) newOutcome(turn *Turn) {
	turn.Outcome = new(int)
}

func (gs *GameState) newDifference(turn *Turn) {
	turn.Difference = new(int)
}

func (gs *GameState) appendTurn(turn Turn) {
	gs.Turns = append(gs.Turns, turn)
}

func (gs *GameState) compareNumbers(turn Turn) {
	switch {
	case turn.GuessNumber > gs.RandomNumber:
		*turn.Outcome = -1
	case turn.GuessNumber < gs.RandomNumber:
		*turn.Outcome = 1
	default:
		*turn.Outcome = 0
	}
}

func (gs *GameState) getDifference(turn Turn) {
	difference := int(math.Abs(float64(gs.RandomNumber - turn.GuessNumber)))
	*turn.Difference = difference
}
