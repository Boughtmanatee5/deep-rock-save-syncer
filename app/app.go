package app

import (
	"context"
	"encoding/json"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"

	"github.com/manifoldco/promptui"
	"github.com/vimeo/alog"
)

type Config struct {
	XboxSavePath  string `json:"xbox_save_path"`
	SteamSavePath string `json:"steam_save_path"`
}

type App struct {
	configPath string
	config     Config
	logger     *alog.Logger
}

func NewApp(configPath string, logger *alog.Logger) *App {
	return &App{
		configPath: configPath,
		logger:     logger,
	}
}

func (a *App) Start(ctx context.Context) error {

	_, statErr := os.Stat(a.configPath)
	if statErr != nil && !os.IsNotExist(statErr) {
		return fmt.Errorf("error can't access config file %w", statErr)
	} else if os.IsNotExist(statErr) {
		a.logger.Print(ctx, "writing empty config")
		emptyConfig := &Config{}
		emptyConfigJson, writeJsonErr := json.Marshal(emptyConfig)
		if writeJsonErr != nil {
			return fmt.Errorf("error marshalling config json", writeJsonErr)
		}

		writeErr := ioutil.WriteFile(a.configPath, emptyConfigJson, fs.ModePerm)
		if writeErr != nil {
			return fmt.Errorf("error creating initial config file %w", writeErr)
		}
	}

	configBytes, readErr := ioutil.ReadFile(a.configPath)
	if readErr != nil {
		return fmt.Errorf("error reading config file %w", readErr)
	}

	var config Config
	jsonParseErr := json.Unmarshal(configBytes, &config)
	if jsonParseErr != nil {
		return fmt.Errorf("error parsing config json %w", jsonParseErr)
	}

	a.logger.Printf(ctx, "config %+v", config)

	a.config = config

	return nil
}

func (a *App) HomePrompt(ctx context.Context) error {

	prompt := promptui.Select{
		Label: "Select an option",
		Items: []string{"Set up", "Sync Xbox to Steam", "Sync Steam to Xbox"},
	}

	index, _, promptErr := prompt.Run()
	if promptErr != nil {
		return fmt.Errorf("error prompting user %w", promptErr)
	}

	switch index {
	case 0:
		config, setupErr := a.setupPrompt(ctx)
		if setupErr != nil {
			return fmt.Errorf("Error setting up config file %w", setupErr)
		}

		writeErr := a.writeConfigFile(ctx, config)
		if writeErr != nil {
			return fmt.Errorf("error saving config changes %w", writeErr)
		}

		break
	case 1:
		a.logger.Printf(ctx, "Backing up xbox save")
		backupErr := backupFile(a.config.XboxSavePath)
		if backupErr != nil {
			return fmt.Errorf("Error backing up save %w", backupErr)
		}
		a.logger.Printf(ctx, "Copying xbox save to steam")
		replaceErr := replaceFile(a.config.XboxSavePath, a.config.SteamSavePath)
		if replaceErr != nil {
			return fmt.Errorf("Error replacing save file %w", replaceErr)
		}
		break
	case 2:
		a.logger.Printf(ctx, "Backing up xbox save")
		backupErr := backupFile(a.config.SteamSavePath)
		if backupErr != nil {
			return fmt.Errorf("Error backing up save %w", backupErr)
		}
		a.logger.Printf(ctx, "Copying xbox save to steam")
		replaceErr := replaceFile(a.config.SteamSavePath, a.config.XboxSavePath)
		if replaceErr != nil {
			return fmt.Errorf("Error replacing save file %w", replaceErr)
		}
	}

	return nil
}

func (a *App) setupPrompt(ctx context.Context) (*Config, error) {

	xboxSavePrompt := promptui.Prompt{
		Label:    "Enter path to Xbox save file",
		Validate: savePathValidator,
	}

	xboxSavePath, promptErr := xboxSavePrompt.Run()
	if promptErr != nil {
		return nil, fmt.Errorf("Error prompting user for xbox save path %w", promptErr)
	}

	a.logger.Printf(ctx, "xboxSavePath %s", xboxSavePath)

	steamSavePrompt := promptui.Prompt{
		Label:    "Enter path to Steam save File",
		Validate: savePathValidator,
	}

	steamSavePath, promptErr := steamSavePrompt.Run()
	if promptErr != nil {
		return nil, fmt.Errorf("Error prompting user for steam save path %w", promptErr)
	}

	a.logger.Printf(ctx, "steamSavePath %s", steamSavePath)

	return &Config{XboxSavePath: xboxSavePath, SteamSavePath: steamSavePath}, nil
}

func (a *App) writeConfigFile(ctx context.Context, config *Config) error {
	a.logger.Printf(ctx, "Config being written %+v", config)
	jsonData, jsonErr := json.Marshal(config)
	if jsonErr != nil {
		return fmt.Errorf("error converting config to json %w", jsonErr)
	}
	a.logger.Printf(ctx, "data being writter %s", jsonData)
	a.logger.Printf(ctx, "writting to configPath %s", a.configPath)
	writeErr := ioutil.WriteFile(a.configPath, jsonData, fs.ModePerm)
	if writeErr != nil {
		fmt.Errorf("error writing config to config file %w", writeErr)
	}

	return nil
}

func savePathValidator(s string) error {
	if len(s) == 0 {
		return fmt.Errorf("You must enter a path to your save file")
	}

	_, statErr := os.Stat(s)

	switch true {
	case os.IsNotExist(statErr):
		return fmt.Errorf("The file path you entered does not exist please check again")
	case os.IsPermission(statErr):
		return fmt.Errorf("Deep Rock Galactic Save Sync does not have permission to use this file")
	case os.IsTimeout(statErr):
		return fmt.Errorf("Timed out while reading file try again")
	case statErr != nil:
		return fmt.Errorf("Encountered some kind of error try again")
	default:
		return nil
	}
}

func backupFile(filePath string) error {
	fileBytes, readError := ioutil.ReadFile(filePath)
	if readError != nil {
		return fmt.Errorf("error reading file %w", readError)
	}

	backupPath := filePath + "_backup"
	writeErr := ioutil.WriteFile(backupPath, fileBytes, fs.ModePerm)
	if writeErr != nil {
		fmt.Errorf("error writing backup file %w", writeErr)
	}

	return nil
}

func replaceFile(inputPath string, outputPath string) error {
	fileBytes, readError := ioutil.ReadFile(inputPath)
	if readError != nil {
		return fmt.Errorf("error reading file %w", readError)
	}

	writeErr := ioutil.WriteFile(outputPath, fileBytes, fs.ModePerm)
	if writeErr != nil {
		fmt.Errorf("error writing backup file %w", writeErr)
	}

	return nil
}
