package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	Timeout int      `json:"timeout"`
	Targets []string `json:"targets"`
}

// Bu metot JSON dosyasını okuyup struct'a çevirir
func LoadConfig(filename string) (*Config, error) {
	file, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var cfg Config
	err = json.Unmarshal(file, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
