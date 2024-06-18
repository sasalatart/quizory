package main

import (
	"log"
	"log/slog"
	"os"
	"os/exec"

	"github.com/sasalatart/quizory/config"
)

func main() {
	slog.Info("Running code generation...")
	defer slog.Info("Code generation complete.")

	dbCfg := config.NewConfig().DB

	os.Setenv("PSQL_USER", dbCfg.User)
	os.Setenv("PSQL_PASS", dbCfg.Password)
	os.Setenv("PSQL_HOST", dbCfg.Host)
	os.Setenv("PSQL_PORT", dbCfg.Port)
	os.Setenv("PSQL_DBNAME", dbCfg.Name)

	sqlboilerCmd := exec.Command("sqlboiler", "psql", "-c", "./cmd/codegen/sqlboiler.toml")
	sqlboilerCmd.Stdout = os.Stdout
	sqlboilerCmd.Stderr = os.Stderr
	if err := sqlboilerCmd.Run(); err != nil {
		log.Fatal(err)
	}

	codegenCmd := exec.Command("go", "generate", "./...")
	codegenCmd.Stdout = os.Stdout
	codegenCmd.Stderr = os.Stderr
	if err := codegenCmd.Run(); err != nil {
		log.Fatal(err)
	}
}
