package cli_test

import (
	"bytes"
	"testing"

	c "github.com/go-number-guessing-game/internal/cli"
	"github.com/stretchr/testify/assert"
)

func TestUnitDisplayMessage(t *testing.T) {
	t.Run("display message", func(t *testing.T) {
		testCases := []struct {
			description string
			input       any
			want        string
		}{
			{
				description: "for a string",
				input:       "test",
				want:        "test",
			},
			{
				description: "for a slice of string",
				input:       []string{"test1", "test2"},
				want:        "test1test2",
			},
		}

		for _, tc := range testCases {
			t.Run(tc.description, func(t *testing.T) {
				buffer := &bytes.Buffer{}
				c.Display(buffer, tc.input)

				got := buffer.String()
				assert.Equal(t, tc.want, got)
			})
		}
	})

	t.Run("empty message return error when", func(t *testing.T) {
		testCases := []struct {
			description string
			input       any
			want        error
		}{
			{
				description: "empty string",
				input:       "",
				want:        c.NewEmptyError(c.EmptyMessage["value"]),
			},
			{
				description: "empty slice of string",
				input:       []string{},
				want:        c.NewEmptyError(c.EmptyMessage["value"]),
			},
			{
				description: "slice with one empty string",
				input:       []string{""},
				want:        c.NewEmptyError(c.EmptyMessage["value"]),
			},
			{
				description: "any type not being string or slice of string",
				input:       0,
				want:        c.NewValueTypeError(0),
			},
		}

		for _, tc := range testCases {
			t.Run(tc.description, func(t *testing.T) {
				buffer := &bytes.Buffer{}

				got := c.Display(buffer, tc.input)

				assert.NotNil(t, got)
				assert.ErrorAs(t, got, &tc.want)
			})
		}
	})
}

func TestUnitCliInput(t *testing.T) {
	t.Run("input tests", func(t *testing.T) {
		testCases := []struct {
			description string
			input       string
			want        error
			isGuess     bool
		}{
			{
				description: "valid next difficulty input",
				input:       "test",
				want:        nil,
				isGuess:     false,
			},
			{
				description: "empty nexy difficulty",
				input:       "",
				want:        c.NewEmptyError(c.EmptyMessage["empty_input"]),
				isGuess:     false,
			},
			{
				description: "valid next guess number input",
				input:       "test",
				want:        nil,
				isGuess:     true,
			},
			{
				description: "empty next guess number input",
				input:       "",
				want:        c.NewEmptyError(c.EmptyMessage["empty_input"]),
				isGuess:     true,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.description, func(t *testing.T) {
				buffer := bytes.NewBufferString(tc.input)
				cli := c.CliInput{Source: buffer}
				var got error
				var input string

				if tc.isGuess {
					input, got = cli.NextGuessNumberInput()
				} else {
					input, got = cli.NextDifficultyInput()
				}

				if tc.want == nil {
					assert.NoError(t, got)
					assert.Equal(t, tc.input, input)
				} else {
					assert.NotNil(t, got)
					assert.ErrorAs(t, got, &tc.want)
				}
			})
		}
	})
}
