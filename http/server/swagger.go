package server

import (
	"log"
	"os"

	"github.com/ghodss/yaml"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	swagger "github.com/swaggo/http-swagger"
)

func (s *Server) mustRegisterSwaggerHandlers() {
	// Load the OpenAPI schema from YAML file
	schemaYAML, err := os.ReadFile(s.cfg.SchemaDir)
	if err != nil {
		log.Fatal(err)
	}

	var schema map[string]interface{}
	err = yaml.Unmarshal(schemaYAML, &schema)
	if err != nil {
		log.Fatal(err)
	}

	// Serve the OpenAPI specification in JSON format
	s.App.Get("/openapi/schema.json", func(c *fiber.Ctx) error {
		return c.JSON(schema)
	})

	// Serve the Swagger UI (e.g. http://localhost:8080/openapi/swagger/index.html)
	s.App.Get("/openapi/swagger/*", adaptor.HTTPHandler(swagger.Handler(
		swagger.URL("/openapi/schema.json"),
	)))
}
