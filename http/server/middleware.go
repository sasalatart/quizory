package server

import (
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/logger"

	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/sasalatart.com/quizory/config"
)

// applyDefaultMiddleware applies default middleware to an app according to a given server config.
func applyDefaultMiddleware(app *fiber.App, cfg config.ServerConfig) {
	app.Use(recover.New(recover.Config{
		EnableStackTrace: true,
	}))

	app.Use(helmet.New())
	app.Use(cors.New())

	app.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(cfg.JWTSecret)},
	}))

	app.Use(logger.New())
}

// GetUserID returns the ID of the user that made a request.
func GetUserID(c *fiber.Ctx) (uuid.UUID, error) {
	user := c.Locals("user").(*jwt.Token)
	userID := user.Claims.(jwt.MapClaims)["sub"].(string)
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return uuid.Nil, errors.Wrapf(err, "parsing user ID %s", userID)
	}
	return userUUID, nil
}
