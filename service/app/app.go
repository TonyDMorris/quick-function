package app

import (
	"io"
	"net/http"
	"time"

	"github.com/TonyDMorris/quick-function/pkg/logging"
	"github.com/TonyDMorris/quick-function/service/github/client/models"
	"github.com/gin-gonic/gin"
	"github.com/hashicorp/go-retryablehttp"
	"go.uber.org/zap"
)

type Config struct {
	Port int
}

type githubService interface {
	GetInstallationsForApp() ([]models.Installation, error)
}

type App struct {
	server        *gin.Engine
	client        *http.Client
	githubService githubService
	port          int
}

func NewApi(c Config, githubService githubService) *App {
	return &App{
		server:        gin.Default(),
		client:        retryablehttp.NewClient().StandardClient(),
		githubService: githubService,
		port:          c.Port,
	}
}

func (a *App) Run() error {
	a.SetupRoutes()
	errChan := make(chan error)
	go func() {
		if err := a.job(); err != nil {
			logging.Logger.Error("error running job", zap.Error(err))
			errChan <- err
		}
	}()

	go func() {
		if err := a.server.Run(); err != nil {
			logging.Logger.Error("error running server", zap.Error(err))
			errChan <- err
		}
	}()

	return <-errChan
}

func (a *App) SetupRoutes() {
	a.server.GET("/installations", a.getInstallations)
	a.server.POST("/installations", a.PostInstallation)

}

func (a *App) PostInstallation(c *gin.Context) {

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {

		logging.Logger.Error("error reading body", zap.Error(err))
	}

	logging.Logger.Info("body", zap.String("body", string(body)))
}

func (a *App) getInstallations(c *gin.Context) {
	installations, err := a.githubService.GetInstallationsForApp()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, installations)
}

func (a *App) job() error {

	for range time.Tick(time.Minute * 10) {
		installations, err := a.githubService.GetInstallationsForApp()
		if err != nil {
			return err
		}

		logging.Logger.Info("installations", zap.Any("installations", installations))

	}
	return nil
}
