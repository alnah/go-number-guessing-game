package parser

import (
	"fmt"
	"strconv"
)

type ParsePlayerError struct{}

const ParsePlayerMessage = "It must be 20 characters at most."

func (e *ParsePlayerError) Error() string {
	return fmt.Sprint(ParsePlayerMessage)
}

func NewParsePlayerError() error {
	return &ParsePlayerError{}
}

func ParsePlayerInput(s string) (string, error) {
	if len(s) > 20 {
		return "", NewParsePlayerError()
	}
	return s, nil
}

type ParseNumberError struct{}

const ParseNumberMessage = "It must be an integer."

func (e *ParseNumberError) Error() string {
	return fmt.Sprint(ParseNumberMessage)
}

func NewParseNumberError() error {
	return &ParseNumberError{}
}

type NumberRangeError struct {
	Min int
	Max int
}

const NumberRangeMessage = "It must be an integer between %d and %d."

func (e *NumberRangeError) Error() string {
	return fmt.Sprintf(NumberRangeMessage, e.Min, e.Max)
}

func NewNumberRangeError(min, max int) error {
	return &NumberRangeError{Min: min, Max: max}
}

func ParseGuessNumberInput(s string) (int, error) {
	return validateInputNumber(s, 1, 100)
}

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

type DifficultyToMaxAttempts map[string]int

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
