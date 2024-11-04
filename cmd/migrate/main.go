package main

import (
	"log"

	"github.com/sasalatart/quizory/config"
	"github.com/sasalatart/quizory/db/migrations"
)

func main() {
	log.Println("Running migrations...")
	defer log.Println("Migrations complete.")

	dbCfg := config.NewConfig().DB
	if err := migrations.Up(dbCfg.URL, dbCfg.MigrationsDir()); err != nil {
		log.Fatal("unable to run migrations: ", err)
	}
}
