// Package env provide access to env variable
package env

import (
	"os"
	"strconv"
)

// GetString return string env variable or passed default value
func GetString(key, defaultValue string) string {
	value, exist := os.LookupEnv(key)
	if !exist {
		return defaultValue
	}

	return value
}

// GetBool return bool env variable or passed default value
func GetBool(key string, defaultValue bool) bool {
	value, exist := os.LookupEnv(key)
	if !exist {
		return defaultValue
	}

	boolValue, err := strconv.ParseBool(value)
	if err != nil {
		return defaultValue
	}

	return boolValue
}
