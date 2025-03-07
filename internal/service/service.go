// Package service orchestrates the game logic and integrates other packages.
// It provides the PlayGame method to manage gameplay and user interactions.
package service

import (
	"fmt"
	"io"
	"time"

	"github.com/go-number-guessing-game/internal/cli"
	"github.com/go-number-guessing-game/internal/game"
	"github.com/go-number-guessing-game/internal/parser"
	"github.com/go-number-guessing-game/internal/store"
	"github.com/go-number-guessing-game/internal/timer"
)

// Game encapsulates the writer and input source interfaces for testing.
// It also holds a configuration map for displaying messages in the CLI.
type Game struct {
	Writer      io.Writer
	InputSource cli.InputSource
	GameConfig  map[string]string
}

// PlayGame initiates the game with a random number and a store interface.
// It manages user inputs, orchestrates game logic, and persists scores.
// If the user guesses correctly, a new random number is generated.
func (g *Game) PlayGame(randomNumber int, store store.Store) {
	cli.Display(g.Writer, []string{
		g.GameConfig["greeting"],
		g.GameConfig["spacer"],
	})

gameLoop:
	for {
		cli.Display(g.Writer, g.GameConfig["player"])
		player := g.getPlayerInput()

		cli.Display(g.Writer, g.GameConfig["difficulty"])
		level, maxAttempts := g.getUserDifficultyInput()

		gameState := g.initGameState(level, maxAttempts, randomNumber)
		found, attempts, time := g.playTurns(gameState)

		if found {
			g.displayScores(player, level, attempts, time, store)
		}

		playAgain := g.getPlayAgainInput()
		switch {
		case playAgain && found:
			randomNumber = game.NewRandomNumber()
			continue gameLoop

		case playAgain && !found:
			continue gameLoop

		default:
			cli.Display(g.Writer, []string{
				g.GameConfig["bye"],
				g.GameConfig["newline"],
			})
			break gameLoop
		}
	}
}

func (g *Game) initGameState(
	level string,
	maxAttempts,
	randomNumber int,
) game.GameState {
	return game.GameState{
		Level:        level,
		MaxAttempts:  maxAttempts,
		RandomNumber: randomNumber,
		Turns:        game.Turns{},
	}
}

func (g *Game) playTurns(gameState game.GameState) (bool, int, time.Duration) {
	var found bool
	var attempts int
	var gameTime time.Duration

	gameTimer := timer.NewGameTimer()
	gameTimer.Start()

turnLoop:
	for {
		if gameState.NoMoreAttempts() {
			cli.Display(g.Writer, []string{
				g.GameConfig["max_attempts"],
				g.GameConfig["newline"],
			})
			gameTimer.End()
			break turnLoop
		}

		guessNumber := g.getUserGuessNumberInput()
		err := gameState.PlayTurn(game.Turn{GuessNumber: guessNumber})
		if err != nil {
			cli.Display(g.Writer, []string{
				err.Error(),
				g.GameConfig["newline"],
			})
		}

		lastTurn, _ := gameState.GetLastTurn()

		switch *lastTurn.Outcome {
		case 1:
			cli.Display(g.Writer, []string{
				fmt.Sprintf(g.GameConfig["greater"], guessNumber),
				g.GameConfig["newline"],
				g.giveHint(lastTurn),
				g.GameConfig["spacer"],
			})
			continue turnLoop

		case -1:
			cli.Display(g.Writer, []string{
				fmt.Sprintf(g.GameConfig["less"], guessNumber),
				g.GameConfig["newline"],
				g.giveHint(lastTurn),
				g.GameConfig["spacer"],
			})
			continue turnLoop

		case 0:
			found = true
			gameTime = gameTimer.End()
			stringTime := gameTime.String()
			attempts = gameState.GetAttempts()

			cli.Display(g.Writer, []string{
				fmt.Sprintf(g.GameConfig["equal"], stringTime, attempts),
				g.GameConfig["newline"],
			})
			break turnLoop
		}
	}

	return found, attempts, gameTime
}

func (g *Game) getPlayerInput() string {
	var player string

playerLoop:
	for {
		input, err := g.InputSource.NextPlayerInput()
		if err != nil {
			cli.Display(g.Writer, []string{
				err.Error(),
				g.GameConfig["spacer"],
				g.GameConfig["player"],
			})
			continue playerLoop
		}

		player, err = parser.ParsePlayerInput(input)
		if err != nil {
			cli.Display(g.Writer, []string{
				err.Error(),
				g.GameConfig["spacer"],
				g.GameConfig["player"],
			})
			continue playerLoop
		}
		break playerLoop
	}
	return player
}

func (g *Game) getUserDifficultyInput() (string, int) {
	var level string
	var maxAttempts int

difficultyLoop:
	for {
		input, err := g.InputSource.NextDifficultyInput()
		if err != nil {
			cli.Display(g.Writer, []string{
				err.Error(),
				g.GameConfig["spacer"],
				g.GameConfig["difficulty"],
			})
			continue difficultyLoop
		}

		level, maxAttempts, err = parser.ParseDifficultyInput(input)
		if err != nil {
			cli.Display(g.Writer, []string{
				err.Error(),
				g.GameConfig["spacer"],
				g.GameConfig["difficulty"],
			})
			continue difficultyLoop
		}
		break difficultyLoop
	}
	cli.Display(g.Writer, []string{
		fmt.Sprintf(g.GameConfig["level"], level),
		g.GameConfig["spacer"],
	})

	return level, maxAttempts
}

func (g *Game) getUserGuessNumberInput() int {
	var guessNumber int

guessNumberLoop:
	for {
		cli.Display(g.Writer, g.GameConfig["guess"])
		input, err := g.InputSource.NextGuessNumberInput()
		if err != nil {
			cli.Display(g.Writer, []string{
				err.Error(),
				g.GameConfig["spacer"],
			})
			continue guessNumberLoop
		}

		guessNumber, err = parser.ParseGuessNumberInput(input)
		if err != nil {
			cli.Display(g.Writer, []string{
				err.Error(),
				g.GameConfig["spacer"],
			})
			continue guessNumberLoop
		}
		break guessNumberLoop
	}

	return guessNumber
}

func (g *Game) getPlayAgainInput() bool {
	var playAgain bool

playAgainLoop:
	for {
		cli.Display(g.Writer, g.GameConfig["again"])

		input, err := g.InputSource.NextPlayAgainInput()
		if err != nil {
			cli.Display(g.Writer, []string{
				err.Error(),
				g.GameConfig["spacer"],
			})
			continue playAgainLoop
		}

		playAgain, err = parser.ParsePlayAgainInput(input)
		if err != nil {
			cli.Display(g.Writer, []string{
				err.Error(),
				g.GameConfig["spacer"],
			})
			continue playAgainLoop
		}
		break playAgainLoop
	}

	return playAgain
}

func (g *Game) giveHint(lastTurn game.Turn) string {
	var hint string

	switch {
	case *lastTurn.Difference == 1:
		hint = g.GameConfig["very_close_1"]
	case *lastTurn.Difference == 2:
		hint = g.GameConfig["very_close_2"]
	case *lastTurn.Difference == 3:
		hint = g.GameConfig["very_close_3"]
	case *lastTurn.Difference == 4:
		hint = g.GameConfig["close_1"]
	case *lastTurn.Difference == 5:
		hint = g.GameConfig["close_2"]
	case *lastTurn.Difference > 5 && *lastTurn.Difference < 10:
		hint = g.GameConfig["far"]
	default:
		hint = g.GameConfig["very_far"]
	}

	return hint
}

func (g *Game) displayScores(player string,
	level string,
	attempts int,
	time time.Duration,
	gameStore store.Store,
) {
	score := store.Score{
		Player:   player,
		Level:    level,
		Attempts: attempts,
		Time:     time,
	}
	scores, _ := gameStore.Add(score)
	cli.Display(g.Writer, []string{
		g.GameConfig["spacer"],
		scores.String(),
		g.GameConfig["spacer"],
	})
}
