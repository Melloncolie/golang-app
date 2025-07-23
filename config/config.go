package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	Port     string `json:"port"`
	Todolist string `json:"todolist"`
	Userlist string `json:"userlist"`
	HomePage string `json:"home_page"`
	LogPage  string `json:"log_page"`
	RegPage  string `json:"reg_page"`
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
