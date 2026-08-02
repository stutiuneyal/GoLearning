package config

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

const (
	defaultServerPort      = "8080"
	defaultStorageProvider = "local"
	defaultUploadDirectory = "./uploads"
	defaultUploadURLBase   = "/uploads"
	defaultAWSRegion       = "ap-south-1"

	defaultMaxOpenConns = 10
	defaultMaxIdleConns = 5
)

type Config struct {
	ServerPort                 string
	DatabaseURL                string
	JWTSecret                  string
	StorageProvider            string
	UploadDirectory            string
	UploadURLBase              string
	AWSRegion                  string
	S3Bucket                   string
	DatabaseMaxOpenConnections int
	DatabaseMaxIdleConnections int
}

func Load() (Config, error) {

	maxOpenConns, err := intFromEnvironemt("DB_MAX_OPEN_CONNECTIONS", defaultMaxOpenConns)
	if err != nil {
		return Config{}, err
	}

	maxIdleConns, err := intFromEnvironemt("DB_MAX_IDLE_CONNECTIONS", defaultMaxIdleConns)
	if err != nil {
		return Config{}, err
	}

	cfg := Config{
		ServerPort:                 valueOrDefault("APP_PORT", defaultServerPort),
		DatabaseURL:                strings.TrimSpace(os.Getenv("DATABASE_URL")),
		JWTSecret:                  strings.TrimSpace(os.Getenv("JWT_SECRET")),
		StorageProvider:            strings.ToLower(valueOrDefault("STORAGE_PROVIDER", defaultStorageProvider)),
		UploadDirectory:            valueOrDefault("UPLOAD_DIR", defaultUploadDirectory),
		UploadURLBase:              valueOrDefault("UPLOAD_URL_BASE", defaultUploadURLBase),
		AWSRegion:                  valueOrDefault("AWS_REGION", defaultAWSRegion),
		S3Bucket:                   strings.TrimSpace(os.Getenv("S3_BUCKET")),
		DatabaseMaxOpenConnections: maxOpenConns,
		DatabaseMaxIdleConnections: maxIdleConns,
	}

	if err := cfg.validate(); err != nil {
		return Config{}, err
	}

	return cfg, nil

}

func (c Config) validate() error {
	if strings.TrimSpace(c.ServerPort) == "" {
		return fmt.Errorf("APP_PORT is required")
	}

	if strings.TrimSpace(c.DatabaseURL) == "" {
		return fmt.Errorf("DATABASE_URL is required")
	}

	if strings.TrimSpace(c.JWTSecret) == "" {
		return fmt.Errorf("JWT_SECRET is required")
	}

	if len(c.JWTSecret) < 32 {
		return fmt.Errorf(
			"JWT_SECRET must contain at least 32 characters",
		)
	}

	if c.DatabaseMaxOpenConnections <= 0 {
		return fmt.Errorf(
			"DB_MAX_OPEN_CONNECTIONS must be greater than zero",
		)
	}

	if c.DatabaseMaxIdleConnections < 0 {
		return fmt.Errorf(
			"DB_MAX_IDLE_CONNECTIONS cannot be negative",
		)
	}

	if c.DatabaseMaxIdleConnections >
		c.DatabaseMaxOpenConnections {
		return fmt.Errorf(
			"DB_MAX_IDLE_CONNECTIONS cannot exceed " +
				"DB_MAX_OPEN_CONNECTIONS",
		)
	}

	switch c.StorageProvider {
	case "local":
		if strings.TrimSpace(c.UploadDirectory) == "" {
			return fmt.Errorf(
				"UPLOAD_DIR is required when " +
					"STORAGE_PROVIDER=local",
			)
		}

		if strings.TrimSpace(c.UploadURLBase) == "" {
			return fmt.Errorf(
				"UPLOAD_URL_BASE is required when " +
					"STORAGE_PROVIDER=local",
			)
		}

	case "s3":
		if strings.TrimSpace(c.AWSRegion) == "" {
			return fmt.Errorf(
				"AWS_REGION is required when " +
					"STORAGE_PROVIDER=s3",
			)
		}

		if strings.TrimSpace(c.S3Bucket) == "" {
			return fmt.Errorf(
				"S3_BUCKET is required when " +
					"STORAGE_PROVIDER=s3",
			)
		}

	default:
		return fmt.Errorf(
			"unsupported STORAGE_PROVIDER %q; expected local or s3",
			c.StorageProvider,
		)
	}

	return nil
}

func ConnectToDatabase(cfg Config) (*sql.DB, error) {

	db, err := sql.Open("pgx", cfg.DatabaseURL)
	if err != nil {
		return nil, err
	}

	connectionReady := false
	defer func() {
		if !connectionReady {
			db.Close()
		}
	}()

	db.SetMaxOpenConns(cfg.DatabaseMaxOpenConnections)
	db.SetMaxIdleConns(cfg.DatabaseMaxIdleConnections)
	db.SetConnMaxLifetime(30 * time.Minute)
	db.SetConnMaxIdleTime(5 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}

	connectionReady = true

	return db, nil

}

func valueOrDefault(name string, defaultValue string) string {
	value := strings.TrimSpace(os.Getenv(name))

	if value == "" {
		return defaultValue
	}

	return value
}

func intFromEnvironemt(name string, defaultValue int) (int, error) {
	rawValue := strings.TrimSpace(os.Getenv(name))

	if rawValue == "" {
		return defaultValue, nil
	}

	value, err := strconv.Atoi(rawValue)
	if err != nil {
		return 0, fmt.Errorf("%s must be a valid integer: %w", rawValue, err)
	}

	return value, nil
}

func normalizeURLBase(value string) string {
	value = strings.TrimSpace(value)

	if value == "" {
		return ""
	}

	return "/" + strings.Trim(value, "/")
}
