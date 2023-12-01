package app

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/TonyDMorris/quick-function/constants"
	gptModels "github.com/TonyDMorris/quick-function/pkg/gpt/models"
	"github.com/TonyDMorris/quick-function/pkg/logging"
	"github.com/TonyDMorris/quick-function/pkg/strapi/models"
	strapiModels "github.com/TonyDMorris/quick-function/pkg/strapi/models"
	"github.com/gin-gonic/gin"
	"github.com/google/go-github/v56/github"
	"go.uber.org/zap/zapcore"
)

func (a *App) HandleStrapiWebhook(c *gin.Context) {

	var webhook models.StrapiWebhookPayload

	if err := c.BindJSON(&webhook); err != nil {
		logging.Logger.Error("Error handling repository configuration", zapcore.Field{Key: "error", Type: zapcore.ErrorType, Interface: err})
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	switch webhook.Model {
	case "repository-configuration":
		if err := a.handleRepositoryConfiguration(webhook); err != nil {
			logging.Logger.Error("Error handling repository configuration", zapcore.Field{Key: "error", Type: zapcore.ErrorType, Interface: err})
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}

	default:
		c.JSON(200, gin.H{
			"success": "true", "message": "ok",
		})
	}

	c.JSON(200, gin.H{
		"success": "true", "message": "ok",
	})
}

func (a *App) handleRepositoryConfiguration(webhook models.StrapiWebhookPayload) error {

	switch webhook.Event {
	case "entry.create":
		bytes, err := json.Marshal(webhook.Entry)
		if err != nil {
			return fmt.Errorf("error marshalling entry: %w", err)
		}

		var repositoryConfiguration models.RepositoryConfiguration

		if err = json.Unmarshal(bytes, &repositoryConfiguration); err != nil {
			return fmt.Errorf("error unmarshalling entry: %w", err)
		}

		fullRepositoryConfiguration, err := a.strapiClient.GetRepositoryConfiguration(repositoryConfiguration.ID)
		if err != nil {
			return fmt.Errorf("error getting repository configuration: %w", err)
		}

		a.WorkerPool.RepostioryConfigurationCreated <- *fullRepositoryConfiguration

		return nil

		return nil
	default:
		return nil
	}

}

func (a *App) HandleRepositoryConfigurationCreatedJob(job strapiModels.RepositoryConfiguration) error {

	ctx := context.Background()

	installation := job.Installation
	repo := job.Repository

	if installation == nil || repo == nil {
		return fmt.Errorf("installation or repository is nil")
	}

	installationID := installation.InstallationID

	userClient, err := a.getUserClientFromInstallation(ctx, installationID)
	if err != nil {
		return fmt.Errorf("error getting user client from installation: %w", err)
	}

	repoinfo, _, err := userClient.Repositories.Get(ctx, installation.Username, repo.Name)
	if err != nil {
		return fmt.Errorf("error getting repository info: %w", err)
	}

	defaultBranch := repoinfo.GetDefaultBranch()

	tree, _, err := userClient.Git.GetTree(ctx, installation.Username, repo.Name, defaultBranch, true)
	if err != nil {
		return fmt.Errorf("error getting tree: %w", err)
	}

	var files []string

	for _, entry := range tree.Entries {
		if entry.GetSize() != 0 {
			files = append(files, entry.GetPath())
		}
	}

	if len(files) == 0 {
		return fmt.Errorf("no files found")
	}

}

func (a *App) getInterestedFiles(repoName string) ([]string, error) {
	intestestFilesPrompts := []gptModels.Message{
		{
			Role:    gptModels.RoleSystem,
			Content: constants.InterestedFiles,
		},
		{
			Role:    gptModels.RoleUser,
			Content: fmt.Sprintf("%s", repoName),
		},
	}

	resp, err := a.chatGptClient.Chat(intestestFilesPrompts)

	if err != nil {
		return nil, fmt.Errorf("error chatting with gpt: %w", err)
	}

	var files []string

	for _, choice := range resp.Choices {

	}

}

func (a *App) getUserClientFromInstallation(ctx context.Context, installationID string) (*github.Client, error) {

	instID, err := strconv.ParseInt(installationID, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("error converting installation id to int: %w", err)
	}

	installationToken, _, err := a.githubClient.Apps.CreateInstallationToken(ctx, instID, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating installation token: %w", err)
	}

	return github.NewClient(http.DefaultClient).WithAuthToken(installationToken.GetToken()), nil

}
