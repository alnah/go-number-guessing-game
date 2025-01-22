package parser_test

import (
	"testing"

	p "github.com/go-number-guessing-game/internal/parser"
	"github.com/stretchr/testify/assert"
)

func TestUnitParsePlayerInput(t *testing.T) {
	t.Run("return parsed player", func(t *testing.T) {
		value := "test"

		got, err := p.ParsePlayerInput(value)
		assert.NoError(t, err)

		assert.Equal(t, value, got)
	})

	t.Run("return error when player is more than 30 chars", func(t *testing.T) {
		var value string
		for i := 1; i <= 31; i++ {
			value += "a"
		}

		want := p.NewParsePlayerError()
		_, got := p.ParsePlayerInput(value)
		assert.NotNil(t, got)
		assert.ErrorAs(t, got, &want)
	})
}

func TestUnitParseGuessNumberInput(t *testing.T) {
	t.Run("return parsed guess number between 1 and 100",
		func(t *testing.T) {
			testCases := []struct {
				description  string
				stringValue  string
				integerValue int
			}{
				{
					description:  "in range",
					stringValue:  "50",
					integerValue: 50,
				},
				{
					description:  "on min",
					stringValue:  "1",
					integerValue: 1,
				},
				{
					description:  "on max",
					stringValue:  "100",
					integerValue: 100,
				},
			}
			for _, tc := range testCases {
				t.Run(tc.description, func(t *testing.T) {
					got, err := p.ParseGuessNumberInput(tc.stringValue)

					assert.NoError(t, err)
					assert.Equal(t, tc.integerValue, got)
					assert.IsType(t, int(tc.integerValue), got)
				})
			}
		})

	t.Run("return error when string value not an integer", func(t *testing.T) {
		want := p.NewParseNumberError()
		_, got := p.ParseGuessNumberInput("not int")

		assert.NotNil(t, got)
		assert.ErrorAs(t, got, &want)
	})

	t.Run("return error when number is out of range", func(t *testing.T) {
		testCases := []struct {
			description string
			value       string
		}{
			{
				description: "less than 1",
				value:       "0",
			},
			{
				description: "greater than 100",
				value:       "101",
			},
		}

		for _, tc := range testCases {
			t.Run(tc.description, func(t *testing.T) {
				want := p.NewNumberRangeError(1, 100)
				_, got := p.ParseGuessNumberInput(tc.value)

				assert.NotNil(t, got)
				assert.ErrorAs(t, got, &want)
			})
		}
	})
}

func TestUnitParsePlayAgainInput(t *testing.T) {
	t.Run("return parsed play again as bool",
		func(t *testing.T) {
			testCases := []struct {
				description string
				stringValue string
				boolean     bool
			}{
				{
					description: "on min",
					stringValue: "1",
					boolean:     true,
				},
				{
					description: "on max",
					stringValue: "2",
					boolean:     false,
				},
			}
			for _, tc := range testCases {
				t.Run(tc.description, func(t *testing.T) {
					got, err := p.ParsePlayAgainInput(tc.stringValue)

					assert.NoError(t, err)
					assert.Equal(t, tc.boolean, got)
					assert.IsType(t, bool(tc.boolean), got)
				})
			}
		})

	t.Run("return error when string value not an integer", func(t *testing.T) {
		want := p.NewParseNumberError()
		_, got := p.ParsePlayAgainInput("not int")

		assert.NotNil(t, got)
		assert.ErrorAs(t, got, &want)
	})

	t.Run("return error when number is out of range", func(t *testing.T) {
		testCases := []struct {
			description string
			value       string
		}{
			{
				description: "less than 1",
				value:       "0",
			},
			{
				description: "greater than 2",
				value:       "3",
			},
		}

		for _, tc := range testCases {
			t.Run(tc.description, func(t *testing.T) {
				want := p.NewNumberRangeError(1, 100)
				_, got := p.ParsePlayAgainInput(tc.value)

				assert.NotNil(t, got)
				assert.ErrorAs(t, got, &want)
			})
		}
	})
}

func TestUnitParseDifficultyInput(t *testing.T) {
	t.Run("return parsed difficulty level and max attempt", func(t *testing.T) {
		testCases := []struct {
			description     string
			value           string
			wantLevel       string
			wantMaxAttempts int
		}{
			{
				description:     "easy level is 10 max attempts",
				value:           "1",
				wantLevel:       "Easy",
				wantMaxAttempts: 10,
			},
			{
				description:     "medium level is 5 max attempts",
				value:           "2",
				wantLevel:       "Medium",
				wantMaxAttempts: 5,
			},
			{
				description:     "hard level is 3 max attempts",
				value:           "3",
				wantLevel:       "Hard",
				wantMaxAttempts: 3,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.description, func(t *testing.T) {
				gotLevel, gotMaxAttempts, err := p.ParseDifficultyInput(tc.value)

				assert.NoError(t, err)
				assert.Equal(t, tc.wantLevel, gotLevel)
				assert.Equal(t, tc.wantMaxAttempts, gotMaxAttempts)
			})
		}
	})

	t.Run("return error when string value not an integer", func(t *testing.T) {
		want := p.NewParseNumberError()
		_, _, got := p.ParseDifficultyInput("not int")

		assert.NotNil(t, got)
		assert.ErrorAs(t, got, &want)
	})

	t.Run("return error when difficulty out of range", func(t *testing.T) {
		testCases := []struct {
			description string
			value       string
		}{
			{
				description: "less than 1",
				value:       "0",
			},
			{
				description: "greater than 3",
				value:       "4",
			},
		}

		for _, tc := range testCases {
			want := p.NewNumberRangeError(1, 3)
			_, _, got := p.ParseDifficultyInput(tc.value)

			assert.NotNil(t, got)
			assert.ErrorAs(t, got, &want)
		}
	})
}
