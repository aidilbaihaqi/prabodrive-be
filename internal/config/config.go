package config

import (
	"log"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Port            string
	GinMode         string
	MaintenanceMode bool

	DB  DBConfig
	JWT JWTConfig
	AWS AWSConfig
	SES SESConfig
	S3  S3Config

	BaseURL string
}

type DBConfig struct {
	Host     string
	Port     string
	Name     string
	User     string
	Password string
	SSLMode  string
}

type JWTConfig struct {
	Secret        string
	AccessExpiry  time.Duration
	RefreshExpiry time.Duration
}

type AWSConfig struct {
	Region          string
	AccessKeyID     string
	SecretAccessKey string
}

type SESConfig struct {
	FromEmail string
}

type S3Config struct {
	Bucket        string
	PresignExpiry time.Duration
}

func Load() *Config {
	secret := getEnv("JWT_SECRET", "")
	if len(secret) < 32 {
		log.Fatal("JWT_SECRET must be at least 32 characters")
	}

	return &Config{
		Port:            getEnv("PORT", "8080"),
		GinMode:         getEnv("GIN_MODE", "debug"),
		MaintenanceMode: getEnvBool("MAINTENANCE_MODE", false),

		DB: DBConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			Name:     getEnv("DB_NAME", "prabodrive"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "secret"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},

		JWT: JWTConfig{
			Secret:        secret,
			AccessExpiry:  getEnvDuration("JWT_ACCESS_EXPIRY", time.Hour),
			RefreshExpiry: getEnvDuration("JWT_REFRESH_EXPIRY", 7*24*time.Hour),
		},

		AWS: AWSConfig{
			Region:          getEnv("AWS_REGION", "ap-southeast-1"),
			AccessKeyID:     getEnv("AWS_ACCESS_KEY_ID", ""),
			SecretAccessKey: getEnv("AWS_SECRET_ACCESS_KEY", ""),
		},

		SES: SESConfig{
			FromEmail: getEnv("SES_FROM_EMAIL", "no-reply@prabodrive.com"),
		},

		S3: S3Config{
			Bucket:        getEnv("S3_BUCKET", "prabodrive-prod"),
			PresignExpiry: getEnvDuration("S3_PRESIGN_EXPIRY", 15*time.Minute),
		},

		BaseURL: getEnv("BASE_URL", ""),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvBool(key string, fallback bool) bool {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	b, err := strconv.ParseBool(v)
	if err != nil {
		return fallback
	}
	return b
}

func getEnvDuration(key string, fallback time.Duration) time.Duration {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	d, err := time.ParseDuration(v)
	if err != nil {
		return fallback
	}
	return d
}
