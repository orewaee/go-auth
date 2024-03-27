package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Activation struct {
	Id        primitive.ObjectID `json:"id" bson:"_id"`
	User      primitive.ObjectID `json:"user" bson:"user"`
	Secret    string             `json:"secret" bson:"secret"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
}
