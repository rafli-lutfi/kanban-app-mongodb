package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	Id       primitive.ObjectID `json:"id" bson:"_id, omitempty"`
	Fullname string             `json:"fullname" bson:"fullname, omitempty"`
	Email    string             `json:"email" bson:"email, omitempty"`
	Password string             `json:"-" bson:"password, omitempty"`
}

type UserLogin struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserRegister struct {
	Fullname string `json:"fullname" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}
