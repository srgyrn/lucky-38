package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

// Config holds required environment variables
type Config struct {
	Driver string
	Source string
}

// Load sets content of configuration file to ENV, reads them and returns Config
func Load(path string) (Config, error) {
	appEnv := os.Getenv("APP_ENV")

	// set APP_ENV to "dev" if not set
	if appEnv == "" {
		os.Setenv("APP_ENV", "development")
	}

	filename := ".env"
	if "development" == appEnv {
		filename = getLocalFilename(path)
	}

	if appEnv == "test" {
		filename = ".env.test"
	}

	if err := godotenv.Overload(filepath.Join(path, filename)); err != nil {
		return Config{}, fmt.Errorf("failed at loading .env file: %v", err)
	}

	return Config{
		Driver: os.Getenv("DB_DRIVER"),
		Source: os.Getenv("DB_SOURCE"),
	}, nil
}

// getLocalFilename checks env files for local environment with following order:
// .env
// .env.development
// .env.local
func getLocalFilename(path string) string {
	for _, f := range []string{".env", ".env.development", ".env.local"} {
		if _, err := os.Stat(filepath.Join(path, f)); !os.IsNotExist(err) {
			return f
		}
	}

	return ""
}
