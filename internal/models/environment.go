package models

// Environment represents a development environment
type Environment struct {
	Name     string   `json:"name" bson:"name"`
	Services []string `json:"services" bson:"services"`
	Stable   bool     `json:"stable" bson:"stable"`
}
