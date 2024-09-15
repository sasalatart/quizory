package main

import (
	"log"
)

func main() {
	log.Println("Running code generation...")
	defer log.Println("Code generation complete.")

	runSQLBoiler()
	runGoGenerate()
}
