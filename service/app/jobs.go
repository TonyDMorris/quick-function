package app

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/TonyDMorris/quick-function/constants"
	gpt "github.com/TonyDMorris/quick-function/pkg/gpt/client"
	gptModels "github.com/TonyDMorris/quick-function/pkg/gpt/models"
	"github.com/TonyDMorris/quick-function/pkg/logging"
	strapiModels "github.com/TonyDMorris/quick-function/pkg/strapi/models"
	"github.com/google/go-github/v56/github"
	"github.com/pkoukk/tiktoken-go"
	"go.uber.org/zap/zapcore"
)

func (a *App) HandleRepositoryConfigurationCreatedJob(job strapiModels.RepositoryConfiguration) error {

	defer func() {
		if err := recover(); err != nil {
			logging.Logger.Error("panic in HandleRepositoryConfigurationCreatedJob", zapcore.Field{Key: "error", Type: zapcore.ErrorType, Interface: err}, zapcore.Field{Key: "job", Type: zapcore.Int64Type, Interface: job.ID})
		}
	}()

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

	interestedFiles, err := a.getInterestedFiles(repo.Name, files)

	if err != nil {
		return fmt.Errorf("error getting interested files: %w", err)
	}

	if len(interestedFiles) == 0 {
		return fmt.Errorf("no interested files found")
	}
	contents, err := a.getContents(ctx, userClient, installation.Username, repo.Name, interestedFiles)
	if err != nil {
		return fmt.Errorf("error getting contents: %w", err)
	}

	tokensPerContent := gpt.MaxTokens / len(contents)

	var trimmedContents = make(map[string]string)

	for path, content := range contents {
		trimmedContent, err := a.trimContent(content, tokensPerContent)
		if err != nil {
			return fmt.Errorf("error trimming content: %w", err)
		}

		trimmedContents[path] = trimmedContent

	}

	var contentsToSend []string

	for path, content := range trimmedContents {
		contentsToSend = append(contentsToSend, fmt.Sprintf("%s\n%s", path, content))
	}

	contentsToSendString := strings.Join(contentsToSend, "\n")

	contentMessage := fmt.Sprintf(constants.ContentMessage, contentsToSendString)

	contentMessagePrompts := []gptModels.Message{
		{
			Role:    gptModels.RoleSystem,
			Content: contentMessage,
		},
	}

	resp, err := a.chatGptClient.Chat(contentMessagePrompts)

	if err != nil {
		return fmt.Errorf("error chatting with gpt: %w", err)
	}

	var content string

	for _, choice := range resp.Choices {
		content = choice.Message.Content

	}

	logging.Logger.Info(content)

	// trim content

	return nil

}

func (a *App) HandleRepositoryConfigurationScheduledJob(job strapiModels.RepositoryConfiguration) error {

	defer func() {
		if err := recover(); err != nil {
			logging.Logger.Error("panic in HandleRepositoryConfigurationScheduledJob", zapcore.Field{Key: "error", Type: zapcore.ErrorType, Interface: err}, zapcore.Field{Key: "job", Type: zapcore.Int64Type, Interface: job.ID})
		}
	}()
	return nil
}

func (a *App) trimContent(content string, tokens int) (string, error) {
	tikToken, err := tiktoken.GetEncoding("cl100k_base")
	if err != nil {
		logging.Logger.Error("error getting tik token", zapcore.Field{Key: "error", Type: zapcore.ErrorType, Interface: err})
		return "", fmt.Errorf("error getting tik token: %w", err)
	}

	contentTokens := tikToken.Encode(content, nil, nil)

	if len(contentTokens) < tokens {
		return content, nil
	}

	overHangPercentage := float64(len(contentTokens)-tokens) / float64(len(contentTokens))

	// slice string based on percentage

	return content[:int(float64(len(content))*overHangPercentage)/100], nil
}

func (a *App) getContents(ctx context.Context, userClient *github.Client, username string, repo string, interestedFiles []string) (map[string]string, error) {
	var contents = make(map[string]string)
	for _, path := range interestedFiles {
		content, _, _, err := userClient.Repositories.GetContents(ctx, username, repo, path, nil)
		if err != nil {
			logging.Logger.Error("error getting contents", zapcore.Field{Key: "error", Type: zapcore.ErrorType, Interface: err})
			continue
		}

		if content == nil {
			logging.Logger.Warn("content is nil")
			continue
		}

		contentBytes, err := content.GetContent()
		if err != nil {
			logging.Logger.Warn("error getting content bytes", zapcore.Field{Key: "error", Type: zapcore.ErrorType, Interface: err})
			continue
		}

		if len(contentBytes) == 0 {
			logging.Logger.Warn("content bytes is empty")
			continue
		}

		minifiedContent := strings.ReplaceAll(contentBytes, "\n", "")
		minifiedContent = strings.ReplaceAll(minifiedContent, "\t", "")
		minifiedContent = strings.ReplaceAll(minifiedContent, " ", "")

		contents[path] = minifiedContent

	}
	if len(contents) == 0 {
		return nil, fmt.Errorf("no contents found")
	}

	return contents, nil

}

func (a *App) getInterestedFiles(repoName string, allFiles []string) ([]string, error) {
	fileToSend := strings.Join(allFiles, "\n")
	intestestFilesPrompts := []gptModels.Message{
		{
			Role:    gptModels.RoleSystem,
			Content: constants.InterestedFiles,
		},
		{
			Role:    gptModels.RoleUser,
			Content: fileToSend,
		},
	}

	resp, err := a.chatGptClient.Chat(intestestFilesPrompts)

	if err != nil {
		return nil, fmt.Errorf("error chatting with gpt: %w", err)
	}

	var files []string

	for _, choice := range resp.Choices {
		files = append(files, strings.Split(choice.Message.Content, "\n")...)
	}

	logging.Logger.Info(strings.Join(files, "\n"))
	return files, nil

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
