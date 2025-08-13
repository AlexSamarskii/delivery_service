package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	EmailConfig EmailConfig
	RedisConfig RedisConfig

	UserServiceURL       string
	FoodServiceURL       string
	RestaurantServiceURL string
	CartServiceURL       string
	PaymentServiceURL    string

	GrpcCatServiceURL  string
	GrpcFoodServiceURL string

	Port     string
	GRPCPort string
}

type EmailConfig struct {
	SMTPHost     string
	SMTPPort     int
	SMTPUsername string
	SMTPPassword string
}

type RedisConfig struct {
	Host     string
	Password string
}

type GoogleConfig struct {
	ClientId     string
	ClientSecret string
	RedirectUrl  string
}

type MinIoConfig struct {
	AccessKey  string
	BucketName string
	Domain     string
	Region     string
	SecretKey  string
	UseSSL     bool
}

type ElasticSearchConfig struct {
	Addresses []string
	Username  string
	Password  string
	CloudID   string
	APIKey    string
	IndexName string
}

var config *Config

func Load() *Config {
	if config != nil {
		return config
	}

	if err := godotenv.Load(); err != nil {
		log.Println(".env file not found or failed to load — using environment variables")
	}

	config = &Config{
		EmailConfig: loadEmailConfig(),
		RedisConfig: loadRedisConfig(),

		UserServiceURL:       getEnv("USER_SERVICE_URL", ""),
		FoodServiceURL:       getEnv("FOOD_SERVICE_URL", ""),
		RestaurantServiceURL: getEnv("RESTAURANT_SERVICE_URL", ""),
		CartServiceURL:       getEnv("CART_SERVICE_URL", ""),
		PaymentServiceURL:    getEnv("PAYMENT_SERVICE_URL", ""),

		GrpcCatServiceURL:  getEnv("GRPC_CAT_SERVICE_URL", ""),
		GrpcFoodServiceURL: getEnv("GRPC_FOOD_SERVICE_URL", ""),

		Port:     getEnv("PORT", "3000"),
		GRPCPort: getEnv("GRPC_PORT", "6000"),
	}

	validateConfig(config)

	return config
}

func Get() *Config {
	return Load()
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

// парсинг SMTP порта с проверкой
func getEnvInt(key string, fallback int) int {
	if value, exists := os.LookupEnv(key); exists {
		if parsed, err := strconv.Atoi(value); err == nil {
			return parsed
		}
		log.Printf("Invalid value for %s: %s, using fallback: %d", key, value, fallback)
	}
	return fallback
}

func loadEmailConfig() EmailConfig {
	port := getEnvInt("SMTP_PORT", 587) // 587 — стандартный порт
	return EmailConfig{
		SMTPHost:     getEnv("SMTP_HOST", ""),
		SMTPPort:     port,
		SMTPUsername: getEnv("SMTP_USERNAME", ""),
		SMTPPassword: getEnv("SMTP_PASSWORD", ""),
	}
}

func loadRedisConfig() RedisConfig {
	return RedisConfig{
		Host:     getEnv("REDIS_ADDR", "localhost:6379"),
		Password: getEnv("REDIS_PASSWORD", ""),
	}
}

func loadGoogleConfig() GoogleConfig {
	return GoogleConfig{
		ClientId:     getEnv("GG_CLIENT_ID", ""),
		ClientSecret: getEnv("GG_CLIENT_SECRET", ""),
		RedirectUrl:  getEnv("GG_REDIRECT_URL", "http://localhost:3000/auth/google/callback"),
	}
}

func loadMinioConfig() MinIoConfig {
	useSSL := false
	if sslStr := getEnv("MINIO_USE_SSL", "false"); sslStr == "true" || sslStr == "1" {
		useSSL = true
	}

	return MinIoConfig{
		AccessKey:  getEnv("MINIO_ACCESS_KEY", ""),
		BucketName: getEnv("MINIO_BUCKET_NAME", "my-bucket"),
		Domain:     getEnv("MINIO_DOMAIN", ""),
		Region:     getEnv("MINIO_REGION", "us-east-1"),
		SecretKey:  getEnv("MINIO_SECRET_KEY", ""),
		UseSSL:     useSSL,
	}
}

func loadElasticSearchConfig() ElasticSearchConfig {
	address := getEnv("ES_ADDRESS", "http://localhost:9200")
	if address == "" {
		address = "http://localhost:9200"
	}

	return ElasticSearchConfig{
		Addresses: []string{address},
		Username:  getEnv("ES_USERNAME", ""),
		Password:  getEnv("ES_PASSWORD", ""),
		CloudID:   getEnv("ES_CLOUD_ID", ""),
		APIKey:    getEnv("ES_API_KEY", ""),
		IndexName: getEnv("ES_INDEX_NAME", "products"),
	}
}

func validateConfig(cfg *Config) {
	required := map[string]string{
		"UserServiceURL":     cfg.UserServiceURL,
		"GrpcCatServiceURL":  cfg.GrpcCatServiceURL,
		"GrpcFoodServiceURL": cfg.GrpcFoodServiceURL,
		"SMTP_HOST":          cfg.EmailConfig.SMTPHost,
		"SMTP_USERNAME":      cfg.EmailConfig.SMTPUsername,
		"SMTP_PASSWORD":      cfg.EmailConfig.SMTPPassword,
		"REDIS_ADDR":         cfg.RedisConfig.Host,
	}

	for envName, value := range required {
		if value == "" {
			log.Fatalf("Required environment variable missing: %s", envName)
		}
	}
}
