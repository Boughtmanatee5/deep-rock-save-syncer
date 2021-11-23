package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/Boughtmanatee5/deep-rock-save-syncer/app"
	"github.com/vimeo/alog"
)

const (
	EXIT_SUCCESS = 0
	EXIT_ERROR   = 1
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*2)
	defer cancel()

	logger := alog.New(alog.To(os.Stdout))

	homeDir, getDirErr := os.UserHomeDir()
	if getDirErr != nil {
		logger.Printf(ctx, "ERROR can't get home dir %s", getDirErr)
		os.Exit(EXIT_ERROR)
	}

	configPath := fmt.Sprintf("%s/.deepRockSyncConfig", homeDir)

	app := app.NewApp(configPath, logger)

	startErr := app.Start(ctx)
	if startErr != nil {
		logger.Printf(ctx, "Error start error %s", startErr)
		os.Exit(EXIT_ERROR)
	}

	promptError := app.HomePrompt(ctx)
	if promptError != nil {
		fmt.Printf("Error with home prompt %s", promptError)
		os.Exit(EXIT_ERROR)
	}

	os.Exit(EXIT_SUCCESS)
}
