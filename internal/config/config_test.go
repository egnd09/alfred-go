package config

import (
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	// Set test environment variables
	os.Setenv("SERVICE_PORT", "8080")
	os.Setenv("DB_MONGO_URI", "mongodb://test:27017")
	os.Setenv("REDIS_URL", "redis://test:6379")
	os.Setenv("LOG_LEVEL", "debug")
	
	defer func() {
		os.Unsetenv("SERVICE_PORT")
		os.Unsetenv("DB_MONGO_URI")
		os.Unsetenv("REDIS_URL")
		os.Unsetenv("LOG_LEVEL")
	}()
	
	cfg := Load()
	
	if cfg.ServicePort != 8080 {
		t.Errorf("ServicePort = %v, want 8080", cfg.ServicePort)
	}
	if cfg.DBMongoURI != "mongodb://test:27017" {
		t.Errorf("DBMongoURI = %v, want mongodb://test:27017", cfg.DBMongoURI)
	}
	if cfg.RedisURL != "redis://test:6379" {
		t.Errorf("RedisURL = %v, want redis://test:6379", cfg.RedisURL)
	}
	if cfg.LogLevel != "debug" {
		t.Errorf("LogLevel = %v, want debug", cfg.LogLevel)
	}
}

func TestLoadDefaults(t *testing.T) {
	// Clear environment
	os.Clearenv()
	
	cfg := Load()
	
	if cfg.ServicePort != 5500 {
		t.Errorf("default ServicePort = %v, want 5500", cfg.ServicePort)
	}
	if cfg.DBMongoURI != "mongodb://localhost:27017/alfred" {
		t.Errorf("default DBMongoURI = %v", cfg.DBMongoURI)
	}
}
