package main

import (
	"log"
	"os"
	"os/exec"
)

func runGoGenerate() {
	codegenCmd := exec.Command("go", "generate", "./...")
	codegenCmd.Stdout = os.Stdout
	codegenCmd.Stderr = os.Stderr
	if err := codegenCmd.Run(); err != nil {
		log.Fatal("unable to run go generate: ", err)
	}
}
