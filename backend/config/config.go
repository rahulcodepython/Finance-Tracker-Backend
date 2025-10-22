package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type serverConfig struct {
	Host string
	Port string
}

type database struct {
	DBHost     string
	DBUser     string
	DBPassword string
	DBName     string
	DBPort     string
	DBSSMode   string
}

type jwt struct {
	JWTSecret    string
	JWTExpiresAt string
}

type Config struct {
	ServerConfig      serverConfig
	GoogleOauthConfig *oauth2.Config
	Database          database
	JWT               jwt
}

func parseEnv(key string, defaultValue string) string {
	envValue := os.Getenv(key)
	// This checks if the environment variable is empty.
	if envValue == "" {
		// If the environment variable is empty, a warning is logged.
		log.Printf("%s is missing, default value is set.", key)
		// The default value is returned.
		return defaultValue
	}
	// The value of the environment variable is returned.
	return envValue
}

func LoadConfig() *Config {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

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
			JWTExpiresAt: parseEnv("JWT_EXPIRES_AT", "1h"),
		},
	}
}
