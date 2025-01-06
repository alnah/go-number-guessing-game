package game

import "fmt"

// TurnsLengthError represents an error that occurs when the number of turns
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

// NewTurnsLengthError creates a new instance of TurnsLengthError with the
// provided turns and maximum attempts.
func NewTurnsLengthError(turns Turns, maxAttempts int) error {
	return &TurnsLengthError{Turns: turns, MaxAttempts: maxAttempts}
}

// LevelError represents an error that occurs when an invalid game level is set.
type LevelError struct{}

// Error returns a message indicating that the level must be one of the
// predefined values: "Easy", "Medium", or "Hard".
func (e *LevelError) Error() string {
	return fmt.Sprint(`Level must be "Easy", "Medium" or "Hard".`)
}

// NewLevelError creates a new instance of LevelError.
func NewLevelError() error {
	return &LevelError{}
}

// MaxAttemptsError represents an error that occurs when the maximum attempts
// for a level are not set correctly.
type MaxAttemptsError struct{}

// Error returns a message indicating the valid maximum attempts for each level.
func (e *MaxAttemptsError) Error() string {
	return fmt.Sprintf("Max attempts must be 10 (Easy), 5 (Medium) or 3 (Hard).")
}

// NewMaxAttemptsError creates a new instance of MaxAttemptsError.
func NewMaxAttemptsError() error {
	return &MaxAttemptsError{}
}

// RandomNumberFoundError represents an error that occurs when a random number
// has already been found.
type RandomNumberFoundError struct{}

// Error returns a message indicating that the random number has already been found.
func (e *RandomNumberFoundError) Error() string {
	return fmt.Sprintf("Random number already found.")
}

// NewRandomNumberFoundError creates a new instance of RandomNumberFoundError.
func NewRandomNumberFoundError() error {
	return &RandomNumberFoundError{}
}

type EmptyTurnsError struct{}

func (e *EmptyTurnsError) Error() string {
	return fmt.Sprint("Turns must be an non-empty slice.")
}

func NewEmptyTurnsError() error {
	return &EmptyTurnsError{}
}

// Turn represents a single turn in the game, containing the guessed number
// and the outcome of the guess.
type Turn struct {
	GuessNumber int
	Outcome     *int
}

// Turns is a slice of Turn, representing all the turns taken in the game.
type Turns []Turn

// GameState holds the current state of the game, including the level,
// maximum attempts, the random number to guess, and the turns taken.
type GameState struct {
	Level        string
	MaxAttempts  int
	RandomNumber int
	Turns        Turns
}

// PlayTurn processes a player's turn, validating the level and maximum attempts,
// checking if the random number has been found, and updating the game state.
func (gs *GameState) PlayTurn(turn Turn) error {
	if err := gs.validateLevelAndMaxAttempts(); err != nil {
		return err
	}

	if err := gs.validateRandomNumberNotFound(); err != nil {
		return err
	}

	gs.newOutcome(&turn)
	gs.compareNumbers(turn)
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

// GetLastTurn retrieves the last turn from the game's state. If no turns
// exist, it returns an error indicating that the turns slice is empty.
func (gs *GameState) GetLastTurn() (Turn, error) {
	if len(gs.Turns) == 0 {
		return Turn{}, NewEmptyTurnsError()
	}
	return gs.Turns[len(gs.Turns)-1], nil
}

// NoMoreAttempts checks if the player has used all available attempts in the game.
// It returns true if the number of attempts made equals the maximum allowed attempts.
func (gs *GameState) NoMoreAttempts() bool {
	return gs.GetAttempts() == gs.MaxAttempts
}

// validateLevelAndMaxAttempts checks if the current level and maximum attempts
// are consistent with the predefined values for each level.
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

// validateRandomNumberNotFound checks if the last turn's outcome indicates
// that the random number has already been found.
func (gs *GameState) validateRandomNumberNotFound() error {
	if len(gs.Turns) > 0 {
		lastTurn := gs.Turns[len(gs.Turns)-1]
		if *lastTurn.Outcome == 0 {
			return NewRandomNumberFoundError()
		}
	}
	return nil
}

// validateMaxLengthTurn checks if the number of turns exceeds the maximum
// allowed attempts.
func (gs *GameState) validateMaxLengthTurn() error {
	if len(gs.Turns) > gs.MaxAttempts {
		return NewTurnsLengthError(gs.Turns, gs.MaxAttempts)
	}
	return nil
}

// newOutcome initializes the outcome of a turn to a new integer pointer.
func (*GameState) newOutcome(turn *Turn) {
	turn.Outcome = new(int)
}

// appendTurn adds a new turn to the game's list of turns.
func (gs *GameState) appendTurn(turn Turn) {
	gs.Turns = append(gs.Turns, turn)
}

// compareNumbers compares the guessed number with the random number and
// updates the outcome accordingly.
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
