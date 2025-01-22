// Package game orchestrates the core game logic, including managing game states,
// handling player turns, and validating game levels and attempts.
package game

import (
	"fmt"
	"math"
	"math/rand/v2"
)

// TurnsLengthError represents an error occurring when the number of turns
// exceeds the maximum allowed attempts.
type TurnsLengthError struct {
	Turns       Turns
	MaxAttempts int
}

// Error returns a formatted error message indicating the number of turns
// and the maximum attempts allowed.
func (e *TurnsLengthError) Error() string {
	message := "Turns length (%d) must be less than max attempts (%d)"
	return fmt.Sprintf(message, len(e.Turns), e.MaxAttempts)
}

// NewTurnsLengthError creates a new TurnsLengthError with the specified
// turns and maximum attempts.
func NewTurnsLengthError(turns Turns, maxAttempts int) error {
	return &TurnsLengthError{Turns: turns, MaxAttempts: maxAttempts}
}

// LevelError represents an error occurring when an invalid game level is
// provided.
type LevelError struct{}

// Error returns a message indicating the valid levels for the game.
func (e *LevelError) Error() string {
	return `Level must be "Easy", "Medium" or "Hard".`
}

// NewLevelError creates a new LevelError for testing.
func NewLevelError() error {
	return &LevelError{}
}

// MaxAttemptsError represents an error occurring when the maximum attempts
// for a level are not met.
type MaxAttemptsError struct{}

// Error returns a message indicating the valid maximum attempts for each level.
func (e *MaxAttemptsError) Error() string {
	return "Max attempts must be 10 (Easy), 5 (Medium) or 3 (Hard)."
}

// NewMaxAttemptsError creates a new MaxAttemptsError for testing.
func NewMaxAttemptsError() error {
	return &MaxAttemptsError{}
}

// RandomNumberFoundError represents an error occurring when the random number
// has already been found.
type RandomNumberFoundError struct{}

// Error returns a message indicating that the random number has already been found.
func (e *RandomNumberFoundError) Error() string {
	return "Random number already found."
}

// NewRandomNumberFoundError creates a new RandomNumberFoundError for testing.
func NewRandomNumberFoundError() error {
	return &RandomNumberFoundError{}
}

// EmptyTurnsError represents an error occurring when the turns slice is empty.
type EmptyTurnsError struct{}

// Error returns a message indicating that the turns slice must not be empty.
func (e *EmptyTurnsError) Error() string {
	return "Turns must be a non-empty slice."
}

// NewEmptyTurnsError creates a new EmptyTurnsError for testing.
func NewEmptyTurnsError() error {
	return &EmptyTurnsError{}
}

// NewRandomNumber generates and returns a new random number between 1 and 100.
func NewRandomNumber() int {
	return rand.IntN(100) + 1
}

// Turn represents a single turn in the game, it holds the guessed number,
// the outcome of the guess, and the difference from the random number.
type Turn struct {
	GuessNumber int
	Outcome     *int
	Difference  *int
}

// Turns holds a collection of turns.
type Turns []Turn

// GameState holds the current state of the game, including the level,
// maximum attempts, the random number, and the turns taken.
type GameState struct {
	Level        string
	MaxAttempts  int
	RandomNumber int
	Turns        Turns
}

// PlayTurn processes a player's turn, validating the game state and updating
// the outcome and difference for the turn.
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

// GetAttempts returns the number of attempts made in the game.
func (gs *GameState) GetAttempts() int {
	return len(gs.Turns)
}

// GetLastTurn retrieves the last turn made in the game, returning an error
// if no turns have been made.
func (gs *GameState) GetLastTurn() (Turn, error) {
	if len(gs.Turns) == 0 {
		return Turn{}, NewEmptyTurnsError()
	}
	return gs.Turns[len(gs.Turns)-1], nil
}

// NoMoreAttempts checks if the maximum number of attempts has been reached.
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
