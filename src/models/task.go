package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Task struct {
	Id          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Title       string             `json:"title" bson:"title,omitempty"`
	Description string             `json:"description" bson:"description,omitempty"`
	CategoryId  primitive.ObjectID `json:"category_id" bson:"category_id,omitempty"`
	UserId      primitive.ObjectID `json:"user_id" bson:"user_id,omitempty"`
}

type RequestTask struct {
	CategoryID  string `json:"category_id" binding:"required"`
	Title       string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`
}

type TaskCategoryRequest struct {
	ID         primitive.ObjectID `json:"id"`
	CategoryID primitive.ObjectID `json:"category_id" binding:"required"`
}
