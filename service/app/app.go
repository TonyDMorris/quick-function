package app

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	gpt "github.com/TonyDMorris/quick-function/pkg/gpt/client"
	"github.com/gin-gonic/gin"

	"github.com/google/go-github/v56/github"
)

type Config struct {
	Port int
}

type App struct {
	server        *gin.Engine
	githubClient  *github.Client
	chatGptClient *gpt.ChatClient
	port          int
}

func NewApi(c Config, githubClient *github.Client, gptClient *gpt.ChatClient) *App {
	return &App{
		server: gin.Default(),

		githubClient:  githubClient,
		chatGptClient: gptClient,
		port:          c.Port,
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
