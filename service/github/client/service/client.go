package service

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/TonyDMorris/quick-function/service/github/client/models"
	"github.com/hashicorp/go-retryablehttp"
)

const (
	BaseURL           = "https://api.github.com"
	InstallationsPath = "/app/installations"
)

type tokenService interface {
	GetToken(*string) (string, error)
}
type Service struct {
	tokenSerivce tokenService
	activeToken  *string
	httpClient   *http.Client
}

func NewService(tokenService tokenService) *Service {
	return &Service{
		tokenSerivce: tokenService,
		httpClient:   retryablehttp.NewClient().StandardClient(),
	}
}

func (s *Service) checkToken() error {
	token, err := s.tokenSerivce.GetToken(s.activeToken)
	if err != nil {
		return err
	}
	s.activeToken = &token
	return nil
}

func (s *Service) GetInstallationsForApp() ([]models.AppInstallation, error) {
	if err := s.checkToken(); err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodGet, BaseURL+InstallationsPath, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+*s.activeToken)
	req.Header.Add("Accept", "application/vnd.github+json")
	req.Header.Add("X-GitHub-Api-Version", "2022-11-28")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var installations []models.AppInstallation
	if err := json.Unmarshal(bytes, &installations); err != nil {
		return nil, err
	}
	return installations, nil
}
