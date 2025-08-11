package config

import "os"

type Config struct {
	Port     string
	GRPCPort string
}

// Load загружает конфигурацию из переменных окружения
func Load() *Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	grpcPort := os.Getenv("GRPC_PORT")
	if grpcPort == "" {
		grpcPort = "6000"
	}

	return &Config{
		Port:     port,
		GRPCPort: grpcPort,
	}
}
