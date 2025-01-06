package main

import (
	_ "embed"
	"os"

	"github.com/go-number-guessing-game/internal/cli"
	"github.com/go-number-guessing-game/internal/config"
	"github.com/go-number-guessing-game/internal/number"
	"github.com/go-number-guessing-game/internal/service"
)

func main() {
	gameConfig := config.LoadConfig("yaml", "configs/data.yaml")
	cliInputSource := &cli.CliInput{Source: os.Stdin}
	game := service.Game{
		Writer:      os.Stdout,
		InputSource: cliInputSource,
		GameConfig:  gameConfig,
	}
	randomNumber := number.NewRandomNumber()
	game.Play(randomNumber)
}
