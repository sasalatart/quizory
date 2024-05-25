package server

import (
	"os"
	"strings"

	"github.com/ghodss/yaml"
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/sasalatart.com/quizory/http/oapi"
	swagger "github.com/swaggo/http-swagger"

	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

// registerMiddlewares applies default middleware to an app according to a given server config.
func (s *Server) registerMiddlewares() {
	s.App.Use(recover.New(recover.Config{
		EnableStackTrace: true,
	}))

	s.App.Use(helmet.New())
	s.App.Use(cors.New())

	s.App.Use(func(c *fiber.Ctx) error {
		if strings.HasPrefix(c.Path(), "/openapi") {
			return c.Next()
		}
		return jwtware.New(jwtware.Config{
			SigningKey: jwtware.SigningKey{Key: []byte(s.cfg.JWTSecret)},
		})(c)
	})

	s.App.Use(logger.New())
}

func (s *Server) registerAppHandlers() error {
	if err := s.registerSwaggerHandlers(); err != nil {
		return err
	}
	oapi.RegisterHandlers(s.App, s)
	return nil
}

func (s *Server) registerSwaggerHandlers() error {
	// Load the OpenAPI schema from YAML file
	schemaYAML, err := os.ReadFile(s.cfg.SchemaDir)
	if err != nil {
		return err
	}

	var schema map[string]interface{}
	err = yaml.Unmarshal(schemaYAML, &schema)
	if err != nil {
		return err
	}

	// Serve the OpenAPI specification in JSON format
	s.App.Get("/openapi/schema.json", func(c *fiber.Ctx) error {
		return c.JSON(schema)
	})

	// Serve the Swagger UI (e.g. http://localhost:8080/openapi/swagger/index.html)
	s.App.Get("/openapi/swagger/*", adaptor.HTTPHandler(swagger.Handler(
		swagger.URL("/openapi/schema.json"),
	)))

	return nil
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
