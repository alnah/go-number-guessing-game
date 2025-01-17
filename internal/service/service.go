package service

import (
	"fmt"
	"io"
	"time"

	"github.com/go-number-guessing-game/internal/cli"
	"github.com/go-number-guessing-game/internal/game"
	"github.com/go-number-guessing-game/internal/number"
	"github.com/go-number-guessing-game/internal/timer"
)

// Newline is a constant string representing a newline character.
const Newline string = "\n"

// Game represents the state and configuration of the number guessing game.
type Game struct {
	Writer      io.Writer
	InputSource cli.InputSource
	GameConfig  map[string]string
}

// Play starts the game by displaying the greeting and difficulty options,
// then initializing the game state and entering the game loop.
func (g *Game) Play(randomNumber int) {
	cli.Display(g.Writer, []string{
		g.GameConfig["greeting"],
		g.GameConfig["spacer"],
	})

gameLoop:
	for {
		cli.Display(g.Writer, g.GameConfig["difficulty"])

		level, maxAttempts := g.getUserDifficultyInput()
		gameState := g.initGameState(level, maxAttempts, randomNumber)
		g.playTurns(gameState)

		playGain := g.getPlayAgainInput()
		if playGain {
			continue gameLoop
		} else {
			cli.Display(g.Writer, []string{
				g.GameConfig["bye"],
				g.GameConfig["newline"],
			})
			break gameLoop
		}
	}
}

// initGameState initializes the game state with the selected difficulty level,
// maximum attempts, and the random number to be guessed.
func (g *Game) initGameState(level string, maxAttempts, randomNumber int) game.GameState {
	return game.GameState{
		Level:        level,
		MaxAttempts:  maxAttempts,
		RandomNumber: randomNumber,
		Turns:        game.Turns{},
	}
}

// playTurns manages the game loop, processing user guesses and displaying
// feedback until the game ends.
func (g *Game) playTurns(gameState game.GameState) {
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
			gameTime = gameTimer.End()
			time := gameTime.String()
			attempts := gameState.GetAttempts()

			cli.Display(g.Writer, []string{
				fmt.Sprintf(g.GameConfig["equal"], time, attempts),
				g.GameConfig["newline"],
			})
			break turnLoop
		}
	}
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

// getUserDifficultyInput prompts the user for difficulty input and returns
// the selected level and maximum attempts.
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

		level, maxAttempts, err = number.ParseDifficultyInput(input)
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

// getUserGuessNumberInput prompts the user for a guess number and returns
// the parsed integer value.
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

		guessNumber, err = number.ParseGuessNumberInput(input)
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

// getPlayAgainInput prompt the user for a play gain number and returns
// the parsed boolean.
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

		playAgain, err = number.ParsePlayAgainInput(input)
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
