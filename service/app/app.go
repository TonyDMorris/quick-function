package app

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	gpt "github.com/TonyDMorris/quick-function/pkg/gpt/client"
	"github.com/TonyDMorris/quick-function/pkg/logging"
	strapi "github.com/TonyDMorris/quick-function/pkg/strapi/client"
	strapiModels "github.com/TonyDMorris/quick-function/pkg/strapi/models"
	"github.com/gin-gonic/gin"
	"github.com/go-co-op/gocron"
	"go.uber.org/zap/zapcore"
	"golang.org/x/sync/errgroup"

	"github.com/google/go-github/v56/github"
)

type Config struct {
	Port int
}

type App struct {
	server        *gin.Engine
	githubClient  *github.Client
	chatGptClient *gpt.ChatClient
	strapiClient  *strapi.Client
	cron          *gocron.Scheduler
	jobs          map[string]*gocron.Job
	port          int
	WorkerPool    *WorkerPool
}

type WorkerPool struct {
	RepostioryConfigurationCreated   chan strapiModels.RepositoryConfiguration
	RepositoryConfigurationScheduled chan strapiModels.RepositoryConfiguration
}

func (a *App) startWorkerPool() {

	go func() {
		group := errgroup.Group{}
		group.SetLimit(10)
		for repoCreateJob := range a.WorkerPool.RepostioryConfigurationCreated {
			repoCreateJob := repoCreateJob
			group.Go(func() error {
				err := a.HandleRepositoryConfigurationCreatedJob(repoCreateJob)
				if err != nil {
					logging.Logger.Error("error handling repository configuration created job", zapcore.Field{Key: "error", Type: zapcore.ErrorType, Interface: err})
				}
				return nil
			})
		}
	}()
	go func() {
		group := errgroup.Group{}
		group.SetLimit(10)
		for repoCreateJob := range a.WorkerPool.RepositoryConfigurationScheduled {
			repoCreateJob := repoCreateJob
			group.Go(func() error {
				err := a.HandleRepositoryConfigurationScheduledJob(repoCreateJob)
				if err != nil {
					logging.Logger.Error("error handling repository configuration created job", zapcore.Field{Key: "error", Type: zapcore.ErrorType, Interface: err})
				}
				return nil
			})
		}
	}()

}

func (a *App) loadSchedules() error {
	repositoryConfigurations, err := a.strapiClient.GetRepositoryConfigurations()
	if err != nil {
		return fmt.Errorf("error getting repository configurations: %w", err)
	}

	for _, repositoryConfiguration := range repositoryConfigurations {

		if repositoryConfiguration.Cron == "" || repositoryConfiguration.NextGeneration == nil {
			continue
		}
		logging.Logger.Info(fmt.Sprintf("scheduling job for repository configuration: %d", repositoryConfiguration.ID))
		fullRepositoryConfiguration, err := a.strapiClient.GetRepositoryConfiguration(repositoryConfiguration.ID)
		if err != nil {
			return fmt.Errorf("error getting repository configuration: %w", err)
		}

		scheduleStrings := strings.Split(repositoryConfiguration.Cron, " ")
		numberString, interval := scheduleStrings[0], scheduleStrings[1]
		number, err := strconv.Atoi(numberString)
		if err != nil {
			return fmt.Errorf("error converting string to number: %w", err)
		}
		switch interval {
		case "days":
			job, err := a.cron.
				Every(int(number)).
				Day().
				Tag(fmt.Sprint(fullRepositoryConfiguration.ID)).
				StartAt(*fullRepositoryConfiguration.NextGeneration).
				Do(func() {
					a.WorkerPool.RepositoryConfigurationScheduled <- *fullRepositoryConfiguration
				})
			if err != nil {
				return fmt.Errorf("error scheduling job: %w", err)
			}
			a.jobs[fmt.Sprint(fullRepositoryConfiguration.ID)] = job
			nextRun := job.NextRun()
			repositoryConfiguration.NextGeneration = &nextRun
			_, err = a.strapiClient.UpdateRepositoryConfiguration(repositoryConfiguration)
			if err != nil {
				return fmt.Errorf("error updating repository configuration: %w", err)
			}

		case "weeks":
			job, err := a.cron.
				Every(int(number)).
				Week().
				Tag(fmt.Sprint(fullRepositoryConfiguration.ID)).
				StartAt(*fullRepositoryConfiguration.NextGeneration).
				Do(func() {
					a.WorkerPool.RepositoryConfigurationScheduled <- *fullRepositoryConfiguration
				})
			if err != nil {
				return fmt.Errorf("error scheduling job: %w", err)
			}
			a.jobs[fmt.Sprint(fullRepositoryConfiguration.ID)] = job
			nextRun := job.NextRun()
			repositoryConfiguration.NextGeneration = &nextRun
			_, err = a.strapiClient.UpdateRepositoryConfiguration(repositoryConfiguration)
			if err != nil {
				return fmt.Errorf("error updating repository configuration: %w", err)
			}

		default:
			return fmt.Errorf("invalid interval: %s", interval)

		}

	}

	return nil
}

func (a *App) HandleGetJobs(c *gin.Context) {
	var resp []map[string]interface{}

	for _, job := range a.jobs {
		resp = append(resp, map[string]interface{}{
			"next_run": job.NextRun(),
			"tags":     job.Tags(),
			"name":     job.GetName(),
		})
	}

	c.JSON(http.StatusOK, resp)
}

func NewApi(c Config, githubClient *github.Client, gptClient *gpt.ChatClient, strapiClient *strapi.Client) *App {
	return &App{
		server: gin.Default(),

		githubClient:  githubClient,
		chatGptClient: gptClient,
		strapiClient:  strapiClient,
		port:          c.Port,
		cron:          gocron.NewScheduler(time.UTC),
		jobs:          make(map[string]*gocron.Job),
		WorkerPool: &WorkerPool{
			RepostioryConfigurationCreated: make(chan strapiModels.RepositoryConfiguration),
		},
	}
}

var num = 1

func (a *App) Test2(c *gin.Context) {
	body, _ := io.ReadAll(c.Request.Body)
	path := c.Request.URL.Path
	query := c.Request.URL.Query()
	newBody := make(map[string]interface{})
	newBody["body"] = string(body)
	newBody["path"] = path
	newBody["query"] = query

	fileName := fmt.Sprintf("%d.json", num)
	num++
	file, err := os.Create(fileName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	defer file.Close()

	_, err = file.Write(body)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

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

		repos, _, err := userClient.Apps.ListRepos(c.Request.Context(), &github.ListOptions{
			Page:    1,
			PerPage: 100,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		for _, repo := range repos.Repositories {
			if repo.GetName() != "discwordle" {
				continue
			}
			tree, _, err := userClient.Git.GetTree(c.Request.Context(), repo.GetOwner().GetLogin(), repo.GetName(), "master", true)
			if err != nil {
				continue
			}
			file, err := os.Create(fmt.Sprintf("./data/%s-tree.txt", repo.GetName()))

			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})

			}

			for _, entry := range tree.Entries {
				if entry.GetSize() != 0 {
					file.WriteString(fmt.Sprintf("%s:%d", entry.GetPath(), entry.GetSize()))
				}
			}
			defer file.Close()

			list := []string{"README.md:3359",
				"src/App.js:4149",
				"src/components/guesses.jsx:794",
				"src/components/header.jsx:1054",
				"src/components/input.jsx:2603",
				"src/components/quoteCarosel.jsx:1402",
				"src/components/result.jsx:1674",
				"src/data/main.go:2790",
				"src/data/quotes.json:181032",
				"src/index.js:535"}

			contentsFile, err := os.Create(fmt.Sprintf("./data/%s-contents.txt", repo.GetName()))

			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})

			}

			defer contentsFile.Close()

			for _, entry := range list {
				path := strings.Split(entry, ":")[0]

				contentRdr, _, err := userClient.Repositories.DownloadContents(c.Request.Context(), repo.GetOwner().GetLogin(), repo.GetName(), path, nil)

				if err != nil {

					c.JSON(http.StatusInternalServerError, gin.H{
						"error": err.Error(),
					})
					return
				}

				content, err := io.ReadAll(contentRdr)

				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{
						"error": err.Error(),
					})
					return
				}

				minifiedContent := strings.ReplaceAll(string(content), "\n", "")
				minifiedContent = strings.ReplaceAll(minifiedContent, "\t", "")
				minifiedContent = strings.ReplaceAll(minifiedContent, " ", "")

				contentsFile.WriteString(fmt.Sprintf("%s\n%s", path, minifiedContent))
			}
		}
	}

}
