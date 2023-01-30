package config

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
)

type Config struct {
	Server struct {
		Host            string `json:"host"`
		Port            string `json:"port"`
		ReadTimeout     int    `json:"read_timeout"`
		WriteTimeout    int    `json:"write_timeout"`
		ShutDownTimeout int    `json:"shutdown_timeout"`
	} `json:"server"`
	TokenManager struct {
		SessionExpiringTime int    `json:"session_expiring_time"`
		TokenName           string `json:"token_name"`
	} `json:"token_manager"`
}

func LoadConfig(filename string) (Config, error) {
	config := Config{}
	configFile, err := os.Open(filename)
	if err != nil {
		return config, err
	}

	data, err := io.ReadAll(configFile)
	defer configFile.Close()
	if err != nil {
		return config, err
	}

	err = json.Unmarshal(data, &config)
	if err != nil {
		return config, err
	}
	return config, nil
}

func ReadEnv(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("config - ReadEnv - Open: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		sl := strings.Split(line, "=")
		key := sl[0]
		value := sl[1]
		os.Setenv(key, value)
	}
	return nil
}
