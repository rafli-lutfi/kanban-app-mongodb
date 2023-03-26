package repository

import (
	"context"

	"github.com/rafli-lutfi/kanban-app-mongodb/src/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository interface {
	Register(ctx context.Context, user models.UserRegister) (interface{}, error)
	FindUserByEmail(ctx context.Context, email string) (models.User, error)
}

type userRepository struct {
	db *mongo.Database
}

func NewUserRepository(db *mongo.Database) *userRepository {
	return &userRepository{db}
}

func (r *userRepository) Register(ctx context.Context, user models.UserRegister) (interface{}, error) {
	var collection = r.db.Collection("users")

	result, err := collection.InsertOne(ctx, &user)
	if err != nil {
		return nil, err
	}

	return result.InsertedID, nil
}

func (r *userRepository) FindUserByEmail(ctx context.Context, email string) (models.User, error) {
	var collection = r.db.Collection("users")

	var user models.User

	err := collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}
