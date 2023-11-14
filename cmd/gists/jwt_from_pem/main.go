package main

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

func main() {

	bytes, err := os.ReadFile("../../../secrets/git-history-blog.2023-11-14.private-key.pem")
	if err != nil {

		panic(err)

	}
	signKey, err := jwt.ParseRSAPrivateKeyFromPEM(bytes)
	if err != nil {
		panic(err)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"iat": time.Now().Unix() - 60,
		"exp": time.Now().Add(time.Minute * 10).Unix(),
		"iss": 467659,
	})

	tokenString, err := token.SignedString(signKey)
	if err != nil {
		panic(err)
	}

	println(tokenString)

}
