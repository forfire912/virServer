package config

import (
	"os"
	"strconv"
)

// Config holds application configuration
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Storage  StorageConfig
	Auth     AuthConfig
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Host string
	Port int
	Mode string
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// StorageConfig holds storage configuration
type StorageConfig struct {
	ArtifactPath  string
	SnapshotPath  string
	WorkDir       string
	S3Endpoint    string
	S3Bucket      string
	S3AccessKey   string
	S3SecretKey   string
}

// AuthConfig holds auth configuration
type AuthConfig struct {
	JWTSecret  string
	APIKeyAuth bool
	OAuthURL   string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Host: getEnv("SERVER_HOST", "0.0.0.0"),
			Port: getEnvInt("SERVER_PORT", 8080),
			Mode: getEnv("SERVER_MODE", "debug"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnvInt("DB_PORT", 5432),
			User:     getEnv("DB_USER", "virserver"),
			Password: getEnv("DB_PASSWORD", "virserver"),
			DBName:   getEnv("DB_NAME", "virserver"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		Storage: StorageConfig{
			ArtifactPath:  getEnv("ARTIFACT_PATH", "./artifacts"),
			SnapshotPath:  getEnv("SNAPSHOT_PATH", "./snapshots"),
			WorkDir:       getEnv("WORK_DIR", "./work"),
			S3Endpoint:    getEnv("S3_ENDPOINT", ""),
			S3Bucket:      getEnv("S3_BUCKET", "virserver"),
			S3AccessKey:   getEnv("S3_ACCESS_KEY", ""),
			S3SecretKey:   getEnv("S3_SECRET_KEY", ""),
		},
		Auth: AuthConfig{
			JWTSecret:  getEnv("JWT_SECRET", "change-me-in-production"),
			APIKeyAuth: getEnvBool("API_KEY_AUTH", true),
			OAuthURL:   getEnv("OAUTH_URL", ""),
		},
	}
}

// Helper functions

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func getEnvInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	intValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return intValue
}

func getEnvBool(key string, defaultValue bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	boolValue, err := strconv.ParseBool(value)
	if err != nil {
		return defaultValue
	}
	return boolValue
}
