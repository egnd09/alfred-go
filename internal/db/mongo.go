package db

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// NewMongoClient creates a new MongoDB client
func NewMongoClient(uri string) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	// Ping to verify connection
	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	return client, nil
}

// EnvRepository handles environment collection operations
type EnvRepository struct {
	collection *mongo.Collection
}

// NewEnvRepository creates a new environment repository
func NewEnvRepository(db *mongo.Database) *EnvRepository {
	return &EnvRepository{
		collection: db.Collection("envs"),
	}
}

// JobRepository handles job collection operations
type JobRepository struct {
	collection *mongo.Collection
}

// NewJobRepository creates a new job repository
func NewJobRepository(db *mongo.Database) *JobRepository {
	return &JobRepository{
		collection: db.Collection("jobs"),
	}
}

// UserRepository handles user collection operations
type UserRepository struct {
	collection *mongo.Collection
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *mongo.Database) *UserRepository {
	return &UserRepository{
		collection: db.Collection("users"),
	}
}
