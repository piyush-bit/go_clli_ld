package myutils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// define constanst appname
const APPNAME = "clli_ld"

// define a struct to hold the config
type Config struct {
	BackendURL string `json:"backend_url"`
	// backend tested with default value false
	BackendTested bool `json:"backend_tested"`
}

func LoadConfig() (*Config, error) {

	configDir, err := os.UserConfigDir()
	if err != nil {
		return nil, err
	}
	appConfigDir := filepath.Join(configDir, APPNAME)
	if err := os.MkdirAll(appConfigDir, os.ModePerm); err != nil {
		config, err := InitConfig()
		if err != nil {
			return nil, err
		}
		return config, nil
	}

	// check if config file exists
	configFile := filepath.Join(appConfigDir, "config.json")
	_, err = os.Stat(configFile)
	if err != nil {
		config, err := InitConfig()
		if err != nil {
			return nil, err
		}
		return config, nil
	}

	// read the config file
	file, err := os.ReadFile(configFile)
	if err != nil {
		return nil, err
	}
	var config Config
	err = json.Unmarshal(file, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func InitConfig() (*Config, error) {
	// print user that deafault config is being created use config command to configure the cli
	fmt.Println("No config found, creating default config")
	config := Config{
		BackendURL:    "http://localhost:3000",
		BackendTested: false,
	}
	// save the config to the file
	jsonConfig, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return nil, err
	}
	configDir, err := os.UserConfigDir()
	if err != nil {
		return nil, err
	}
	configFile := filepath.Join(configDir, APPNAME, "config.json")
	err = os.WriteFile(configFile, jsonConfig, 0644)
	if err != nil {
		return nil, err
	}
	// print the config was created at  directory and stae the command to reconfigure the cli
	fmt.Println("Config created at ", configFile)
	fmt.Println("Run config command to reconfigure the cli")
	return &config, nil
}

func UpdateConfig(config *Config) error {
	// save the config to the file
	jsonConfig, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}
	configDir, err := os.UserConfigDir()
	if err != nil {
		return err
	}
	configFile := filepath.Join(configDir, APPNAME, "config.json")
	err = os.WriteFile(configFile, jsonConfig, 0644)
	if err != nil {
		return err
	}
	return nil
}

	func TestBackend(config *Config , backendStatus chan bool)  {
		// test the backend
		fmt.Println("Testing backend")
		//create a http client with 5 sec timeout
		client := &http.Client{
			Timeout: 5 * time.Second,
		}

		// send a get request to the backend
		resp, err := client.Get(config.BackendURL)
		if err != nil {
			fmt.Println("Error: ", err)
			backendStatus <- false
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			fmt.Println("Backend is not running")
			backendStatus <- false
			return
		}
		// Check for the X-Server header
		xServerHeader := resp.Header.Get("X-Server")
		if xServerHeader != "Lab-Digitization" {
			fmt.Println("wrong server configuration")
			backendStatus <- false
			return
		}
		
		backendStatus <- true
		fmt.Println("Backend tested successfully")
	}
