package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AccessLog struct {
	ID        primitive.ObjectID `bson:"ID,omitempty"`
	UserID    primitive.ObjectID `bson:"UserId,omitempty"`
	Timestamp time.Time          `bson:"Timestamp"`
	Route     string             `bson:"Route"`
	Method    string             `bson:"Method"`
	Data      interface{}        `bson:"Data,omitempty"`
}
