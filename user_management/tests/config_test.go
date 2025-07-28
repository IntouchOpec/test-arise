package tests

import (
	"os"
	"testing"

	"github.com/IntouchOpec/user_management/config"
	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	tests := []struct {
		name     string
		envVars  map[string]string
		expected config.Config
	}{
		{
			name:    "default configuration",
			envVars: map[string]string{},
			expected: config.Config{
				Database: config.DatabaseConfig{
					Host:     "localhost",
					User:     "postgres",
					Password: "password",
					Name:     "users_db",
					Port:     "5432",
					SSLMode:  "disable",
				},
				Server: config.ServerConfig{
					Port: "8080",
				},
				Redis: config.RedisConfig{
					Host:     "redis",
					Port:     "6379",
					Password: "",
					DB:       0,
				},
			},
		},
		{
			name: "custom configuration from env vars",
			envVars: map[string]string{
				"DB_HOST":        "custom-host",
				"DB_USER":        "custom-user",
				"DB_PASSWORD":    "custom-password",
				"DB_NAME":        "custom-db",
				"DB_PORT":        "5433",
				"DB_SSLMODE":     "require",
				"SERVER_PORT":    "3000",
				"REDIS_HOST":     "custom-redis",
				"REDIS_PORT":     "6380",
				"REDIS_PASSWORD": "redis-pass",
			},
			expected: config.Config{
				Database: config.DatabaseConfig{
					Host:     "custom-host",
					User:     "custom-user",
					Password: "custom-password",
					Name:     "custom-db",
					Port:     "5433",
					SSLMode:  "require",
				},
				Server: config.ServerConfig{
					Port: "3000",
				},
				Redis: config.RedisConfig{
					Host:     "custom-redis",
					Port:     "6380",
					Password: "redis-pass",
					DB:       0,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup environment variables
			for key, value := range tt.envVars {
				os.Setenv(key, value)
			}

			// Clean up after test
			defer func() {
				for key := range tt.envVars {
					os.Unsetenv(key)
				}
			}()

			// Execute
			cfg := config.LoadConfig()

			// Assertions
			assert.Equal(t, tt.expected.Database.Host, cfg.Database.Host)
			assert.Equal(t, tt.expected.Database.User, cfg.Database.User)
			assert.Equal(t, tt.expected.Database.Password, cfg.Database.Password)
			assert.Equal(t, tt.expected.Database.Name, cfg.Database.Name)
			assert.Equal(t, tt.expected.Database.Port, cfg.Database.Port)
			assert.Equal(t, tt.expected.Database.SSLMode, cfg.Database.SSLMode)
			assert.Equal(t, tt.expected.Server.Port, cfg.Server.Port)
			assert.Equal(t, tt.expected.Redis.Host, cfg.Redis.Host)
			assert.Equal(t, tt.expected.Redis.Port, cfg.Redis.Port)
			assert.Equal(t, tt.expected.Redis.Password, cfg.Redis.Password)
			assert.Equal(t, tt.expected.Redis.DB, cfg.Redis.DB)
		})
	}
}

func TestDatabaseConfig_GetDSN(t *testing.T) {
	tests := []struct {
		name     string
		config   config.DatabaseConfig
		expected string
	}{
		{
			name: "default configuration",
			config: config.DatabaseConfig{
				Host:     "localhost",
				User:     "postgres",
				Password: "password",
				Name:     "users_db",
				Port:     "5432",
				SSLMode:  "disable",
			},
			expected: "host=localhost user=postgres password=password dbname=users_db port=5432 sslmode=disable",
		},
		{
			name: "custom configuration",
			config: config.DatabaseConfig{
				Host:     "custom-host",
				User:     "custom-user",
				Password: "custom-password",
				Name:     "custom-db",
				Port:     "5433",
				SSLMode:  "require",
			},
			expected: "host=custom-host user=custom-user password=custom-password dbname=custom-db port=5433 sslmode=require",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Execute
			dsn := tt.config.GetDSN()

			// Assertions
			assert.Equal(t, tt.expected, dsn)
		})
	}
}

func TestGetEnv(t *testing.T) {
	// Test getEnv function indirectly through LoadConfig
	t.Run("environment variable exists", func(t *testing.T) {
		os.Setenv("DB_HOST", "test-host")
		defer os.Unsetenv("DB_HOST")

		cfg := config.LoadConfig()
		assert.Equal(t, "test-host", cfg.Database.Host)
	})

	t.Run("environment variable does not exist uses default", func(t *testing.T) {
		os.Unsetenv("DB_HOST")

		cfg := config.LoadConfig()
		assert.Equal(t, "localhost", cfg.Database.Host) // default value
	})

	t.Run("empty environment variable uses default", func(t *testing.T) {
		os.Setenv("DB_HOST", "")
		defer os.Unsetenv("DB_HOST")

		cfg := config.LoadConfig()
		assert.Equal(t, "localhost", cfg.Database.Host) // should use default when empty
	})
}
