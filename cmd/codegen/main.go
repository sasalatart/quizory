package main

import (
	"log"
	"log/slog"
	"net"
	"net/url"
	"os"
	"os/exec"
	"strings"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	slog.Info("Running code generation...")
	defer slog.Info("Code generation complete.")

	u, err := url.Parse(os.Getenv("DB_URL"))
	if err != nil {
		log.Fatal(err)
	}

	psqlUser := u.User.Username()
	psqlPass, _ := u.User.Password()
	psqlHost, psqlPort, err := net.SplitHostPort(u.Host)
	if err != nil {
		log.Fatal(err)
	}

	os.Setenv("PSQL_USER", psqlUser)
	os.Setenv("PSQL_PASS", psqlPass)
	os.Setenv("PSQL_HOST", psqlHost)
	os.Setenv("PSQL_PORT", psqlPort)
	os.Setenv("PSQL_DBNAME", strings.TrimPrefix(u.Path, "/"))

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
