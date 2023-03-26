package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Task struct {
	Id          primitive.ObjectID `json:"id" bson:"_id, omitempty"`
	Title       string             `json:"title" bson:"title"`
	Description string             `json:"description" bson:"description"`
}

type RequestTask struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`
}

type TaskCategoryRequest struct {
	ID         int `json:"id"`
	CategoryID int `json:"category_id" binding:"required"`
}
