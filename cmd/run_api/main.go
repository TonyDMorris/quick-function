package main

import (
	"os"

	tokenservice "github.com/TonyDMorris/quick-function/pkg/github_token_service/client"
	"github.com/TonyDMorris/quick-function/service/app"
	"github.com/TonyDMorris/quick-function/service/github/client/service"
	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"

	"go.uber.org/zap"
)

type Config struct {
	RSA   string `env:"RSA"`
	AppID int    `env:"GITHUB_APP_ID"`
}

func main() {

	var config Config
	log, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	err = godotenv.Load("../../.env")
	if err != nil {
		log.Info(err.Error())
	}
	err = env.Parse(&config)
	if err != nil {
		panic(err)
	}

	tokenService := tokenservice.NewTokenSerivce([]byte(config.RSA), config.AppID)

	gitService := service.NewService(tokenService)

	app := app.NewApi(app.Config{
		Port: 8080,
	}, gitService)

	if err := app.Run(); err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}

}
