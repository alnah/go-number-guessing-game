package cli

import (
	"bufio"
	_ "embed"
	"fmt"
	"io"
)

type InputSource interface {
	NextPlayer() (string, error)
	NextDifficultyInput() (string, error)
	NextGuessNumberInput() (string, error)
	NextPlayAgainInput() (string, error)
}

type CliInput struct {
	Source io.Reader
}

func (c *CliInput) NextPlayer() (string, error) {
	return GetUserInput(c.Source)
}

func (c *CliInput) NextDifficultyInput() (string, error) {
	return GetUserInput(c.Source)
}

func (c *CliInput) NextGuessNumberInput() (string, error) {
	return GetUserInput(c.Source)
}

func (c *CliInput) NextPlayAgainInput() (string, error) {
	return GetUserInput(c.Source)
}

type EmptyError struct {
	Message string
}

var EmptyMessage = map[string]string{
	"value": "Message to display can't be empty.",
	"input": "You must enter something!",
}

func (e *EmptyError) Error() string {
	return e.Message
}

func NewEmptyError(message string) error {
	return &EmptyError{Message: message}
}

type ValueTypeError struct {
	Value any
}

const TextTypeMessage = "Value must be a string or a slice of string, got %+v"

func (e *ValueTypeError) Error() string {
	return fmt.Sprintf(TextTypeMessage, e.Value)
}

func NewValueTypeError(value any) error {
	return &ValueTypeError{Value: value}
}

type Strings []string

func Display(writer io.Writer, value any) {
	switch t := value.(type) {
	case string:
		err := validateEmptyString(t)
		if err != nil {
			fmt.Fprint(writer, EmptyMessage["value"])
		}
		fmt.Fprint(writer, t)

	case []string:
		err := validateEmptySlice(t)
		if err != nil {
			fmt.Fprint(writer, EmptyMessage["value"])
		}

		for _, s := range t {
			err := validateEmptyString(s)
			if err != nil {
				fmt.Fprint(writer, EmptyMessage["value"])
			}
			fmt.Fprint(writer, s)
		}

	default:
		fmt.Fprint(writer, EmptyMessage["value"])
	}
}

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
