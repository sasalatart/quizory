package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/sasalatart.com/quizory/config"
	"github.com/sasalatart.com/quizory/http/oapi"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"server",
	fx.Provide(NewServer),
)

var TestModule = fx.Module(
	"test-server",
	fx.Provide(NewServer),
	fx.Provide(newTestClient),
)

func newTestClient(cfg config.ServerConfig) *oapi.ClientWithResponses {
	client, err := oapi.NewClientWithResponses(
		fmt.Sprintf("http://%s", cfg.Address()),
		func(c *oapi.Client) error {
			c.RequestEditors = append(
				c.RequestEditors,
				func(ctx context.Context, req *http.Request) error {
					req.Header.Set(
						"Authorization",
						fmt.Sprintf("Bearer %s", newTestJWT(cfg.JWTSecret)),
					)
					return nil
				})
			return nil
		},
	)
	if err != nil {
		log.Fatal(err)
	}
	return client
}

func newTestJWT(secret string) string {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["sub"] = uuid.New()
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix() // Expiration time (24 hours from now)

	// Sign the token with a secret key
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		log.Fatal(err)
	}
	return tokenString
}
