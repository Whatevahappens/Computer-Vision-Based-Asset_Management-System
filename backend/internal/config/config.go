package config

import (
	"os"
	"strconv"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	RedisAddr     string
	RedisPassword string
	RedisDB       int

	JWTSecret      string
	JWTExpiryHours int

	ServerPort   string
	GinMode      string
	CVServiceURL string

	AdminEmail    string
	AdminPassword string
}

func Load() *Config {
	return &Config{
		DBHost:         getEnv("DB_HOST", "localhost"),
		DBPort:         getEnv("DB_PORT", "5432"),
		DBUser:         getEnv("DB_USER", "asset_admin"),
		DBPassword:     getEnv("DB_PASSWORD", "changeme_secret"),
		DBName:         getEnv("DB_NAME", "asset_management"),
		DBSSLMode:      getEnv("DB_SSLMODE", "disable"),
		RedisAddr:      getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPassword:  getEnv("REDIS_PASSWORD", ""),
		RedisDB:        getEnvInt("REDIS_DB", 0),
		JWTSecret:      getEnv("JWT_SECRET", "your-256-bit-secret-change-in-production"),
		JWTExpiryHours: getEnvInt("JWT_EXPIRY_HOURS", 24),
		ServerPort:     getEnv("SERVER_PORT", "8080"),
		GinMode:        getEnv("GIN_MODE", "debug"),
		CVServiceURL:   getEnv("CV_SERVICE_URL", "http://localhost:8000"),
		AdminEmail:     getEnv("ADMIN_EMAIL", "admin@must.edu.mn"),
		AdminPassword:  getEnv("ADMIN_PASSWORD", "Admin@123"),
	}
}

func (c *Config) DSN() string {
	return "host=" + c.DBHost + " port=" + c.DBPort + " user=" + c.DBUser +
		" password=" + c.DBPassword + " dbname=" + c.DBName + " sslmode=" + c.DBSSLMode
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return fallback
}
