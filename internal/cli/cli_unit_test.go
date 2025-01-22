package cli_test

import (
	"bytes"
	"testing"

	"github.com/go-number-guessing-game/internal/cli"
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
				cli.Display(buffer, tc.input)
				got := buffer.String()
				assert.Equal(t, tc.want, got)
			})
		}
	})

	t.Run("empty message display value can't be empty", func(t *testing.T) {
		testCases := []struct {
			description string
			input       any
			want        string
		}{
			{
				description: "empty string",
				input:       "",
				want:        cli.EmptyMessage["value"],
			},
			{
				description: "empty slice of string",
				input:       []string{},
				want:        cli.EmptyMessage["value"],
			},
			{
				description: "slice with one empty string",
				input:       []string{""},
				want:        cli.EmptyMessage["value"],
			},
			{
				description: "any type not being string or slice of string",
				input:       0,
				want:        cli.EmptyMessage["value"],
			},
		}

		for _, tc := range testCases {
			t.Run(tc.description, func(t *testing.T) {
				buffer := &bytes.Buffer{}

				cli.Display(buffer, tc.input)
				got := buffer.String()

				assert.Equal(t, tc.want, got)
			})
		}
	})
}

func TestUnitCliInput(t *testing.T) {
	t.Run("input tests", func(t *testing.T) {
		testCases := []struct {
			description   string
			input         string
			want          error
			isPlayer      bool
			isDifficulty  bool
			isGuessNumber bool
			isPlayAgain   bool
		}{
			{
				description: "valid next player input",
				input:       "test",
				want:        nil,
				isPlayer:    true,
			},
			{
				description: "empty next player input",
				input:       "",
				want:        cli.NewEmptyError(cli.EmptyMessage["empty_input"]),
				isPlayer:    true,
			},
			{
				description:  "valid next difficulty input",
				input:        "1",
				want:         nil,
				isDifficulty: true,
			},
			{
				description:  "empty next difficulty input",
				input:        "",
				want:         cli.NewEmptyError(cli.EmptyMessage["empty_input"]),
				isDifficulty: true,
			},
			{
				description:   "valid next guess number input",
				input:         "1",
				want:          nil,
				isGuessNumber: true,
			},
			{
				description:   "empty next guess number input",
				input:         "",
				want:          cli.NewEmptyError(cli.EmptyMessage["empty_input"]),
				isGuessNumber: true,
			},
			{
				description: "valid play again input",
				input:       "1",
				want:        nil,
				isPlayAgain: true,
			},
			{
				description: "empty play again input",
				input:       "",
				want:        cli.NewEmptyError(cli.EmptyMessage["input"]),
				isPlayAgain: true,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.description, func(t *testing.T) {
				buffer := bytes.NewBufferString(tc.input)
				cli := cli.CliInput{Source: buffer}

				var got error
				var input string
				switch {
				case tc.isPlayer:
					input, got = cli.NextPlayerInput()
				case tc.isPlayAgain:
					input, got = cli.NextPlayAgainInput()
				case tc.isGuessNumber:
					input, got = cli.NextGuessNumberInput()
				case tc.isDifficulty:
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
