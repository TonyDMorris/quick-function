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
	body, err := json.Marshal(repoConfig)
	if err != nil {
		return nil, err
	}
	req, err := retryablehttp.NewRequest(http.MethodPut, fmt.Sprintf(repositoryConfigurationPath, c.baseURL, repoConfig.ID), body)
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
