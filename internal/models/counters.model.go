package models

import "go.mongodb.org/mongo-driver/v2/bson"

type CounterModel struct {
	ID    bson.ObjectID `bson:"_id,omitempty" json:"id"`
	Model string        `bson:"model,omitempty" json:"model"`
	Seq   *int          `bson:"seq,omitempty" json:"seq"`
}

type EnumCounterModel string

const (
	EnumCounterModelReminders     = "reminders"
	EnumCounterModelNotifyRecords = "notify_records"
)
