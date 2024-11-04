package resttest

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/sasalatart/quizory/config"
	"github.com/sasalatart/quizory/http/rest/oapi"
)

// ClientFactory creates an authenticated oapi.ClientWithResponses for a given user ID.
type ClientFactory func(userID uuid.UUID) (*oapi.ClientWithResponses, error)

func newClientFactory(cfg config.ServerConfig) ClientFactory {
	return func(userID uuid.UUID) (*oapi.ClientWithResponses, error) {
		authUserClientOption := func(c *oapi.Client) error {
			c.RequestEditors = append(
				c.RequestEditors,
				func(ctx context.Context, req *http.Request) error {
					req.Header.Set(
						"Authorization",
						fmt.Sprintf("Bearer %s", newJWT(userID, cfg.JWTSecret)),
					)
					return nil
				})
			return nil
		}

		return oapi.NewClientWithResponses(
			fmt.Sprintf("http://%s", cfg.RESTAddress()),
			authUserClientOption,
		)
	}
}

func newJWT(userID uuid.UUID, secret string) string {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["sub"] = userID
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix() // Expiration time (24 hours from now)

	// Sign the token with a secret key
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		log.Fatal(err)
	}
	return tokenString
}
