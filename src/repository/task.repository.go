package repository

import (
	"context"

	"github.com/rafli-lutfi/kanban-app-mongodb/src/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TaskRepository interface {
	GetTaskByID(ctx context.Context, taskID primitive.ObjectID) (models.Task, error)
	StoreTask(ctx context.Context, task models.Task) (interface{}, error)
	UpdateTask(ctx context.Context, task models.Task) error
	DeleteTask(ctx context.Context, taskID primitive.ObjectID) error
	DeleteAllTaskByCategory(ctx context.Context, categoryID primitive.ObjectID) error
	AppendTaskToCategoryTasks(ctx context.Context, task models.Task) error
}

type taskRepository struct {
	db *mongo.Database
}

func NewTaskRepository(db *mongo.Database) *taskRepository {
	return &taskRepository{db}
}

func (r *taskRepository) GetTaskByID(ctx context.Context, taskID primitive.ObjectID) (models.Task, error) {
	collection := r.db.Collection("tasks")

	var task models.Task

	err := collection.FindOne(ctx, bson.M{"_id": taskID}).Decode(&task)
	if err != nil {
		return models.Task{}, nil
	}

	return task, nil
}

func (r *taskRepository) StoreTask(ctx context.Context, task models.Task) (interface{}, error) {
	collection := r.db.Collection("tasks")

	result, err := collection.InsertOne(ctx, &task)
	if err != nil {
		return nil, err
	}

	err = r.AppendTaskToCategoryTasks(ctx, task)
	if err != nil {
		return nil, err
	}

	return result.InsertedID, nil
}

func (r *taskRepository) UpdateTask(ctx context.Context, task models.Task) error {
	collection := r.db.Collection("tasks")

	filter := bson.D{{Key: "_id", Value: task.Id}}
	update := bson.M{"$set": task}

	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (r *taskRepository) DeleteTask(ctx context.Context, taskID primitive.ObjectID) error {
	collection := r.db.Collection("tasks")

	filter := bson.D{{Key: "_id", Value: taskID}}

	result, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return models.ErrRecordNotFound
	}

	return nil
}

func (r *taskRepository) DeleteAllTaskByCategory(ctx context.Context, categoryID primitive.ObjectID) error {
	collection := r.db.Collection("tasks")

	filter := bson.D{{Key: "category_id", Value: categoryID}}

	_, err := collection.DeleteMany(ctx, filter)
	if err != nil {
		return err
	}

	return nil
}

func (r *taskRepository) AppendTaskToCategoryTasks(ctx context.Context, task models.Task) error {
	collection := r.db.Collection("categories")

	filter := bson.D{{Key: "_id", Value: task.CategoryId}}

	data := models.Task{
		Id:          task.Id,
		Title:       task.Title,
		Description: task.Description,
	}

	_, err := collection.UpdateMany(ctx, filter, bson.M{"$push": bson.M{"tasks": data}})
	if err != nil {
		return err
	}

	return nil
}
