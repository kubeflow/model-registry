package main

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"
)

func getEnvAsInt(name string, defaultVal int) int {
	if value, exists := os.LookupEnv(name); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultVal
}

func getEnvAsString(name string, defaultVal string) string {
	if value, exists := os.LookupEnv(name); exists {
		return value
	}
	return defaultVal
}

func parseLevel(s string) slog.Level {
	var level slog.Level
	err := level.UnmarshalText([]byte(s))
	if err != nil {
		panic(fmt.Errorf("invalid log level: %s, valid levels are: error, warn, info, debug", s))
	}
	return level
}

func newOriginParser(allowList *[]string, defaultVal string) func(s string) error {
	return func(s string) error {
		value := defaultVal

		if s != "" {
			value = s
		}

		if value == "" {
			return nil
		}

		for _, str := range strings.Split(s, ",") {
			*allowList = append(*allowList, strings.TrimSpace(str))
		}

		return nil
	}
}
