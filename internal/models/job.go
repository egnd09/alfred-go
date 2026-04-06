package models

// Job represents a CI/CD job configuration
type Job struct {
	Repo     string   `json:"repo" bson:"repo"`
	Active   bool     `json:"active" bson:"active"`
	Builds   []string `json:"builds" bson:"builds"`
	Branches []string `json:"branches" bson:"branches"`
	Access   string   `json:"access" bson:"access"`
}
