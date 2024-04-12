package main

import (
	"log"
	"net"
	"net/url"
	"os"
	"os/exec"
	"strings"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

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

	cmd := exec.Command("sqlboiler", "psql", "-c", "./cmd/codegen/sqlboiler.toml")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}
