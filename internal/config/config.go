package config

import (
	"errors"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Env          string
	Port         int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration

	DatabaseURL string

	RedisAddr     string
	RedisPassword string
	RedisDB       int
}

func Load() (Config, error) {
	env := getEnv("ENV", "dev")
	portStr := getEnv("PORT", "8080")
	port, err := strconv.Atoi(portStr)
	if err != nil || port <= 0 {
		return Config{}, errors.New("invalid PORT")
	}

	dbURL := getEnv("DATABASE_URL", "")
	if dbURL == "" {
		return Config{}, errors.New("DATABASE_URL is required")
	}

	redisDBstr := getEnv("REDIS_DB", "0")
	redisDB, err := strconv.Atoi(redisDBstr)
	if err != nil || redisDB < 0 {
		return Config{}, errors.New("Invalid REDIS_DB")
	}

	return Config{
		Env:           env,
		Port:          port,
		ReadTimeout:   10 * time.Second,
		WriteTimeout:  10 * time.Second,
		IdleTimeout:   60 * time.Second,
		DatabaseURL:   dbURL,
		RedisAddr:     getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		RedisDB:       redisDB,
	}, nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
