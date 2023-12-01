package main

import (
	"net/http"
	"os"

	gpt "github.com/TonyDMorris/quick-function/pkg/gpt/client"
	"github.com/TonyDMorris/quick-function/pkg/logging"
	strapi "github.com/TonyDMorris/quick-function/pkg/strapi/client"
	"github.com/TonyDMorris/quick-function/service/app"
	"github.com/bradleyfalzon/ghinstallation/v2"
	"github.com/caarlos0/env/v10"
	"github.com/golang-jwt/jwt"
	"github.com/google/go-github/v56/github"
	"github.com/joho/godotenv"
)

type Config struct {
	RSA           string `env:"RSA,required"`
	AppID         int64  `env:"GITHUB_APP_ID,required"`
	ChatGPTAPIKey string `env:"CHAT_GPT_API_KEY,required"`
	StrapiAPIKey  string `env:"STRAPI_API_KEY,required"`
	StrapiBaseURL string `env:"STRAPI_BASE_URL,required"`
}

func main() {

	var config Config

	err := godotenv.Load("../../.env")
	if err != nil {
		logging.Logger.Info(err.Error())
	}
	err = env.Parse(&config)
	if err != nil {
		logging.Logger.Error(err.Error())
		os.Exit(1)
	}

	key, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(config.RSA))
	if err != nil {
		logging.Logger.Error(err.Error())
		os.Exit(1)
	}

	itr := ghinstallation.NewAppsTransportFromPrivateKey(http.DefaultTransport, config.AppID, key)

	client := github.NewClient(&http.Client{
		Transport: itr,
	})

	gptClient := gpt.NewChatClient(config.ChatGPTAPIKey)

	strapiClient := strapi.NewClient(config.StrapiAPIKey, config.StrapiBaseURL)

	app := app.NewApi(
		app.Config{
			Port: 8080,
		},
		client, gptClient,
		strapiClient,
	)

	if err := app.Run(); err != nil {
		logging.Logger.Error(err.Error())
		os.Exit(1)
	}

}
