package game_test

import (
	"testing"

	"github.com/go-number-guessing-game/internal/game"
	"github.com/stretchr/testify/assert"
)

func TestUnitNewRandomNumber(t *testing.T) {
	t.Run("return random number between 1 and 100", func(t *testing.T) {
		got := game.NewRandomNumber()

		assert.Greater(t, got, 0, "random number greater than 0")
		assert.LessOrEqual(t, got, 100, "random number less or equal than 100")
		assert.IsType(t, int(got), got)
	})
}

func TestUnitPlayTurn(t *testing.T) {
	t.Run("return complete game state", func(t *testing.T) {
		testCases := []struct {
			description string
			level       string
			maxAttempts int
			beforeTurns game.Turns
			afterTurns  game.Turns
		}{
			{
				description: "for hard level",
				level:       "Hard",
				maxAttempts: 3,
				beforeTurns: game.Turns{
					{GuessNumber: 25},
					{GuessNumber: 75},
					{GuessNumber: 50},
				},
				afterTurns: game.Turns{
					{
						GuessNumber: 25,
						Outcome:     toPointer(1),
						Difference:  toPointer(25),
					},
					{
						GuessNumber: 75,
						Outcome:     toPointer(-1),
						Difference:  toPointer(25),
					},
					{
						GuessNumber: 50,
						Outcome:     toPointer(0),
						Difference:  toPointer(0),
					},
				},
			},
			{
				description: "for medium level",
				level:       "Medium",
				maxAttempts: 5,
				beforeTurns: game.Turns{
					{GuessNumber: 30},
					{GuessNumber: 70},
					{GuessNumber: 40},
					{GuessNumber: 60},
					{GuessNumber: 50},
				},
				afterTurns: game.Turns{
					{
						GuessNumber: 30,
						Outcome:     toPointer(1),
						Difference:  toPointer(20),
					},
					{
						GuessNumber: 70,
						Outcome:     toPointer(-1),
						Difference:  toPointer(20),
					},
					{
						GuessNumber: 40,
						Outcome:     toPointer(1),
						Difference:  toPointer(10),
					},
					{
						GuessNumber: 60,
						Outcome:     toPointer(-1),
						Difference:  toPointer(10),
					},
					{
						GuessNumber: 50,
						Outcome:     toPointer(0),
						Difference:  toPointer(0),
					},
				},
			},
			{
				description: "for easy level",
				level:       "Easy",
				maxAttempts: 10,
				beforeTurns: game.Turns{
					{GuessNumber: 25},
					{GuessNumber: 70},
					{GuessNumber: 30},
					{GuessNumber: 65},
					{GuessNumber: 35},
					{GuessNumber: 60},
					{GuessNumber: 40},
					{GuessNumber: 55},
					{GuessNumber: 45},
					{GuessNumber: 50},
				},
				afterTurns: game.Turns{
					{
						GuessNumber: 25,
						Outcome:     toPointer(1),
						Difference:  toPointer(25),
					},
					{
						GuessNumber: 70,
						Outcome:     toPointer(-1),
						Difference:  toPointer(20),
					},
					{
						GuessNumber: 30,
						Outcome:     toPointer(1),
						Difference:  toPointer(20),
					},
					{
						GuessNumber: 65,
						Outcome:     toPointer(-1),
						Difference:  toPointer(15),
					},
					{
						GuessNumber: 35,
						Outcome:     toPointer(1),
						Difference:  toPointer(15),
					},
					{
						GuessNumber: 60,
						Outcome:     toPointer(-1),
						Difference:  toPointer(10),
					},
					{
						GuessNumber: 40,
						Outcome:     toPointer(1),
						Difference:  toPointer(10),
					},
					{
						GuessNumber: 55,
						Outcome:     toPointer(-1),
						Difference:  toPointer(5),
					},
					{
						GuessNumber: 45,
						Outcome:     toPointer(1),
						Difference:  toPointer(5),
					},
					{
						GuessNumber: 50,
						Outcome:     toPointer(0),
						Difference:  toPointer(0),
					},
				},
			},
		}

		for _, tc := range testCases {
			t.Run(tc.description, func(t *testing.T) {
				gameState := game.GameState{
					Level:        tc.level,
					MaxAttempts:  tc.maxAttempts,
					RandomNumber: 50,
					Turns:        game.Turns{},
				}

				for _, turn := range tc.beforeTurns {
					err := gameState.PlayTurn(turn)
					assert.NoError(t, err)
				}

				got := gameState
				want := game.GameState{
					Level:        tc.level,
					MaxAttempts:  tc.maxAttempts,
					RandomNumber: 50,
					Turns:        tc.afterTurns,
				}

				assert.Equal(t, want, got)
			})
		}
	})
	t.Run("return error when turns length greater than max attempts", func(t *testing.T) {
		testCases := []struct {
			description string
			level       string
			maxAttempts int
			turns       game.Turns
		}{
			{
				description: "for hard level",
				level:       "Hard",
				maxAttempts: 3,
				turns: game.Turns{
					{GuessNumber: 25},
					{GuessNumber: 75},
					{GuessNumber: 51},
					{GuessNumber: 50}, // too many turns
				},
			},
			{
				description: "for medium level",
				level:       "Medium",
				maxAttempts: 5,
				turns: game.Turns{
					{GuessNumber: 30},
					{GuessNumber: 70},
					{GuessNumber: 40},
					{GuessNumber: 60},
					{GuessNumber: 51},
					{GuessNumber: 50}, // too many turns
				},
			},
			{
				description: "for easy level",
				level:       "Easy",
				maxAttempts: 10,
				turns: game.Turns{
					{GuessNumber: 25},
					{GuessNumber: 70},
					{GuessNumber: 30},
					{GuessNumber: 65},
					{GuessNumber: 35},
					{GuessNumber: 60},
					{GuessNumber: 40},
					{GuessNumber: 55},
					{GuessNumber: 45},
					{GuessNumber: 51},
					{GuessNumber: 50}, // too many turns
				},
			},
		}

		for _, tc := range testCases {
			t.Run(tc.description, func(t *testing.T) {
				gameState := game.GameState{
					Level:        tc.level,
					MaxAttempts:  tc.maxAttempts,
					RandomNumber: 50,
					Turns:        game.Turns{},
				}

				want := game.NewTurnsLengthError(
					gameState.Turns,
					gameState.MaxAttempts,
				)

				var got error
				for _, turn := range tc.turns {
					got = gameState.PlayTurn(turn)
					if len(gameState.Turns)-1 < gameState.MaxAttempts {
						assert.NoError(t, got)
					}
				}

				assert.NotNil(t, got)
				assert.ErrorAs(t, got, &want)
			})
		}
	})

	t.Run("return error when incorrect level", func(t *testing.T) {
		gameState := game.GameState{
			Level:        "incorrect",
			MaxAttempts:  3,
			RandomNumber: 50,
			Turns:        game.Turns{},
		}

		turn := game.Turn{GuessNumber: 1}

		want := game.NewLevelError()
		got := gameState.PlayTurn(turn)

		assert.NotNil(t, want)
		assert.ErrorAs(t, got, &want)
	})

	t.Run("return error when incorrect max attempts", func(t *testing.T) {
		testCases := []struct {
			description string
			level       string
			maxAttempts int
		}{
			{
				description: "not matching between easy and 10",
				level:       "Easy",
				maxAttempts: 3,
			},
			{
				description: "not matching between medium and 5",
				level:       "Medium",
				maxAttempts: 3,
			},
			{
				description: "not matching between hard and 3",
				level:       "Hard",
				maxAttempts: 10,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.description, func(t *testing.T) {
				turn := game.Turn{GuessNumber: 1}
				gameState := game.GameState{
					Level:        tc.level,
					MaxAttempts:  tc.maxAttempts,
					RandomNumber: 50,
					Turns:        game.Turns{},
				}

				want := game.NewMaxAttemptsError()
				got := gameState.PlayTurn(turn)

				assert.NotNil(t, want)
				assert.ErrorAs(t, got, &want)
			})
		}
	})

	t.Run("return error when random number already found", func(t *testing.T) {
		gameState := game.GameState{
			Level:        "Hard",
			MaxAttempts:  3,
			RandomNumber: 50,
			Turns:        game.Turns{},
		}

		err := gameState.PlayTurn(game.Turn{GuessNumber: 50})
		assert.NoError(t, err)

		want := game.NewRandomNumberFoundError()
		got := gameState.PlayTurn(game.Turn{GuessNumber: 51})
		assert.ErrorAs(t, got, &want)
	})
}

func TestUnitGetAttempts(t *testing.T) {
	gameState := game.GameState{
		Level:        "Easy",
		MaxAttempts:  10,
		RandomNumber: 50,
		Turns: game.Turns{{
			GuessNumber: 50,
			Outcome:     toPointer(0),
			Difference:  toPointer(0),
		}},
	}

	got := gameState.GetAttempts()
	assert.Equal(t, 1, got, "should be 1")
}

func TestUnitGetLastTurn(t *testing.T) {
	t.Run("return last turn", func(t *testing.T) {
		gameState := game.GameState{
			Level:        "Easy",
			MaxAttempts:  10,
			RandomNumber: 50,
			Turns: game.Turns{{
				GuessNumber: 50,
				Outcome:     toPointer(0),
				Difference:  toPointer(0),
			}},
		}

		got, err := gameState.GetLastTurn()

		assert.NoError(t, err)
		assert.Equal(t, game.Turn{
			GuessNumber: 50,
			Outcome:     toPointer(0),
			Difference:  toPointer(0),
		}, got)
	})

	t.Run("error when no turns", func(t *testing.T) {
		gameState := game.GameState{
			Level:        "Easy",
			MaxAttempts:  10,
			RandomNumber: 50,
			Turns:        game.Turns{},
		}

		want := game.NewEmptyTurnsError()
		_, got := gameState.GetLastTurn()

		assert.NotNil(t, got)
		assert.ErrorAs(t, got, &want)
	})
}

func TestUnitNoMoreAttempts(t *testing.T) {
	t.Run("return", func(t *testing.T) {
		testCases := []struct {
			description string
			turns       game.Turns
			want        bool
		}{
			{
				description: "true when no more attempts",
				turns: game.Turns{
					{
						GuessNumber: 53,
						Outcome:     toPointer(-1),
						Difference:  toPointer(3),
					},
					{
						GuessNumber: 52,
						Outcome:     toPointer(-1),
						Difference:  toPointer(2),
					},
					{
						GuessNumber: 51,
						Outcome:     toPointer(-1),
						Difference:  toPointer(1),
					},
				},
				want: true,
			},
			{
				description: "false when more attempts",
				turns: game.Turns{
					{
						GuessNumber: 53,
						Outcome:     toPointer(-1),
						Difference:  toPointer(3),
					},
					{
						GuessNumber: 52,
						Outcome:     toPointer(-1),
						Difference:  toPointer(2),
					},
				},
				want: false,
			},
		}
		for _, tc := range testCases {
			gameState := game.GameState{
				Level:        "Hard",
				MaxAttempts:  3,
				RandomNumber: 50,
				Turns:        tc.turns,
			}

			got := gameState.NoMoreAttempts()
			assert.Equal(t, tc.want, got)
		}
	})
}

func toPointer(value int) *int {
	return &value
}
