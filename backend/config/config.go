package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type ServerConfig struct {
	Host string
	Port string
}

type GoogleOAuth2Config struct {
	RedirectUrl  string
	ClientId     string
	ClientSecret string
}

type Database struct {
	DBHost     string
	DBUser     string
	DBPassword string
	DBName     string
	DBPort     string
	DBSSMode   string
}

type JWT struct {
	JWTSecret    string
	JWTExpiresAt string
}

type Config struct {
	ServerConfig       ServerConfig
	GoogleOAuth2Config GoogleOAuth2Config
	GoogleOauthConfig  *oauth2.Config
	Database           Database
	JWT                JWT
}

func ParseEnv(key string, defaultValue string) string {
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
		ServerConfig: ServerConfig{
			Host: ParseEnv("SERVER_HOST", "localhost"),
			Port: ParseEnv("SERVER_PORT", "8080"),
		},
		GoogleOAuth2Config: GoogleOAuth2Config{
			RedirectUrl:  ParseEnv("GOOGLE_REDIRECT_URL", ""),
			ClientId:     ParseEnv("GOOGLE_CLIENT_ID", ""),
			ClientSecret: ParseEnv("GOOGLE_CLIENT_SECRET", ""),
		},
		GoogleOauthConfig: &oauth2.Config{
			RedirectURL:  ParseEnv("GOOGLE_REDIRECT_URL", ""),
			ClientID:     ParseEnv("GOOGLE_CLIENT_ID", ""),
			ClientSecret: ParseEnv("GOOGLE_CLIENT_SECRET", ""),
			Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
			Endpoint:     google.Endpoint,
		},
		Database: Database{
			DBHost:     ParseEnv("DB_HOST", "localhost"),
			DBUser:     ParseEnv("DB_USER", "postgres"),
			DBPassword: ParseEnv("DB_PASSWORD", "postgres"),
			DBName:     ParseEnv("DB_NAME", "finance-tracker"),
			DBPort:     ParseEnv("DB_PORT", "5432"),
			DBSSMode:   ParseEnv("DB_SSL_MODE", "disable"),
		},
		JWT: JWT{
			JWTSecret:    ParseEnv("JWT_SECRET", "secret"),
			JWTExpiresAt: ParseEnv("JWT_EXPIRES_AT", "1h"),
		},
	}
}
