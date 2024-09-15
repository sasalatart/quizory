package main

import (
	"log"
	"net"
	"net/url"

	"github.com/sasalatart/quizory/config"
	"github.com/volatiletech/sqlboiler/v4/boilingcore"
	"github.com/volatiletech/sqlboiler/v4/importers"

	_ "github.com/volatiletech/sqlboiler/v4/drivers/sqlboiler-psql/driver"
)

func runSQLBoiler() {
	dbCfg := config.NewConfig().DB

	u, err := url.Parse(dbCfg.URL)
	if err != nil {
		log.Fatal("unable to parse connection URL: ", err)
	}

	username := u.User.Username()
	password, _ := u.User.Password()

	host, port, err := net.SplitHostPort(u.Host)
	if err != nil {
		log.Fatal("unable to parse host and port: ", err)
	}

	dbname := u.Path
	if len(dbname) > 0 && dbname[0] == '/' {
		dbname = dbname[1:] // Remove leading slash
	}

	driverConfig := map[string]interface{}{
		"dbname":    dbname,
		"host":      host,
		"port":      port,
		"user":      username,
		"pass":      password,
		"sslmode":   "disable",
		"schema":    "public",
		"blacklist": []string{"schema_migrations"},
	}

	cfg := &boilingcore.Config{
		DriverName:   "psql",
		DriverConfig: driverConfig,
		OutFolder:    "db/model",
		PkgName:      "models",
		Imports:      importers.NewDefaultImports(),
		Wipe:         true,
		NoTests:      true,
		AddEnumTypes: true,
	}
	core, err := boilingcore.New(cfg)
	if err != nil {
		log.Fatal("Failed to initialize sqlboiler core: ", err)
	}
	if err := core.Run(); err != nil {
		log.Fatal("Failed to run sqlboiler code generation: ", err)
	}
}
