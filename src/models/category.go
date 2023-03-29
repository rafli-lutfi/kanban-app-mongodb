package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Category struct {
	Id     primitive.ObjectID `json:"id" bson:"_id, omitempty"`
	Type   string             `json:"type" bson:"type"`
	UserId primitive.ObjectID `json:"user_id" bson:"user_id,omitempty"`
	Tasks  []Task             `json:"tasks,omitempty" bson:"tasks,omitempty"`
}

type CategoryRequest struct {
	Type string `json:"type" binding:"required"`
}
