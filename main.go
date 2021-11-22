package main

import (
	"fmt"
	"os"

	"github.com/Boughtmanatee5/deep-rock-save-syncer/app"
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

	app := app.NewApp(configPath)

	app.Start()

	promptError := app.HomePrompt()
	if promptError != nil {
		fmt.Printf("Error with home prompt %s", promptError)
		os.Exit(EXIT_ERROR)
	}

	os.Exit(EXIT_SUCCESS)
}
