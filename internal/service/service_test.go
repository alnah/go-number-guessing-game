package service_test

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/go-number-guessing-game/internal/config"
	"github.com/go-number-guessing-game/internal/number"
	s "github.com/go-number-guessing-game/internal/service"
	"github.com/stretchr/testify/assert"
)

const newline string = "\n"

var gameConfig = config.LoadConfig("yaml", "../../configs/data.yaml")

func TestIntegrationGameConfig(t *testing.T) {
	var wantSet = map[string]struct{}{
		"greeting":     {},
		"difficulty":   {},
		"level":        {},
		"guess":        {},
		"greater":      {},
		"less":         {},
		"equal":        {},
		"max_attempts": {},
		"newline":      {},
		"spacer":       {},
	}

	var gotSet = make(map[string]struct{})
	for k := range gameConfig {
		gotSet[k] = struct{}{}
	}
	assert.Equal(t, wantSet, gotSet)
}

func TestIntegrationPlay(t *testing.T) {
	randomNumber := 50

	t.Run("user success", func(t *testing.T) {
		testCases := []struct {
			description     string
			mockInputSource *MockInputSource
			outputStrings   []string
		}{
			{
				description: "hard difficulty",
				mockInputSource: &MockInputSource{
					DifficultyInput:   []string{"3"},
					GuessNumberInputs: []string{"51", "49", "50"},
				},
				outputStrings: []string{
					gameConfig["greeting"],
					gameConfig["newline"],
					gameConfig["difficulty"],
					fmt.Sprintf(gameConfig["level"], "Hard"),
					gameConfig["newline"],
					gameConfig["guess"],
					fmt.Sprintf(gameConfig["less"], 51),
					gameConfig["spacer"],
					gameConfig["guess"],
					fmt.Sprintf(gameConfig["greater"], 49),
					gameConfig["spacer"],
					gameConfig["guess"],
					fmt.Sprintf(gameConfig["equal"], 3),
					gameConfig["newline"],
				},
			},
			{
				description: "medium difficulty",
				mockInputSource: &MockInputSource{
					DifficultyInput:   []string{"2"},
					GuessNumberInputs: []string{"52", "51", "48", "49", "50"},
				},
				outputStrings: []string{
					gameConfig["greeting"],
					gameConfig["newline"],
					gameConfig["difficulty"],
					fmt.Sprintf(gameConfig["level"], "Medium"),
					gameConfig["newline"],
					gameConfig["guess"],
					fmt.Sprintf(gameConfig["less"], 52),
					gameConfig["spacer"],
					gameConfig["guess"],
					fmt.Sprintf(gameConfig["less"], 51),
					gameConfig["spacer"],
					gameConfig["guess"],
					fmt.Sprintf(gameConfig["greater"], 48),
					gameConfig["spacer"],
					gameConfig["guess"],
					fmt.Sprintf(gameConfig["greater"], 49),
					gameConfig["spacer"],
					gameConfig["guess"],
					fmt.Sprintf(gameConfig["equal"], 5),
					gameConfig["newline"],
				},
			},
			{
				description: "easy difficulty",
				mockInputSource: &MockInputSource{
					DifficultyInput: []string{"1"},
					GuessNumberInputs: []string{
						"55", "54", "53", "52", "51",
						"46", "47", "48", "49", "50",
					},
				},
				outputStrings: []string{
					gameConfig["greeting"],
					gameConfig["newline"],
					gameConfig["difficulty"],
					fmt.Sprintf(gameConfig["level"], "Easy"),
					gameConfig["newline"],
					gameConfig["guess"],
					fmt.Sprintf(gameConfig["less"], 55),
					gameConfig["spacer"],
					gameConfig["guess"],
					fmt.Sprintf(gameConfig["less"], 54),
					gameConfig["spacer"],
					gameConfig["guess"],
					fmt.Sprintf(gameConfig["less"], 53),
					gameConfig["spacer"],
					gameConfig["guess"],
					fmt.Sprintf(gameConfig["less"], 52),
					gameConfig["spacer"],
					gameConfig["guess"],
					fmt.Sprintf(gameConfig["less"], 51),
					gameConfig["spacer"],
					gameConfig["guess"],
					fmt.Sprintf(gameConfig["greater"], 46),
					gameConfig["spacer"],
					gameConfig["guess"],
					fmt.Sprintf(gameConfig["greater"], 47),
					gameConfig["spacer"],
					gameConfig["guess"],
					fmt.Sprintf(gameConfig["greater"], 48),
					gameConfig["spacer"],
					gameConfig["guess"],
					fmt.Sprintf(gameConfig["greater"], 49),
					gameConfig["spacer"],
					gameConfig["guess"],
					fmt.Sprintf(gameConfig["equal"], 10),
					gameConfig["newline"],
				},
			},
		}

		for _, tc := range testCases {
			t.Run(tc.description, func(t *testing.T) {
				gotWriter, game := initGame(tc.mockInputSource)
				var wantWriter strings.Builder
				for _, s := range tc.outputStrings {
					wantWriter.WriteString(s)
				}
				game.Play(randomNumber)

				assert.Equal(t, wantWriter.String(), gotWriter.String())
			})
		}
	})

	t.Run("no more attempts", func(t *testing.T) {
		mockInputSource := &MockInputSource{
			DifficultyInput:   []string{"3"},
			GuessNumberInputs: []string{"53", "52", "51"},
		}
		gotWriter, game := initGame(mockInputSource)
		game.Play(randomNumber)

		assert.Contains(t, gotWriter.String(), gameConfig["max_attempts"])
	})

	t.Run("invalid difficulty inputs", func(t *testing.T) {
		testCases := []struct {
			description      string
			invalidInput     string
			wantErrorMessage string
		}{
			{
				description:      "incorrect",
				invalidInput:     "incorrect",
				wantErrorMessage: number.ParseNumberMessage,
			},
			{
				description:      "greater than 3",
				invalidInput:     "4",
				wantErrorMessage: fmt.Sprintf(number.NumberRangeMessage, 1, 3),
			},
			{
				description:      "less than 1",
				invalidInput:     "0",
				wantErrorMessage: fmt.Sprintf(number.NumberRangeMessage, 1, 3),
			},
		}

		for _, tc := range testCases {
			t.Run(tc.description, func(t *testing.T) {
				mockInputSource := &MockInputSource{
					DifficultyInput:   []string{tc.invalidInput, "3"},
					GuessNumberInputs: []string{"50"},
				}
				gotWriter, game := initGame(mockInputSource)
				game.Play(randomNumber)

				gotOutput := gotWriter.String()

				assert.Contains(t, gotOutput, tc.wantErrorMessage)
				assert.Contains(t, gotOutput, gameConfig["difficulty"])
				assert.Contains(t, gotOutput, "Hard")
				assert.Contains(t, gotOutput, gameConfig["guess"])
				assert.Contains(t, gotOutput, fmt.Sprintf(gameConfig["equal"], 1))
			})
		}
	})

	t.Run("invalid guess number inputs", func(t *testing.T) {
		testCases := []struct {
			description      string
			invalidInput     string
			wantErrorMessage string
		}{
			{
				description:      "incorrect",
				invalidInput:     "incorrect",
				wantErrorMessage: number.ParseNumberMessage,
			},
			{
				description:      "less than 1",
				invalidInput:     "0",
				wantErrorMessage: fmt.Sprintf(number.NumberRangeMessage, 1, 100),
			},
			{
				description:      "greater than 100",
				invalidInput:     "101",
				wantErrorMessage: fmt.Sprintf(number.NumberRangeMessage, 1, 100),
			},
		}

		for _, tc := range testCases {
			t.Run(tc.description, func(t *testing.T) {
				mockInputSource := &MockInputSource{
					DifficultyInput:   []string{"3"},
					GuessNumberInputs: []string{tc.invalidInput, "50"},
				}
				gotWriter, game := initGame(mockInputSource)
				game.Play(randomNumber)

				gotOutput := gotWriter.String()

				assert.Contains(t, gotOutput, "Hard")
				assert.Contains(t, gotOutput, tc.wantErrorMessage)
				assert.Contains(t, gotOutput, gameConfig["guess"])
				assert.Contains(t, gotOutput, fmt.Sprintf(gameConfig["equal"], 1))
			})
		}
	})
}

func initGame(mockInputSource *MockInputSource) (*bytes.Buffer, s.Game) {
	gotWriter := &bytes.Buffer{}
	game := s.Game{
		InputSource: mockInputSource,
		GameConfig:  gameConfig,
		Writer:      gotWriter,
	}
	return gotWriter, game
}

type MockInputSource struct {
	DifficultyInput   []string
	difficultyIndex   int
	GuessNumberInputs []string
	guessNumberIndex  int
}

func (m *MockInputSource) NextDifficultyInput() (string, error) {
	if m.difficultyIndex < len(m.DifficultyInput) {
		value := m.DifficultyInput[m.difficultyIndex]
		m.difficultyIndex++
		return value, nil
	}
	return "", fmt.Errorf("no more difficulty inputs")
}

func (m *MockInputSource) NextGuessNumberInput() (string, error) {
	if m.guessNumberIndex < len(m.GuessNumberInputs) {
		value := m.GuessNumberInputs[m.guessNumberIndex]
		m.guessNumberIndex++
		return value, nil
	}

	return "", fmt.Errorf("no more guess number inputs")
}
