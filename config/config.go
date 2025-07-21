package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	Port         string `json:"port"`
	Todolist     string `json:"todolist"`
	TemplateFile string `json:"templateFile"`
}

func (config *Config) LoadConfig() {
	data, err := os.ReadFile("config/config.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(data, &config)
	if err != nil {
		panic(err)
	}

}
