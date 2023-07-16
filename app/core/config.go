package core

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"time"
)

type Config struct {
	EmailHost     string        `json:"email_host"`
	EmailUsername string        `json:"email_username"`
	EmailPassword string        `json:"email_password"`
	EmailPort     int           `json:"email_port"`
	EmailTo       string        `json:"email_to"`
	Port          int           `json:"port"`
	Count         int           `json:"count"`
	Timeout       time.Duration `json:"timeout"`
}

func NewConfig() *Config {
	path := filepath.Join(".", "config.json")
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("cannot open config file: %v", err)
	}
	defer file.Close()

	js := json.NewDecoder(file)

	config := &Config{}
	if err := js.Decode(config); err != nil {
		log.Fatalf("cannot decode file config: %v", err)
	}

	return config
}
