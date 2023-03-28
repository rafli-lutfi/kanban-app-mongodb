package repository

import (
	"context"

	"github.com/rafli-lutfi/kanban-app-mongodb/src/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CategoryRepository interface {
	GetCategoryByID(ctx context.Context, id primitive.ObjectID) (models.Category, error)
	GetCategories(ctx context.Context, userID primitive.ObjectID) ([]models.Category, error)
	StoreCategory(ctx context.Context, category models.Category) (interface{}, error)
	StoreManyCategory(ctx context.Context, categories []interface{}) (interface{}, error)
	DeleteCategory(ctx context.Context, categoryID primitive.ObjectID) error
}

type categoryRepository struct {
	db *mongo.Database
}

func NewCategoryRepository(db *mongo.Database) *categoryRepository {
	return &categoryRepository{db}
}

func (r *categoryRepository) GetCategoryByID(ctx context.Context, id primitive.ObjectID) (models.Category, error) {
	var collection = r.db.Collection("categories")

	var category models.Category

	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&category)
	if err != nil {
		return models.Category{}, err
	}

	return category, nil
}

func (r *categoryRepository) GetCategories(ctx context.Context, userID primitive.ObjectID) ([]models.Category, error) {
	var collection = r.db.Collection("categories")

	filter := bson.D{{Key: "user_id", Value: userID}}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return []models.Category{}, err
	}

	defer cursor.Close(ctx)

	var categories []models.Category

	for cursor.Next(ctx) {
		var category models.Category

		err := cursor.Decode(&category)
		if err != nil {
			return []models.Category{}, err
		}

		categories = append(categories, category)
	}

	if err := cursor.Err(); err != nil {
		return []models.Category{}, err
	}

	return categories, nil
}

func (r *categoryRepository) StoreCategory(ctx context.Context, category models.Category) (interface{}, error) {
	var collection = r.db.Collection("categories")

	result, err := collection.InsertOne(ctx, &category)
	if err != nil {
		return nil, err
	}

	return result.InsertedID, nil
}

func (r *categoryRepository) StoreManyCategory(ctx context.Context, categories []interface{}) (interface{}, error) {
	var collection = r.db.Collection("categories")

	result, err := collection.InsertMany(ctx, categories, options.InsertMany().SetOrdered(true))
	if err != nil {
		return nil, err
	}

	return result.InsertedIDs, nil
}

func (r *categoryRepository) DeleteCategory(ctx context.Context, categoryID primitive.ObjectID) error {
	var collection = r.db.Collection("categories")

	_, err := collection.DeleteOne(ctx, bson.M{"_id": categoryID})
	if err != nil {
		return err
	}
	return nil
}
