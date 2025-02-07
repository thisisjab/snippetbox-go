package config

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strings"
	"unicode"
)

type Config struct {
	databasePath   string
	migrationsPath string
}

func (c *Config) DatabasePath() string   { return c.databasePath }
func (c *Config) MigrationsPath() string { return c.migrationsPath }

func LoadConfig() (*Config, error) {
	cfg := &Config{
		databasePath:   "./db.sql",
		migrationsPath: "./cmd/web/db/versions",
	}

	v := reflect.ValueOf(cfg).Elem()
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		envVarName, err := toUpperSnakeCase(field.Name)
		if err != nil {
			return nil, fmt.Errorf("error converting field name '%s': %w", field.Name, err) // Return error immediately
		}

		envVarValue := os.Getenv(envVarName)

		if envVarValue != "" {
			fieldValue := v.Field(i)
			if fieldValue.Kind() == reflect.String {
				fieldValue.SetString(envVarValue)
			} else {
				return nil, fmt.Errorf("field '%s' is not a string, cannot set from env var", field.Name)
			}
		}

	}

	return cfg, nil
}

func toUpperSnakeCase(key string) (string, error) {
	result := ""
	for i, char := range key {
		if !(unicode.IsLetter(char) || unicode.IsNumber(char)) {
			return "", errors.New("key must only contain letters and numbers")
		}

		if unicode.IsUpper(char) {
			if i > 0 {
				result += "_"
			}
		}
		result += string(char)
	}

	return strings.ToUpper(result), nil
}
