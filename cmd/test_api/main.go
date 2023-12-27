package main

import (
	"net/http"
	"os"
	"time"

	gpt "github.com/TonyDMorris/quick-function/pkg/gpt/client"
	"github.com/TonyDMorris/quick-function/pkg/logging"
	strapi "github.com/TonyDMorris/quick-function/pkg/strapi/client"
	"github.com/TonyDMorris/quick-function/pkg/strapi/models"
	"github.com/TonyDMorris/quick-function/service/app"
	"github.com/bradleyfalzon/ghinstallation/v2"
	"github.com/caarlos0/env/v10"
	"github.com/golang-jwt/jwt"
	"github.com/google/go-github/v56/github"
	"github.com/hashicorp/go-retryablehttp"
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
	trns := retryablehttp.NewClient().HTTPClient.Transport
	itr := ghinstallation.NewAppsTransportFromPrivateKey(trns, config.AppID, key)

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

	lastGen := time.Now().Add(-time.Hour * 24 * 7 * 4)
	job := models.RepositoryConfiguration{
		ID: 55,
		Repository: &models.Repository{
			ID:           150,
			Name:         "quick-function",
			FullName:     "TonyDMorris/quick-function",
			Private:      false,
			RepositoryID: "718219828",
		},
		Installation: &models.Installation{
			ID:             18,
			InstallationID: "44656141",
			Username:       "TonyDMorris",
		},
		Cron: "4 weeks",

		LastGeneration: &lastGen,
	}

	err = app.HandleRepositoryConfigurationScheduledJob(job)
	if err != nil {
		logging.Logger.Error(err.Error())
		os.Exit(1)
	}

}
