package env

import (
	"os"
	"strconv"
	"strings"
)

// Func GetString returns a string value from the env file
func GetString(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// Func GetInt returns an int value from the env file
func GetInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// Function Load loads the .env file and sets the environment variables
func Load() {
	data, err := os.ReadFile(".env")
	if err != nil {
		return
	}

	for line := range strings.SplitSeq(string(data), "\n") {
		if strings.Contains(line, "=") &&
			!strings.HasPrefix(line, "#") {
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				os.Setenv(strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]))
			}
		}
	}
}
