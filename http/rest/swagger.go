package rest

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/ghodss/yaml"
	swagger "github.com/swaggo/http-swagger"
)

func registerSwaggerHandlers(mux *http.ServeMux, schemaDir string) error {
	// Load the OpenAPI schema from YAML file
	schemaYAML, err := os.ReadFile(schemaDir)
	if err != nil {
		return err
	}

	var schema map[string]interface{}
	err = yaml.Unmarshal(schemaYAML, &schema)
	if err != nil {
		return err
	}

	// Serve the OpenAPI specification in JSON format
	mux.HandleFunc("GET /openapi/schema.json", func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(schema)
	})

	// Serve the Swagger UI (e.g. http://localhost:8080/openapi/swagger/index.html)
	mux.HandleFunc("GET /openapi/swagger/*", swagger.Handler(
		swagger.URL("/openapi/schema.json"),
	))

	return nil
}
