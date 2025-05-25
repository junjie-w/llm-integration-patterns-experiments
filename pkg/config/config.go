package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	OpenAIKey string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	config := &Config{}

	config.OpenAIKey = os.Getenv("OPENAI_API_KEY")
	if config.OpenAIKey == "" {
		return nil, fmt.Errorf("OPENAI_API_KEY environment variable is required")
	}

	return config, nil
}
