package service

import (
	"fmt"
	"io"

	"github.com/go-number-guessing-game/internal/cli"
	"github.com/go-number-guessing-game/internal/game"
	"github.com/go-number-guessing-game/internal/number"
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
		g.GameConfig["newline"],
		g.GameConfig["difficulty"],
	})

	level, maxAttempts := g.getUserDifficultyInput()
	gameState := g.initGameState(level, maxAttempts, randomNumber)

	g.playTurns(gameState)
}

// playTurns manages the game loop, processing user guesses and displaying
// feedback until the game ends.
func (g *Game) playTurns(gameState game.GameState) {
turnLoop:
	for {
		if gameState.NoMoreAttempts() {
			cli.Display(g.Writer, []string{
				g.GameConfig["max_attempts"],
				g.GameConfig["newline"],
			})
			break turnLoop
		}

		guessNumber := g.getUserGuessNumberInput()
		gameState.PlayTurn(game.Turn{GuessNumber: guessNumber})
		lastTurn, _ := gameState.GetLastTurn()

		switch *lastTurn.Outcome {
		case 1:
			cli.Display(g.Writer, []string{
				fmt.Sprintf(g.GameConfig["greater"], guessNumber),
				g.GameConfig["spacer"],
			})
			continue

		case -1:
			cli.Display(g.Writer, []string{
				fmt.Sprintf(g.GameConfig["less"], guessNumber),
				g.GameConfig["spacer"],
			})
			continue

		case 0:
			cli.Display(g.Writer, []string{
				fmt.Sprintf(g.GameConfig["equal"], gameState.GetAttempts()),
				g.GameConfig["newline"],
			})
			break turnLoop
		}
	}
}

// getUserDifficultyInput prompts the user for difficulty input and returns
// the selected level and maximum attempts.
func (g *Game) getUserDifficultyInput() (string, int) {
	var level string
	var maxAttempts int

	for {
		input, err := g.InputSource.NextDifficultyInput() // cli.GetUserInput(g.Reader)
		if err != nil {
			cli.Display(g.Writer, []string{
				err.Error(),
				g.GameConfig["spacer"],
				g.GameConfig["difficulty"],
			})
			continue
		}

		level, maxAttempts, err = number.ParseDifficultyInput(input)
		if err != nil {
			cli.Display(g.Writer, []string{
				err.Error(),
				g.GameConfig["spacer"],
				g.GameConfig["difficulty"],
			})
			continue
		}
		break
	}
	cli.Display(g.Writer, []string{
		fmt.Sprintf(g.GameConfig["level"], level),
		g.GameConfig["newline"],
	})

	return level, maxAttempts
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

// getUserGuessNumberInput prompts the user for a guess number and returns
// the parsed integer value.
func (g *Game) getUserGuessNumberInput() int {
	var guessNumber int

	for {
		cli.Display(g.Writer, g.GameConfig["guess"])
		input, err := g.InputSource.NextGuessNumberInput() // cli.GetUserInput(g.Reader)
		if err != nil {
			cli.Display(g.Writer, []string{
				err.Error(),
				g.GameConfig["spacer"],
			})
			continue
		}

		guessNumber, err = number.ParseGuessNumberInput(input)
		if err != nil {
			cli.Display(g.Writer, []string{
				err.Error(),
				g.GameConfig["spacer"],
			})
			continue
		}
		break
	}

	return guessNumber
}
