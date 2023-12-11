package client

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/TonyDMorris/quick-function/pkg/strapi/models"
	"github.com/hashicorp/go-retryablehttp"
)

const repositoryConfigurationPath = "%s/api/internal/repository-configurations/%d"
const repositoryConfigurationsPath = "%s/api/internal/repository-configurations"
const gitBlogPostsPath = "%s/api/git-blog-posts"

type Client struct {
	apiKey         string
	baseURL        string
	retryingClient *retryablehttp.Client
}

func NewClient(apiKey, baseURL string) *Client {
	retryingClient := retryablehttp.NewClient()
	return &Client{
		apiKey:         apiKey,
		baseURL:        baseURL,
		retryingClient: retryingClient,
	}
}

func (c *Client) GetRepositoryConfiguration(id int) (*models.RepositoryConfiguration, error) {
	req, err := retryablehttp.NewRequest(http.MethodGet, fmt.Sprintf(repositoryConfigurationPath, c.baseURL, id), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.retryingClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var repositoryConfiguration models.RepositoryConfiguration

	err = json.NewDecoder(resp.Body).Decode(&repositoryConfiguration)
	if err != nil {
		return nil, err
	}

	return &repositoryConfiguration, nil

}

func (c *Client) GetRepositoryConfigurations() ([]models.RepositoryConfiguration, error) {
	req, err := retryablehttp.NewRequest(http.MethodGet, fmt.Sprintf(repositoryConfigurationsPath, c.baseURL), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.retryingClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var repositoryConfigurations []models.RepositoryConfiguration

	err = json.NewDecoder(resp.Body).Decode(&repositoryConfigurations)
	if err != nil {
		return nil, err
	}

	return repositoryConfigurations, nil

}

func (c *Client) UpdateRepositoryConfiguration(repoConfig models.RepositoryConfiguration) (*models.RepositoryConfiguration, error) {
	carrier := models.Carrier{
		Data: repoConfig,
	}
	body, err := json.Marshal(carrier)
	if err != nil {
		return nil, err
	}
	req, err := retryablehttp.NewRequest(http.MethodPut, fmt.Sprintf(repositoryConfigurationPath, c.baseURL, repoConfig.ID), body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.retryingClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var repositoryConfiguration models.RepositoryConfiguration

	err = json.NewDecoder(resp.Body).Decode(&repositoryConfiguration)
	if err != nil {
		return nil, err
	}

	return &repositoryConfiguration, nil

}

func (c *Client) StandardCreateGitBlogPost(gitBlogPost models.GitBlogPost) (*models.GitBlogPost, error) {
	carrier := models.Carrier{
		Data: gitBlogPost,
	}

	body, err := json.Marshal(carrier)
	if err != nil {
		return nil, err
	}
	req, err := retryablehttp.NewRequest(http.MethodPost, fmt.Sprintf(gitBlogPostsPath, c.baseURL), body)
	if err != nil {
		return nil, err

	}

	resp, err := c.retryingClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var responseCarrier models.Carrier

	err = json.NewDecoder(resp.Body).Decode(&responseCarrier)

	if err != nil {
		return nil, err
	}

	bytes, err := json.Marshal(responseCarrier.Data)
	if err != nil {
		return nil, err
	}
	var respGitBlogPost models.GitBlogPost
	err = json.Unmarshal(bytes, &respGitBlogPost)
	if err != nil {
		return nil, err
	}

	return &gitBlogPost, nil

}
