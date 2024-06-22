package server

import (
	"log/slog"
	"os"
	"strings"
	"time"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"

	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

func init() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))
}

// registerMiddlewares applies default middleware to an app according to a given server config.
func (s *Server) registerMiddlewares() {
	s.App.Use(recover.New(recover.Config{
		EnableStackTrace: true,
	}))

	s.App.Use(helmet.New())
	s.App.Use(cors.New())
	s.App.Use(newJWTMiddleware(s.cfg.JWTSecret))
	s.App.Use(newLoggerMiddleware())
}

func newJWTMiddleware(jwtSecret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if strings.HasPrefix(c.Path(), "/openapi") || strings.HasPrefix(c.Path(), "/health-check") {
			return c.Next()
		}
		return jwtware.New(jwtware.Config{
			SigningKey: jwtware.SigningKey{Key: []byte(jwtSecret)},
		})(c)
	}
}

func newLoggerMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		slog.Info(
			"Handling request",
			slog.Group(
				"http",
				slog.String("method", c.Method()),
				slog.String("path", c.Path()),
				slog.String("ip", c.IP()),
			),
		)

		err := c.Next()

		slog.Info(
			"Completed request",
			slog.Group(
				"http",
				slog.Int("status", c.Response().StatusCode()),
				slog.String("duration", time.Since(start).String()),
			),
		)

		return err
	}
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
