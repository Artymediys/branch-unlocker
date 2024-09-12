package env

import (
	"encoding/json"
	"errors"
	"os"
)

type GitLabConfig struct {
	URL      string   `json:"url"`
	Token    string   `json:"token"`
	GroupID  int      `json:"group_id"`
	Branches []string `json:"branches"`
}

func NewConfig(configPath string) (*GitLabConfig, error) {
	if configPath == "" {
		return nil, errors.New("configuration file path is required")
	}

	var cfg GitLabConfig
	err := cfg.loadConfig(configPath)
	if err != nil {
		return nil, errors.New("error loading config: " + err.Error())
	}

	return &cfg, err
}

func (cfg *GitLabConfig) loadConfig(configPath string) error {
	file, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}

	err = json.Unmarshal(file, cfg)
	if err != nil {
		return err
	}

	return nil
}
