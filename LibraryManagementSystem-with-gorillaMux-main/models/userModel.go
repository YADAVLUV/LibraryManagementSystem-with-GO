package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Id         primitive.ObjectID `bson:"_id, omitempty" json:"user_id"`
	Username   string             `bson:"username" json:"username"`
	Password   string             `bson:"password" json:"password"`
	UserType   string             `bson:"userType" json:"userType"`
	Created_at time.Time          `bson:"created_at" json:"created_at"`
}
