package config

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	ServicePort         int
	DBMongoURI          string
	RedisURL            string
	GitHubClientID      string
	GitHubClientSecret  string
	AWSAccessKeyID      string
	AWSSecretAccessKey  string
	AWSDefaultRegion    string
	EKSClusterName      string
	JWTSecret           string
}

// Load reads configuration from environment variables
func Load() *Config {
	viper.SetDefault("SERVICE_PORT", 5500)
	viper.SetDefault("JWT_SECRET", "your-secret-key-change-in-production")
	viper.SetDefault("AWS_DEFAULT_REGION", "us-east-1")
	viper.SetDefault("DB_MONGO_URI", "mongodb://localhost:27017")
	viper.SetDefault("REDIS_URL", "redis://localhost:6379")

	viper.AutomaticEnv()

	return &Config{
		ServicePort:        viper.GetInt("SERVICE_PORT"),
		DBMongoURI:         viper.GetString("DB_MONGO_URI"),
		RedisURL:           viper.GetString("REDIS_URL"),
		GitHubClientID:     viper.GetString("GITHUB_CLIENTID"),
		GitHubClientSecret: viper.GetString("GITHUB_CLIENTSECRET"),
		AWSAccessKeyID:     viper.GetString("AWS_ACCESS_KEY_ID"),
		AWSSecretAccessKey: viper.GetString("AWS_SECRET_ACCESS_KEY"),
		AWSDefaultRegion:   viper.GetString("AWS_DEFAULT_REGION"),
		EKSClusterName:     viper.GetString("EKS_CLUSTER_NAME"),
		JWTSecret:          viper.GetString("JWT_SECRET"),
	}
}

// GetJWTSecret returns JWT signing secret
func (c *Config) GetJWTSecret() []byte {
	return []byte(c.JWTSecret)
}

// TokenExpiration is JWT token validity duration
const TokenExpiration = 24 * time.Hour
