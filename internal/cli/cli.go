package cli

import (
	"bufio"
	_ "embed"
	"fmt"
	"io"
)

// InputSource defines an interface for obtaining user input for game difficulty
// and guess numbers.
type InputSource interface {
	NextDifficultyInput() (string, error)
	NextGuessNumberInput() (string, error)
}

// CliInput is a structure that implements the InputSource interface,
// allowing the retrieval of user input for game difficulty and guess numbers.
type CliInput struct {
	Source io.Reader // Source is the input stream from which user input is read.
}

// NextDifficultyInput retrieves the next difficulty input from the user.
// It calls GetUserInput with the Source to obtain the input string and any error encountered.
func (c *CliInput) NextDifficultyInput() (string, error) {
	return GetUserInput(c.Source)
}

// NextGuessNumberInput retrieves the next guess number input from the user.
// Similar to NextDifficultyInput, it uses GetUserInput to read from the Source.
func (c *CliInput) NextGuessNumberInput() (string, error) {
	return GetUserInput(c.Source)
}

// EmptyError represents an error when a message is empty.
type EmptyError struct {
	Message string
}

// EmptyMessage holds predefined messages for empty inputs.
var EmptyMessage = map[string]string{
	"value": "Message to display can't be empty.",
	"input": "You must enter something!",
}

// Error returns the error message for EmptyError.
func (e *EmptyError) Error() string {
	return fmt.Sprintf("%s", e.Message)
}

// NewEmptyError creates a new instance of EmptyError with a given message.
func NewEmptyError(message string) error {
	return &EmptyError{Message: message}
}

// ValueTypeError represents an error when the value type is incorrect.
type ValueTypeError struct {
	Value any
}

// TextTypeMessage is the error message for invalid value types.
const TextTypeMessage = "Value must be a string or a slice of string, got %+v"

// Error returns the error message for ValueTypeError.
func (e *ValueTypeError) Error() string {
	return fmt.Sprintf(TextTypeMessage, e.Value)
}

// NewValueTypeError creates a new instance of ValueTypeError with a given value.
func NewValueTypeError(value any) error {
	return &ValueTypeError{Value: value}
}

// Strings is a type alias for a slice of strings.
type Strings []string

// Display prints a string or a slice of strings to the provided writer.
// It validates that the input is not empty and returns an error if it is.
func Display(writer io.Writer, value any) error {
	switch t := value.(type) {
	case string:
		err := validateEmptyString(t)
		if err != nil {
			return err
		}
		fmt.Fprint(writer, t)

	case []string:
		err := validateEmptySlice(t)
		if err != nil {
			return err
		}

		for _, s := range t {
			err := validateEmptyString(s)
			if err != nil {
				return err
			}
			fmt.Fprint(writer, s)
		}

	default:
		return NewValueTypeError(value)
	}

	return nil
}

// GetUserInput reads a line of input from the provided reader.
// It returns an error if the input is empty.
func GetUserInput(reader io.Reader) (string, error) {
	scanner := bufio.NewScanner(reader)
	scanner.Scan()
	if scanner.Text() == "" {
		return "", NewEmptyError(EmptyMessage["input"])
	}

	return scanner.Text(), nil
}

// validateEmptySlice checks if a slice of strings is empty.
// It returns an error if the slice is empty.
func validateEmptySlice(slice []string) error {
	if len(slice) == 0 {
		return NewEmptyError(EmptyMessage["value"])
	}
	return nil
}

// validateEmptyString checks if a string is empty.
// It returns an error if the string is empty.
func validateEmptyString(str any) error {
	if str == "" {
		return NewEmptyError(EmptyMessage["value"])
	}
	return nil
}
