package parser

import (
	"fmt"
	"strconv"
)

// ParsePlayerError represents an error that occurs when the player input
// exceeds the allowed character limit.
type ParsePlayerError struct{}

// ParsePlayerMessage is the error message indicating that the player name
// must be 20 characters or fewer.
const ParsePlayerMessage = "It must be 20 characters at most."

// Error returns the error message for ParsePlayerError.
func (e *ParsePlayerError) Error() string {
	return fmt.Sprint(ParsePlayerMessage)
}

// NewParsePlayerError creates a new instance of ParsePlayerError.
func NewParsePlayerError() error {
	return &ParsePlayerError{}
}

// ParsePlayerInput validates the player input string and ensures it does not
// exceed the maximum allowed length of 20 characters.
func ParsePlayerInput(s string) (string, error) {
	if len(s) > 20 {
		return "", NewParsePlayerError()
	}
	return s, nil
}

// ParseNumberError represents an error that occurs when parsing a number
// fails due to an invalid format.
type ParseNumberError struct{}

// ParseNumberMessage is the error message indicating that the input must
// be an integer.
const ParseNumberMessage = "It must be an integer."

// Error returns the error message for ParseNumberError.
func (e *ParseNumberError) Error() string {
	return fmt.Sprint(ParseNumberMessage)
}

// NewParseNumberError creates a new instance of ParseNumberError.
func NewParseNumberError() error {
	return &ParseNumberError{}
}

// NumberRangeError represents an error that occurs when a number is
// outside the specified range.
type NumberRangeError struct {
	Min int // Minimum valid value
	Max int // Maximum valid value
}

// NumberRangeMessage is the error message indicating that the input must
// be within a specific range.
const NumberRangeMessage = "It must be an integer between %d and %d."

// Error returns the error message for NumberRangeError.
func (e *NumberRangeError) Error() string {
	return fmt.Sprintf(NumberRangeMessage, e.Min, e.Max)
}

// NewNumberRangeError creates a new instance of NumberRangeError with
// the specified minimum and maximum values.
func NewNumberRangeError(min, max int) error {
	return &NumberRangeError{Min: min, Max: max}
}

// ParseGuessNumberInput parses the input string as a guess number and
// validates that it is within the range of 1 to 100.
func ParseGuessNumberInput(s string) (int, error) {
	return validateInputNumber(s, 1, 100)
}

// ParsePlayAgainInput parses the input string as a play again value and
// validates that is is within the range of 1 to 2
func ParsePlayAgainInput(s string) (bool, error) {
	integer, err := validateInputNumber(s, 1, 2)
	if err != nil {
		return false, err
	}

	mapPlayAgain := map[int]bool{
		1: true,
		2: false,
	}

	return mapPlayAgain[integer], nil
}

// DifficultyToMaxAttempts is a mapping of difficulty levels to their
// corresponding maximum attempts.
type DifficultyToMaxAttempts map[string]int

// ParseDifficultyInput parses the input string as a difficulty level and
// returns the corresponding level name and maximum attempts.
func ParseDifficultyInput(s string) (string, int, error) {
	integer, err := validateInputNumber(s, 1, 3)
	if err != nil {
		return "", 0, err
	}

	toDifficulty := map[int]string{
		1: "Easy",
		2: "Medium",
		3: "Hard",
	}

	toMaxAttempts := map[int]int{
		1: 10,
		2: 5,
		3: 3,
	}

	return toDifficulty[integer], toMaxAttempts[integer], nil
}

// validateInputNumber validates that the input string can be converted
// to an integer and is within the specified range.
func validateInputNumber(stringValue string, min, max int) (int, error) {
	integer, err := strconv.Atoi(stringValue)
	if err != nil {
		return 0, NewParseNumberError()
	}

	if integer < min || integer > max {
		return 0, NewNumberRangeError(min, max)
	}
	return integer, nil
}
