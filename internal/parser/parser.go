// Package parser provides utilities for parsing user inputs, ensuring
// validation and error handling for player names, numbers, and game
// difficulty levels.
package parser

import (
	"fmt"
	"strconv"
)

// ParsePlayerError indicates an error when parsing player name input.
type ParsePlayerError struct{}

// ParsePlayerMessage is the message displayed when the player name exceeds
// 20 characters. It is public for testing purposes.
const ParsePlayerMessage = "It must be non-empty, and 20 characters at most."

// Error returns the error message for ParsePlayerError.
func (e *ParsePlayerError) Error() string {
	return fmt.Sprint(ParsePlayerMessage)
}

// NewParsePlayerError creates a new instance of ParsePlayerError for testing.
func NewParsePlayerError() error {
	return &ParsePlayerError{}
}

// ParsePlayerInput validates the player name, ensuring it is not empty and
// does not exceed 20 characters. Returns a custom error if validation fails.
func ParsePlayerInput(s string) (string, error) {
	if len(s) > 20 || s == "" {
		return "", NewParsePlayerError()
	}
	return s, nil
}

// ParseNumberError indicates an error when parsing number input.
type ParseNumberError struct{}

// ParseNumberMessage is the message displayed when the input is not an integer.
// It is public for testing purposes.
const ParseNumberMessage = "It must be an integer."

// Error returns the error message for ParseNumberError.
func (e *ParseNumberError) Error() string {
	return fmt.Sprint(ParseNumberMessage)
}

// NewParseNumberError creates a new instance of ParseNumberError for testing.
func NewParseNumberError() error {
	return &ParseNumberError{}
}

// NumberRangeError indicates an error when a number is out of the range.
type NumberRangeError struct {
	Min int
	Max int
}

// NumberRangeMessage is the message displayed when the number is out of range.
const NumberRangeMessage = "It must be an integer between %d and %d."

// Error returns the error message for NumberRangeError.
func (e *NumberRangeError) Error() string {
	return fmt.Sprintf(NumberRangeMessage, e.Min, e.Max)
}

// NewNumberRangeError creates a new instance of NumberRangeError for testing.
func NewNumberRangeError(min, max int) error {
	return &NumberRangeError{Min: min, Max: max}
}

// ParseGuessNumberInput validates and returns the parsed guess number,
// ensuring it is between 1 and 100. Returns a custom error if validation fails.
func ParseGuessNumberInput(s string) (int, error) {
	return validateInputNumber(s, 1, 100)
}

// ParsePlayAgainInput validates and returns the parsed play again input,
// ensuring the user input is either 1 or 2. Returns a custom error if
// validation fails.
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

// ParseDifficultyInput validates and returns the parsed difficulty level
// and maximum number of attempts based on user input, ensuring the input
// is between 1 and 3. Returns a custom error if validation fails.
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

// validateInputNumber converts a string to an integer and checks if it is
// within the specified range. Returns a custom error if validation fails.
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
