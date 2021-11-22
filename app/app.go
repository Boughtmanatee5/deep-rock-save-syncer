package app

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"

	"github.com/manifoldco/promptui"
)

type Config struct {
	xboxSavePath  string
	steamSavePath string
}

type App struct {
	configPath string
	config     Config
}

func NewApp(configPath string) *App {
	return &App{
		configPath: configPath,
	}
}

func (a *App) Start() error {

	_, statErr := os.Stat(a.configPath)
	if statErr != nil && !os.IsNotExist(statErr) {
		return fmt.Errorf("error can't access config file %w", statErr)
	} else if os.IsNotExist(statErr) {
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
	jsonParseErr := json.Unmarshal(configBytes, config)
	if jsonParseErr != nil {
		return fmt.Errorf("error parsing config json %w", jsonParseErr)
	}

	a.config = config

	return nil
}

func (a *App) HomePrompt() error {

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
		setupErr := a.setupPrompt()
		if setupErr != nil {
			return fmt.Errorf("Error setting up config file %w", setupErr)
		}

		writeErr := a.writeConfigFile()
		if writeErr != nil {
			return fmt.Errorf("error saving config changes %w", writeErr)
		}

		break
	case 1:
		fmt.Print("Copy xbox save to steam save")
		break
	case 2:
		fmt.Printf("Copy steam save to xbox save")
	}

	return nil
}

func (a *App) setupPrompt() error {

	xboxSavePrompt := promptui.Prompt{
		Label:    "Enter path to Xbox save file",
		Validate: savePathValidator,
	}

	xboxSavePath, promptErr := xboxSavePrompt.Run()
	if promptErr != nil {
		return fmt.Errorf("Error prompting user for xbox save path %w", promptErr)
	}

	steamSavePrompt := promptui.Prompt{
		Label:    "Enter path to Steam save File",
		Validate: savePathValidator,
	}

	steamSavePath, promptErr := steamSavePrompt.Run()
	if promptErr != nil {
		return fmt.Errorf("Error prompting user for steam save path %w", promptErr)
	}

	a.config = Config{
		xboxSavePath:  xboxSavePath,
		steamSavePath: steamSavePath,
	}

	return nil
}

func (a *App) writeConfigFile() error {
	jsonData, jsonErr := json.Marshal(a.config)
	if jsonErr != nil {
		return fmt.Errorf("error converting config to json %w", jsonErr)
	}

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
