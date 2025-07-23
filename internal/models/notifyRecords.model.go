package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type NotifyRecordModel struct {
	ID        bson.ObjectID  `bson:"_id,omitempty"`
	UserID    string         `bson:"userId"`
	Content   map[string]any `bson:"content"`
	Status    bool           `bson:"status"`
	Retry     int            `bson:"retry"`
	CreatedAt time.Time      `bson:"createdAt"`
	UpdatedAt time.Time      `bson:"updatedAt"`
}

type InsertNotifyRecord struct {
	UserID  string         `bson:"userId"`
	Content map[string]any `bson:"content"`
	Status  bool           `bson:"status"`
	Retry   int            `bson:"retry"`
}
