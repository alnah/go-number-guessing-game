// Package main initializes the number guessing game, loading the necessary
// configurations, setting up the game state, and handling user input.
package main

import (
	_ "embed"
	"os"

	"github.com/go-number-guessing-game/internal/cli"
	"github.com/go-number-guessing-game/internal/config"
	"github.com/go-number-guessing-game/internal/game"
	"github.com/go-number-guessing-game/internal/service"
	"github.com/go-number-guessing-game/internal/store"
)

// The main function serves as the entry point for the app.
func main() {
	// Load game configuration from a YAML file.
	gameConfig := config.LoadConfig("yaml", "configs/app.yaml")

	// Generate a new random number for the game.
	randomNumber := game.NewRandomNumber()

	// Initialize the scores store with the specified file path.
	gameStore := &store.ScoresStore{FilePath: "internal/data/scores.json"}

	// Create a CLI input source to read user input from standard input.
	cliInputSource := &cli.CliInput{Source: os.Stdin}

	// Set up the game with the writer, input source, and configuration.
	game := service.Game{
		Writer:      os.Stdout,
		InputSource: cliInputSource,
		GameConfig:  gameConfig,
	}

	// Start the game with the generated random number and the scores store.
	game.PlayGame(randomNumber, gameStore)
}
