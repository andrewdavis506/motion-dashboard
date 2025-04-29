package config

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	APIKey          string
	Port            int
	TemplatesDir    string
	RefreshInterval time.Duration
}

// LoadConfig loads environment variables and command-line flags into a Config struct
func LoadConfig() (*Config, error) {
	// Load .env if available
	if err := godotenv.Load(); err != nil {
		slog.Warn("Note: .env file not found. Using environment variables and flags.")
	}

	// Define flags
	apiKeyFlag := flag.String("api-key", "", "Motion API key")
	portFlag := flag.Int("port", 8080, "Web server port")
	templatesDirFlag := flag.String("templates", "templates", "Templates directory")
	refreshFlag := flag.Duration("refresh", 60*time.Second, "Task refresh interval")
	flag.Parse()

	// Create config with initial values from flags
	cfg := &Config{
		APIKey:          *apiKeyFlag,
		Port:            *portFlag,
		TemplatesDir:    *templatesDirFlag,
		RefreshInterval: *refreshFlag,
	}

	// Override with environment variables if needed
	if cfg.APIKey == "" {
		if apiKeyEnv := os.Getenv("MOTION_API_KEY"); apiKeyEnv != "" {
			cfg.APIKey = apiKeyEnv
		} else {
			return nil, fmt.Errorf("API key is required. Set it with -api-key flag or MOTION_API_KEY environment variable")
		}
	}

	if portEnv := os.Getenv("PORT"); portEnv != "" && cfg.Port == 8080 {
		if portNum, err := strconv.Atoi(portEnv); err == nil {
			cfg.Port = portNum
		}
	}

	if refreshEnv := os.Getenv("REFRESH_INTERVAL"); refreshEnv != "" && cfg.RefreshInterval == 60*time.Second {
		if duration, err := time.ParseDuration(refreshEnv); err == nil {
			cfg.RefreshInterval = duration
		}
	}

	return cfg, nil
}
