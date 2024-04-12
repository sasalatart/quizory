package migrations

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/pkg/errors"
)

// Up runs the migrations on the database specified by dbUrl.
func Up(dbUrl string) error {
	moduleRoot, err := getModuleRoot()
	if err != nil {
		return errors.Wrap(err, "getting module root")
	}
	m, err := migrate.New(fmt.Sprintf("file://%s/db/migrations", moduleRoot), dbUrl)
	if err != nil {
		return errors.Wrap(err, "creating migrations instance")
	}
	if err := m.Up(); err != nil {
		return errors.Wrap(err, "executing migrations")
	}
	return nil
}

// getModuleRoot returns the root directory of the module.
func getModuleRoot() (string, error) {
	cmd := exec.Command("go", "list", "-m", "-f", "{{.Dir}}")
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return strings.TrimSpace(stdout.String()), nil
}
