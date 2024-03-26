package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Session struct {
	Id           primitive.ObjectID `json:"id" bson:"_id"`
	User         primitive.ObjectID `json:"user" bson:"user"`
	Ip           string             `json:"ip" bson:"ip"`
	RefreshToken string             `json:"-" bson:"refresh_token"`
	CreatedAt    time.Time          `json:"created_at" bson:"created_at"`
}
