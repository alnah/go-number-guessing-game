// Package cli provides utilities for handling user input in the command-line
// interface of the game. It defines input sources and error types for input
// validation.
package cli

import (
	"bufio"
	_ "embed"
	"fmt"
	"io"
)

// InputSource defines an interface for obtaining various types of user input.
// This abstraction allows for easy mocking in service integration tests,
// where a mock input source simulates user interactions.
type InputSource interface {
	NextPlayerInput() (string, error)
	NextDifficultyInput() (string, error)
	NextGuessNumberInput() (string, error)
	NextPlayAgainInput() (string, error)
}

// CliInput implements the InputSource interface, providing methods to read
// user input from a specified io.Reader source.
type CliInput struct {
	Source io.Reader
}

// NextPlayerInput retrieves the next player input from the source.
func (c *CliInput) NextPlayerInput() (string, error) {
	return GetUserInput(c.Source)
}

// NextDifficultyInput retrieves the next difficulty input from the source.
func (c *CliInput) NextDifficultyInput() (string, error) {
	return GetUserInput(c.Source)
}

// NextGuessNumberInput retrieves the next guess number input from the source.
func (c *CliInput) NextGuessNumberInput() (string, error) {
	return GetUserInput(c.Source)
}

// NextPlayAgainInput retrieves the next play-again input from the source.
func (c *CliInput) NextPlayAgainInput() (string, error) {
	return GetUserInput(c.Source)
}

// EmptyError represents an error that occurs when an expected input is empty.
type EmptyError struct {
	Message string
}

// EmptyMessage provides default messages for empty input errors.
var EmptyMessage = map[string]string{
	"value": "Message to display can't be empty.",
	"input": "You must enter something!",
}

// Error returns the error message for EmptyError.
func (e *EmptyError) Error() string {
	return e.Message
}

// NewEmptyError creates a new instance of EmptyError with the specified message
// for testing.
func NewEmptyError(message string) error {
	return &EmptyError{Message: message}
}

// ValueTypeError indicates an error when an unexpected value type is encountered.
type ValueTypeError struct {
	Value any
}

// TextTypeMessage is the message displayed when the value type is incorrect.
const TextTypeMessage = "Value must be a string or a slice of string, got %+v"

// Error returns the error message for ValueTypeError.
func (e *ValueTypeError) Error() string {
	return fmt.Sprintf(TextTypeMessage, e.Value)
}

// NewValueTypeError creates a new instance of ValueTypeError with the specified
// value for testing.
func NewValueTypeError(value any) error {
	return &ValueTypeError{Value: value}
}

// Strings is a type alias for a slice of strings.
type Strings []string

// Display outputs a string or a slice of strings to the specified writer,
// handling empty values appropriately.
func Display(writer io.Writer, value any) {
	switch t := value.(type) {
	case string:
		if err := validateEmptyString(t); err != nil {
			fmt.Fprint(writer, EmptyMessage["value"])
		}
		fmt.Fprint(writer, t)

	case []string:
		if err := validateEmptySlice(t); err != nil {
			fmt.Fprint(writer, EmptyMessage["value"])
		}

		for _, s := range t {
			if err := validateEmptyString(s); err != nil {
				fmt.Fprint(writer, EmptyMessage["value"])
			}
			fmt.Fprint(writer, s)
		}

	default:
		fmt.Fprint(writer, EmptyMessage["value"])
	}
}

// GetUserInput reads a line of input from the specified reader, returning an
// error if the input is empty.
func GetUserInput(reader io.Reader) (string, error) {
	scanner := bufio.NewScanner(reader)
	scanner.Scan()
	if scanner.Text() == "" {
		return "", NewEmptyError(EmptyMessage["input"])
	}

	return scanner.Text(), nil
}

func validateEmptySlice(slice []string) error {
	if len(slice) == 0 {
		return NewEmptyError(EmptyMessage["value"])
	}

	return nil
}

func validateEmptyString(str any) error {
	if str == "" {
		return NewEmptyError(EmptyMessage["value"])
	}

	return nil
}
