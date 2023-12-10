package app

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/TonyDMorris/quick-function/pkg/logging"
	"github.com/TonyDMorris/quick-function/pkg/strapi/models"
	"github.com/gin-gonic/gin"
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
		c.JSON(204, gin.H{
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

		scheduleStrings := strings.Split(repositoryConfiguration.Cron, " ")
		numberString, interval := scheduleStrings[0], scheduleStrings[1]
		number, err := strconv.Atoi(numberString)
		if err != nil {
			return fmt.Errorf("error converting string to number: %w", err)
		}
		switch interval {
		case "days":
			job, err := a.cron.
				Every(uint64(number)).
				Day().
				Tag(fmt.Sprint(fullRepositoryConfiguration.ID)).
				WaitForSchedule().
				Do(func() {
					a.WorkerPool.RepostioryConfigurationCreated <- *fullRepositoryConfiguration
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
				Every(uint64(number)).
				Week().
				Tag(fmt.Sprint(fullRepositoryConfiguration.ID)).
				WaitForSchedule().
				Do(func() {
					a.WorkerPool.RepostioryConfigurationCreated <- *fullRepositoryConfiguration
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

		a.WorkerPool.RepostioryConfigurationCreated <- *fullRepositoryConfiguration

		return nil

	default:
		return nil
	}

}
