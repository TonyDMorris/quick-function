package app

import (
	"database/sql"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/go-github/v56/github"
	"github.com/hashicorp/go-retryablehttp"
)

type Config struct {
	Port int
}

type App struct {
	server       *gin.Engine
	client       *http.Client
	githubClient *github.Client
	DB           *sql.DB
	port         int
}

func NewApi(c Config, githubClient *github.Client) *App {
	return &App{
		server:       gin.Default(),
		client:       retryablehttp.NewClient().StandardClient(),
		githubClient: githubClient,
		port:         c.Port,
	}
}

func (a *App) Run() error {
	a.SetupRoutes()
	return a.server.Run()
}

func (a *App) SetupRoutes() {

	a.server.POST("/test", a.Test)

}

func (a *App) Test(c *gin.Context) {

	installations, resp, err := a.githubClient.Apps.ListInstallations(c.Request.Context(), nil)

	if err != nil {
		body, _ := io.ReadAll(resp.Body)
		c.JSON(resp.StatusCode, gin.H{
			"error": err.Error(),
			"body":  string(body),
		})
		return
	}

	for _, installation := range installations {

		installationToken, tokenResponse, err := a.githubClient.Apps.CreateInstallationToken(c.Request.Context(), installation.GetID(), nil)

		if err != nil {
			body, _ := io.ReadAll(tokenResponse.Body)
			c.JSON(tokenResponse.StatusCode, gin.H{
				"error": err.Error(),
				"body":  string(body),
			})
			return
		}

		userClient := github.NewClient(http.DefaultClient).WithAuthToken(installationToken.GetToken())

		repos, repoResponse, err := userClient.Apps.ListRepos(c.Request.Context(), &github.ListOptions{
			Page:    1,
			PerPage: 100,
		})

		if err != nil {
			body, _ := io.ReadAll(repoResponse.Body)
			c.JSON(repoResponse.StatusCode, gin.H{
				"error": err.Error(),
				"body":  string(body),
			})
			return
		}
		var returnedEvents []*github.Event
		for _, repo := range repos.Repositories {
			events, eventResponse, err := userClient.Activity.ListRepositoryEvents(c.Request.Context(), repo.GetOwner().GetLogin(), repo.GetName(), nil)
			if err != nil {
				body, _ := io.ReadAll(eventResponse.Body)
				c.JSON(eventResponse.StatusCode, gin.H{
					"error": err.Error(),
					"body":  string(body),
				})
				return
			}
			returnedEvents = append(returnedEvents, events...)
		}

		c.JSON(http.StatusOK, gin.H{
			"events": returnedEvents,
		})
	}

}
