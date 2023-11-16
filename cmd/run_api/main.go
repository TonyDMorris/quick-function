package main

import (
	"net/http"
	"os"

	"github.com/TonyDMorris/quick-function/pkg/logging"
	"github.com/TonyDMorris/quick-function/service/app"
	"github.com/bradleyfalzon/ghinstallation/v2"
	"github.com/caarlos0/env/v10"
	"github.com/golang-jwt/jwt"
	"github.com/google/go-github/v56/github"
	"github.com/joho/godotenv"
)

type Config struct {
	RSA   string `env:"RSA"`
	AppID int64  `env:"GITHUB_APP_ID"`
}

func main() {

	var config Config

	err := godotenv.Load("../../.env")
	if err != nil {
		logging.Logger.Info(err.Error())
	}
	err = env.Parse(&config)
	if err != nil {
		panic(err)
	}

	key, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(config.RSA))
	if err != nil {
		panic(err)
	}

	itr := ghinstallation.NewAppsTransportFromPrivateKey(http.DefaultTransport, config.AppID, key)

	client := github.NewClient(&http.Client{
		Transport: itr,
	})

	app := app.NewApi(app.Config{
		Port: 8080,
	},
		client)

	if err := app.Run(); err != nil {
		logging.Logger.Error(err.Error())
		os.Exit(1)
	}

}
