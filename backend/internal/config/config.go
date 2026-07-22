package config

import (
	"log"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Port           string
	MongoURI       string
	MongoDBName    string
	RedisAddr      string
	RedisPassword  string
	KafkaBrokers   string
	KafkaTopic     string
	JWTSecret      string
	GoogleClientID string
	SeatLockTTL    time.Duration
	FrontendURL    string
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func mustEnv(key string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}

	log.Fatalf("environment variable %s is required", key)
	return ""
}

func Load() *Config {
	seatLockSeconds, err := strconv.Atoi(getEnv("SEAT_LOCK_SECONDS", "300"))
	if err != nil {
		seatLockSeconds = 300
	}

	kafkaBrokers := mustEnv("KAFKA_BROKERS")
	log.Printf(
		"KAFKA_BROKERS=%s",
		kafkaBrokers,
	)

	return &Config{
		// Defaults
		Port:        getEnv("PORT", "8080"),
		MongoDBName: getEnv("MONGO_DB_NAME", "cinema"),
		KafkaTopic:  getEnv("KAFKA_BOOKING_TOPIC", "booking-events"),
		FrontendURL: getEnv("FRONTEND_URL", "http://localhost:5173"),

		// Required
		MongoURI:       mustEnv("MONGO_URI"),
		RedisAddr:      mustEnv("REDIS_ADDR"),
		KafkaBrokers:   mustEnv("KAFKA_BROKERS"),
		JWTSecret:      mustEnv("JWT_SECRET"),
		GoogleClientID: mustEnv("GOOGLE_CLIENT_ID"),

		// Optional
		RedisPassword: getEnv("REDIS_PASSWORD", ""),

		SeatLockTTL: time.Duration(seatLockSeconds) * time.Second,
	}
}
