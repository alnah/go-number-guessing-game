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

func main() {
	gameConfig := config.LoadConfig("yaml", "configs/app.yaml")
	randomNumber := game.NewRandomNumber()

	gameStore := &store.ScoresStore{FilePath: "internal/data/scores.json"}
	cliInputSource := &cli.CliInput{Source: os.Stdin}
	game := service.Game{
		Writer:      os.Stdout,
		InputSource: cliInputSource,
		GameConfig:  gameConfig,
	}

	game.PlayGame(randomNumber, gameStore)
}
