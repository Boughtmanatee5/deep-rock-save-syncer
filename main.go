package main

import (
	"fmt"
	"os"

	app "github.com/boughtmanatee5/deep-rock-save-syncer"
)

const (
	EXIT_SUCCESS = 0
	EXIT_ERROR   = 1
)

func main() {
	homeDir, getDirErr := os.UserHomeDir()
	if getDirErr != nil {
		fmt.Printf("ERROR can't get home dir %s", getDirErr)
		os.Exit(EXIT_ERROR)
	}

	configPath := fmt.Sprintf("%s/.deepRockSyncConfig", homeDir)

	app := app.New(configPath)

	app.Start()

	os.Exit(EXIT_SUCCESS)
}
