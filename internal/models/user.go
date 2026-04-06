package models

// User represents an authenticated user
type User struct {
	ID     string `json:"id" bson:"_id"`
	User   string `json:"user" bson:"user"`
	Active bool   `json:"active" bson:"active"`
	Token  string `json:"token,omitempty" bson:"token,omitempty"`
}
