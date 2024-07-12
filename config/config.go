package config

import (
	"ethereum-parser/shared"
	"os"
)

// Config is an alias for the shared Config type
type Config = shared.Config

// LoadConfig loads the configuration from environment variables, providing default values if not set
func LoadConfig() Config {
	return Config{
		ServerPort: getEnv("SERVER_PORT", "8080"),
		RPCURL:     getEnv("RPC_URL", "https://cloudflare-eth.com"),
	}
}

// getEnv retrieves the value of an environment variable or returns a default value if the variable is not set
func getEnv(key, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}
