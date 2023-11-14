package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hashicorp/go-retryablehttp"
)

type Config struct {
	Port                  int
	GithubAppClientID     string
	GithubAppClientSecret string
}

type Api struct {
	config Config
	server *gin.Engine
	client *http.Client
}

func NewApi(c Config) *Api {
	return &Api{
		config: c,
		server: gin.Default(),
		client: retryablehttp.NewClient().StandardClient(),
	}
}

func (a *Api) Run() error {
	return a.server.Run()
}
