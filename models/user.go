package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type User struct {
	Id        primitive.ObjectID `json:"id" bson:"_id"`
	Email     string             `json:"email" bson:"email"`
	Name      string             `json:"name" bson:"name"`
	Password  string             `json:"-" bson:"password"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
}
