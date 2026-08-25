package config

import (
	"os"
	"strconv"
)

func Port() string {
	value := os.Getenv("PORT")
	if value == "" {
		return "8080"
	}
	port, err := strconv.Atoi(value)
	if err == nil && port < 65536 {
		return value
	}
	return "8080"
}
