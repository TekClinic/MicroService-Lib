package ms

import (
	"fmt"
	"os"
)

// GetRequiredEnv retrieves the value of the environment variable named by the key.
// If the variable is not present in the environment an error is returned.
func GetRequiredEnv(key string) (string, error) {
	value, set := os.LookupEnv(key)
	if !set {
		return "", fmt.Errorf("%s environment variable is missing", key)
	}
	return value, nil
}

// GetOptionalEnv retrieves the value of the environment variable named by the key.
// If the variable is not present in the environment the def values is returned.
func GetOptionalEnv(key string, def string) string {
	value, set := os.LookupEnv(key)
	if set {
		return value
	}
	return def
}
