package client

import (
	"crypto/rsa"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/hashicorp/go-retryablehttp"
)

type TokenSerivce struct {
	signKey *rsa.PrivateKey
	appID   int
	client  *http.Client
}

func NewTokenSerivce(signKey []byte, appID int) *TokenSerivce {
	key, err := jwt.ParseRSAPrivateKeyFromPEM(signKey)
	if err != nil {
		panic(err)
	}

	return &TokenSerivce{
		signKey: key,
		appID:   appID,
		client:  retryablehttp.NewClient().StandardClient(),
	}
}

func (t *TokenSerivce) getToken() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"iat": time.Now().Unix() - 60,
		"exp": time.Now().Add(time.Minute * 10).Unix(),
		"iss": t.appID,
	})

	tokenString, err := token.SignedString(t.signKey)
	if err != nil {
		return "", fmt.Errorf("error signing token: %w", err)
	}
	return tokenString, nil
}

func (t *TokenSerivce) GetToken(activeToken *string) (string, error) {
	if activeToken == nil {
		return t.getToken()
	}
	token, err := jwt.Parse(*activeToken, func(token *jwt.Token) (interface{}, error) {
		return t.signKey, nil
	})
	if err != nil {
		return "", err
	}
	exp := token.Claims.(jwt.MapClaims)["exp"].(int)
	// 2 minutes before expiration
	if time.Now().Unix()+2*60 > int64(exp) {
		return t.getToken()
	}
	return *activeToken, nil
}
