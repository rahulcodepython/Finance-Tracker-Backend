// Package config provides functionality for loading and managing application configuration.
package config

import (
	"log"
	"os"

	"github.com/joho/godotenv" // Used for loading environment variables from a .env file
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// serverConfig holds the configuration for the HTTP server.
type serverConfig struct {
	Host string // Server host
	Port string // Server port
}

// database holds the configuration for the database connection.
type database struct {
	DBHost     string // Database host
	DBUser     string // Database user
	DBPassword string // Database password
	DBName     string // Database name
	DBPort     string // Database port
	DBSSMode   string // Database SSL mode
}

// jwt holds the configuration for JWT authentication.
type jwt struct {
	JWTSecret    string // Secret key for signing JWT tokens
	JWTExpiresIn string // Expiration duration for JWT tokens
}

// Config holds the entire application configuration.
type Config struct {
	ServerConfig      serverConfig      // Server configuration
	GoogleOauthConfig *oauth2.Config    // Google OAuth2 configuration
	Database          database          // Database configuration
	JWT               jwt               // JWT configuration
}

// parseEnv retrieves the value of an environment variable or returns a default value if the variable is not set.
func parseEnv(key string, defaultValue string) string {
	envValue := os.Getenv(key)
	if envValue == "" {
		log.Printf("%s is missing, default value is set.", key)
		return defaultValue
	}
	return envValue
}

// LoadConfig loads the application configuration from a .env file and the environment.
func LoadConfig() *Config {
	// Load environment variables from the .env file.
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Return a new Config struct populated with values from the environment.
	return &Config{
		ServerConfig: serverConfig{
			Host: parseEnv("HOST", "localhost"),
			Port: parseEnv("PORT", "8000"),
		},
		GoogleOauthConfig: &oauth2.Config{
			RedirectURL:  parseEnv("GOOGLE_OAUTH_REDIRECT_URL", ""),
			ClientID:     parseEnv("GOOGLE_CLIENT_ID", ""),
			ClientSecret: parseEnv("GOOGLE_CLIENT_SECRET", ""),
			Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
			Endpoint:     google.Endpoint,
		},
		Database: database{
			DBHost:     parseEnv("DB_HOST", "localhost"),
			DBUser:     parseEnv("DB_USER", "postgres"),
			DBPassword: parseEnv("DB_PASSWORD", "admin"),
			DBName:     parseEnv("DB_NAME", "finance_tracker"),
			DBPort:     parseEnv("DB_PORT", "5432"),
			DBSSMode:   parseEnv("DB_SSL_MODE", "disable"),
		},
		JWT: jwt{
			JWTSecret:    parseEnv("JWT_SECRET", "secret"),
			JWTExpiresIn: parseEnv("JWT_EXPIRES_IN", "1h"),
		},
	}
}