package app

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
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
