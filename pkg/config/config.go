package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort     string
	GrpcPort    string
	AppEnv      string
	DBHost      string
	DBPort      string
	DBUser      string
	DBPassword  string
	DBName      string
	DatabaseDSN string

	JWTSecret     string
	JWTExpiration int // in seconds

	MongoURI        string
	KafkaBrokerURL  string
	KafkaEventTopic string
	// DBName   string
}

func LoadConfig(env string) *Config {
	if err := godotenv.Load("../../.env"); err != nil {
		log.Println("No .env file found, using system env", err)
	}

	mongoURI := getFirstEnv(
		[]string{"MONGO_URI", "MONGODB_URI", "MongoURI"},
		"mongodb://localhost:27017",
	)

	dbName := getFirstEnv(
		[]string{"MONGO_DB_NAME", "MONGODB_DATABASE", "DB_NAME"},
		"notification_pref",
	)

	return &Config{
		AppPort:         getEnv("APP_PORT", "8000"),
		AppEnv:          getEnv("APP_ENV", "development"),
		JWTSecret:       getEnv("JWT_SECRET", "changeme"),
		JWTExpiration:   getEnvAsInt("JWT_EXPIRATION", 3600),
		MongoURI:        mongoURI,
		DBName:          dbName,
		KafkaBrokerURL:  getEnv("KAFKA_BROKER_URL", "localhost:9092"),
		KafkaEventTopic: getEnv("KAFKA_EVENT_TOPIC", "social_events"),
	}
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func getEnvAsInt(key string, fallback int) int {
	if val := os.Getenv(key); val != "" {
		if parsed, err := strconv.Atoi(val); err == nil {
			return parsed
		}
	}
	return fallback
}

func getFirstEnv(keys []string, fallback string) string {
	for _, key := range keys {
		if val := os.Getenv(key); val != "" {
			return val
		}
	}

	return fallback
}
