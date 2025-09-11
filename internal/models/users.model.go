package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type UserModel struct {
	ID          bson.ObjectID `bson:"_id,omitempty" json:"id"`
	LineID      string        `bson:"lineId" json:"lineId"`
	Name        string        `bson:"name" json:"name"`
	Deleted     bool          `bson:"deleted" json:"deleted"`
	CreatedAt   time.Time     `bson:"createdAt" json:"createdAt"`
	UpdatedAt   time.Time     `bson:"updatedAt" json:"updatedAt"`
	LastLoginAt time.Time     `bson:"lastLoginAt" json:"lastLoginAt"`
}
