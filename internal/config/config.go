package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	ServerPort string
	Database   DatabaseConfig
	JWTSecret  string
}

type DatabaseConfig struct {
	DSN string
}

func Load() *Config {
	viper.AutomaticEnv()
	viper.SetDefault("PORT", "8080")
	viper.SetDefault("DATABASE_URL", "postgres://user:password@localhost:5432/ratix?sslmode=disable")
	viper.SetDefault("JWT_SECRET", "supersecret")

	// Allow reading from a .env file if it exists, but don't fail if it doesn't
	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Note: No .env file found or error reading it: %v", err)
	}

	config := &Config{
		ServerPort: viper.GetString("PORT"),
		Database: DatabaseConfig{
			DSN: viper.GetString("DATABASE_URL"),
		},
		JWTSecret: viper.GetString("JWT_SECRET"),
	}

	log.Printf("Config loaded: Port=%s", config.ServerPort)
	return config
}
