package service_test

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/go-number-guessing-game/internal/config"
	"github.com/go-number-guessing-game/internal/parser"
	s "github.com/go-number-guessing-game/internal/service"
	"github.com/go-number-guessing-game/internal/store"
	"github.com/stretchr/testify/assert"
)

var (
	gameConfig       = config.LoadConfig("yaml", "../../configs/appconfig.yaml")
	stubScoreStore   = &StubScoreStore{isEmpty: false}
	fakeRandomNumber = 50
	fakeScores       = stubScoreStore.Load().String()
)

func TestIntegrationGameConfig(t *testing.T) {
	wantSet := map[string]struct{}{
		"greeting":     {},
		"player":       {},
		"difficulty":   {},
		"level":        {},
		"guess":        {},
		"greater":      {},
		"less":         {},
		"equal":        {},
		"max_attempts": {},
		"very_close_1": {},
		"very_close_2": {},
		"very_close_3": {},
		"close_1":      {},
		"close_2":      {},
		"far":          {},
		"very_far":     {},
		"again":        {},
		"bye":          {},
		"newline":      {},
		"spacer":       {},
	}

	gotSet := make(map[string]struct{})
	for k := range gameConfig {
		gotSet[k] = struct{}{}
	}
	assert.Equal(t, wantSet, gotSet)
}

func TestIntegrationPlay(t *testing.T) {
	t.Run("user success", func(t *testing.T) {
		testCases := []struct {
			description     string
			mockInputSource *MockInputSource
			outputStrings   []string
		}{
			{
				description: "hard difficulty",
				mockInputSource: &MockInputSource{
					PlayerInput:       []string{"test"},
					DifficultyInput:   []string{"3"},
					GuessNumberInputs: []string{"48", "51", "50"},
					PlayAgainInput:    []string{"2"},
				},
				outputStrings: []string{
					gameConfig["greeting"],
					gameConfig["spacer"],
					gameConfig["player"],
					gameConfig["difficulty"],
					fmt.Sprintf(gameConfig["level"], "Hard"),
					gameConfig["spacer"],
					gameConfig["guess"],
					fmt.Sprintf(gameConfig["greater"], 48),
					gameConfig["newline"],
					gameConfig["very_close_2"],
					gameConfig["spacer"],
					gameConfig["guess"],
					fmt.Sprintf(gameConfig["less"], 51),
					gameConfig["newline"],
					gameConfig["very_close_1"],
					gameConfig["spacer"],
					gameConfig["guess"],
					fmt.Sprintf(gameConfig["equal"], "0s", 3),
					gameConfig["newline"],
					gameConfig["spacer"],
					fakeScores,
					gameConfig["spacer"],
					gameConfig["again"],
					gameConfig["bye"],
					gameConfig["newline"],
				},
			},
			{
				description: "medium difficulty",
				mockInputSource: &MockInputSource{
					PlayerInput:       []string{"test"},
					DifficultyInput:   []string{"2"},
					GuessNumberInputs: []string{"46", "53", "48", "51", "50"},
					PlayAgainInput:    []string{"2"},
				},
				outputStrings: []string{
					gameConfig["greeting"],
					gameConfig["spacer"],
					gameConfig["player"],
					gameConfig["difficulty"],
					fmt.Sprintf(gameConfig["level"], "Medium"),
					gameConfig["spacer"],
					gameConfig["guess"],
					fmt.Sprintf(gameConfig["greater"], 46),
					gameConfig["newline"],
					gameConfig["close_1"],
					gameConfig["spacer"],
					gameConfig["guess"],
					fmt.Sprintf(gameConfig["less"], 53),
					gameConfig["newline"],
					gameConfig["very_close_3"],
					gameConfig["spacer"],
					gameConfig["guess"],
					fmt.Sprintf(gameConfig["greater"], 48),
					gameConfig["newline"],
					gameConfig["very_close_2"],
					gameConfig["spacer"],
					gameConfig["guess"],
					fmt.Sprintf(gameConfig["less"], 51),
					gameConfig["newline"],
					gameConfig["very_close_1"],
					gameConfig["spacer"],
					gameConfig["guess"],
					fmt.Sprintf(gameConfig["equal"], "0s", 5),
					gameConfig["newline"],
					gameConfig["spacer"],
					fakeScores,
					gameConfig["spacer"],
					gameConfig["again"],
					gameConfig["bye"],
					gameConfig["newline"],
				},
			},
			{
				description: "easy difficulty",
				mockInputSource: &MockInputSource{
					PlayerInput:     []string{"test"},
					DifficultyInput: []string{"1"},
					GuessNumberInputs: []string{
						"61", "39", "56", "44", "54",
						"45", "47", "48", "49", "50",
					},
					PlayAgainInput: []string{"2"},
				},
				outputStrings: []string{
					gameConfig["greeting"],
					gameConfig["spacer"],
					gameConfig["player"],
					gameConfig["difficulty"],
					fmt.Sprintf(gameConfig["level"], "Easy"),
					gameConfig["spacer"],
					gameConfig["guess"],
					fmt.Sprintf(gameConfig["less"], 61),
					gameConfig["newline"],
					gameConfig["very_far"],
					gameConfig["spacer"],
					gameConfig["guess"],
					fmt.Sprintf(gameConfig["greater"], 39),
					gameConfig["newline"],
					gameConfig["very_far"],
					gameConfig["spacer"],
					gameConfig["guess"],
					fmt.Sprintf(gameConfig["less"], 56),
					gameConfig["newline"],
					gameConfig["far"],
					gameConfig["spacer"],
					gameConfig["guess"],
					fmt.Sprintf(gameConfig["greater"], 44),
					gameConfig["newline"],
					gameConfig["far"],
					gameConfig["spacer"],
					gameConfig["guess"],
					fmt.Sprintf(gameConfig["less"], 54),
					gameConfig["newline"],
					gameConfig["close_1"],
					gameConfig["spacer"],
					gameConfig["guess"],
					fmt.Sprintf(gameConfig["greater"], 45),
					gameConfig["newline"],
					gameConfig["close_2"],
					gameConfig["spacer"],
					gameConfig["guess"],
					fmt.Sprintf(gameConfig["greater"], 47),
					gameConfig["newline"],
					gameConfig["very_close_3"],
					gameConfig["spacer"],
					gameConfig["guess"],
					fmt.Sprintf(gameConfig["greater"], 48),
					gameConfig["newline"],
					gameConfig["very_close_2"],
					gameConfig["spacer"],
					gameConfig["guess"],
					fmt.Sprintf(gameConfig["greater"], 49),
					gameConfig["newline"],
					gameConfig["very_close_1"],
					gameConfig["spacer"],
					gameConfig["guess"],
					fmt.Sprintf(gameConfig["equal"], "0s", 10),
					gameConfig["newline"],
					gameConfig["spacer"],
					fakeScores,
					gameConfig["spacer"],
					gameConfig["again"],
					gameConfig["bye"],
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
				game.Play(fakeRandomNumber, stubScoreStore)

				assert.Contains(t, wantWriter.String(), gotWriter.String())
			})
		}
	})

	t.Run("play again", func(t *testing.T) {
		testCases := []struct {
			description    string
			playAgain      bool
			playAgainInput []string
		}{
			{
				description:    "play again",
				playAgain:      true,
				playAgainInput: []string{"1", "2"},
			},
			{
				description:    "don't play again",
				playAgain:      false,
				playAgainInput: []string{"2"},
			},
		}

		for _, tc := range testCases {
			t.Run(tc.description, func(t *testing.T) {
				mockInputSource := &MockInputSource{
					PlayerInput:     []string{"test1", "test2"},
					DifficultyInput: []string{"3", "3"},
					GuessNumberInputs: []string{
						"51", "49", "50",
						"51", "49", "50",
					},
					PlayAgainInput: tc.playAgainInput,
				}

				var wantWriter strings.Builder

				for _, s := range []string{
					gameConfig["greeting"],
					gameConfig["spacer"],
					gameConfig["player"],
					gameConfig["difficulty"],
					fmt.Sprintf(gameConfig["level"], "Hard"),
					gameConfig["spacer"],
					gameConfig["guess"],
					fmt.Sprintf(gameConfig["less"], 51),
					gameConfig["newline"],
					gameConfig["very_close_1"],
					gameConfig["spacer"],
					gameConfig["guess"],
					fmt.Sprintf(gameConfig["greater"], 49),
					gameConfig["newline"],
					gameConfig["very_close_1"],
					gameConfig["spacer"],
					gameConfig["guess"],
					fmt.Sprintf(gameConfig["equal"], "0s", 3),
					gameConfig["newline"],
					gameConfig["spacer"],
					fakeScores,
					gameConfig["spacer"],
					gameConfig["again"],
				} {
					wantWriter.WriteString(s)
				}

				if tc.playAgain {
					for _, s := range []string{
						gameConfig["player"],
						gameConfig["difficulty"],
						fmt.Sprintf(gameConfig["level"], "Hard"),
						gameConfig["spacer"],
						gameConfig["guess"],
						fmt.Sprintf(gameConfig["less"], 51),
						gameConfig["newline"],
						gameConfig["very_close_1"],
						gameConfig["spacer"],
						gameConfig["guess"],
						fmt.Sprintf(gameConfig["greater"], 49),
						gameConfig["newline"],
						gameConfig["very_close_1"],
						gameConfig["spacer"],
						gameConfig["guess"],
						fmt.Sprintf(gameConfig["equal"], "0s", 3),
						gameConfig["newline"],
						gameConfig["spacer"],
						fakeScores,
						gameConfig["spacer"],
						gameConfig["again"],
					} {
						wantWriter.WriteString(s)
					}
				}

				wantWriter.WriteString(gameConfig["bye"])
				wantWriter.WriteString(gameConfig["newline"])

				gotWriter, game := initGame(mockInputSource)
				game.Play(fakeRandomNumber, stubScoreStore)

				assert.Equal(t, wantWriter.String(), gotWriter.String())
			})
		}
	})

	t.Run("no more attempts", func(t *testing.T) {
		mockInputSource := &MockInputSource{
			PlayerInput:       []string{"test"},
			DifficultyInput:   []string{"3"},
			GuessNumberInputs: []string{"53", "52", "51"},
			PlayAgainInput:    []string{"2"},
		}
		gotWriter, game := initGame(mockInputSource)
		game.Play(fakeRandomNumber, stubScoreStore)

		assert.Contains(t, gotWriter.String(), gameConfig["max_attempts"])
		assert.Contains(t, gotWriter.String(), gameConfig["bye"])
		assert.Contains(t, gotWriter.String(), gameConfig["newline"])
	})

	t.Run("invalid difficulty inputs", func(t *testing.T) {
		testCases := []struct {
			description      string
			invalidInput     string
			wantErrorMessage string
		}{
			{
				description:      "non-numeric input",
				invalidInput:     "incorrect",
				wantErrorMessage: parser.ParseNumberMessage,
			},
			{
				description:      "greater than max difficulty",
				invalidInput:     "4",
				wantErrorMessage: fmt.Sprintf(parser.NumberRangeMessage, 1, 3),
			},
			{
				description:      "less than min difficulty",
				invalidInput:     "0",
				wantErrorMessage: fmt.Sprintf(parser.NumberRangeMessage, 1, 3),
			},
		}

		for _, tc := range testCases {
			t.Run(tc.description, func(t *testing.T) {
				mockInputSource := &MockInputSource{
					PlayerInput:       []string{"test"},
					DifficultyInput:   []string{tc.invalidInput, "3"},
					GuessNumberInputs: []string{"50"},
					PlayAgainInput:    []string{"2"},
				}
				gotWriter, game := initGame(mockInputSource)
				game.Play(fakeRandomNumber, stubScoreStore)

				gotOutput := gotWriter.String()

				assert.Contains(t, gotOutput, tc.wantErrorMessage)
				assert.Contains(t, gotOutput, gameConfig["difficulty"])
				assert.Contains(t, gotOutput, fmt.Sprintf(gameConfig["level"], "Hard"))
				assert.Contains(t, gotOutput, gameConfig["guess"])
				assert.Contains(t, gotOutput, fmt.Sprintf(gameConfig["equal"], "0s", 1))
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
				description:      "non-numeric input",
				invalidInput:     "incorrect",
				wantErrorMessage: parser.ParseNumberMessage,
			},
			{
				description:      "less than min range",
				invalidInput:     "0",
				wantErrorMessage: fmt.Sprintf(parser.NumberRangeMessage, 1, 100),
			},
			{
				description:      "greater than max range",
				invalidInput:     "101",
				wantErrorMessage: fmt.Sprintf(parser.NumberRangeMessage, 1, 100),
			},
		}

		for _, tc := range testCases {
			t.Run(tc.description, func(t *testing.T) {
				mockInputSource := &MockInputSource{
					PlayerInput:       []string{"test"},
					DifficultyInput:   []string{"3"},
					GuessNumberInputs: []string{tc.invalidInput, "50"},
					PlayAgainInput:    []string{"2"},
				}
				gotWriter, game := initGame(mockInputSource)
				game.Play(fakeRandomNumber, stubScoreStore)

				gotOutput := gotWriter.String()

				assert.Contains(t, gotOutput, fmt.Sprintf(gameConfig["level"], "Hard"))
				assert.Contains(t, gotOutput, tc.wantErrorMessage)
				assert.Contains(t, gotOutput, gameConfig["guess"])
				assert.Contains(t, gotOutput, fmt.Sprintf(gameConfig["equal"], "0s", 1))
			})
		}
	})

	t.Run("invalid play again inputs", func(t *testing.T) {
		testCases := []struct {
			description      string
			invalidInput     string
			wantErrorMessage string
		}{
			{
				description:      "non-numeric input",
				invalidInput:     "invalid",
				wantErrorMessage: parser.ParseNumberMessage,
			},
			{
				description:      "less than valid range",
				invalidInput:     "0",
				wantErrorMessage: fmt.Sprintf(parser.NumberRangeMessage, 1, 2),
			},
			{
				description:      "greater than valid range",
				invalidInput:     "3",
				wantErrorMessage: fmt.Sprintf(parser.NumberRangeMessage, 1, 2),
			},
		}

		for _, tc := range testCases {
			t.Run(tc.description, func(t *testing.T) {
				mockInputSource := &MockInputSource{
					PlayerInput:       []string{"test"},
					DifficultyInput:   []string{"3"},
					GuessNumberInputs: []string{"50"},
					PlayAgainInput:    []string{tc.invalidInput, "2"},
				}
				gotWriter, game := initGame(mockInputSource)
				game.Play(fakeRandomNumber, stubScoreStore)

				gotOutput := gotWriter.String()

				assert.Contains(t, gotOutput, fmt.Sprintf(gameConfig["level"], "Hard"))
				assert.Contains(t, gotOutput, fmt.Sprintf(gameConfig["equal"], "0s", 1))
				assert.Contains(t, gotOutput, tc.wantErrorMessage)
				assert.Contains(t, gotOutput, gameConfig["again"])
				assert.Contains(t, gotOutput, gameConfig["bye"])
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

type StubScoreStore struct {
	isEmpty bool
}

func (s *StubScoreStore) Load() store.Scores {
	if s.isEmpty {
		return store.Scores{}
	}

	return store.Scores{
		{
			Player:   "Test",
			Level:    "Hard",
			Attempts: 3,
			Time:     30 * time.Second,
		},
	}
}

func (s *StubScoreStore) Add(store.Score) (store.Scores, error) {
	return store.Scores{
		{
			Player:   "Test",
			Level:    "Hard",
			Attempts: 3,
			Time:     30 * time.Second,
		},
	}, nil
}

type MockInputSource struct {
	PlayerInput []string
	playerIndex int

	DifficultyInput []string
	difficultyIndex int

	GuessNumberInputs []string
	guessNumberIndex  int

	PlayAgainInput []string
	playAgainIndex int
}

func (m *MockInputSource) NextPlayer() (string, error) {
	if m.playerIndex < len(m.PlayerInput) {
		value := m.PlayerInput[m.playerIndex]
		m.playerIndex++
		return value, nil
	}

	return "", fmt.Errorf("no more player inputs")
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

func (m *MockInputSource) NextPlayAgainInput() (string, error) {
	if m.playAgainIndex < len(m.PlayAgainInput) {
		value := m.PlayAgainInput[m.playAgainIndex]
		m.playAgainIndex++
		return value, nil
	}

	return "", fmt.Errorf("no more play again inputs")
}
