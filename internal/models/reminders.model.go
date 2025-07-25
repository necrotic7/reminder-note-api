package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type ReminderModel struct {
	ID         bson.ObjectID       `bson:"_id,omitempty"`
	UserID     string              `bson:"userId"`
	Title      string              `bson:"title"`
	Content    string              `bson:"content"`
	Frequency  EnumRemindFrequency `bson:"frequency"`
	RemindTime RemindTime          `bson:"remindTime"`
	Deleted    bool                `bson:"deleted"`
	CreatedAt  time.Time           `bson:"createdAt"`
	UpdatedAt  time.Time           `bson:"updatedAt"`
}

type EnumRemindFrequency string

const (
	EnumRemindFrequencyOnce     EnumRemindFrequency = "Once"     // 2025/07/17 09:00
	EnumRemindFrequencyDaily    EnumRemindFrequency = "Daily"    // 09:00
	EnumRemindFrequencyWeekly   EnumRemindFrequency = "Weekly"   // Thu 09:00
	EnumRemindFrequencyMonthly  EnumRemindFrequency = "Monthly"  // 07/17 09:00
	EnumRemindFrequencyAnnually EnumRemindFrequency = "Annually" // 07/17 09:00
)

type RemindTime struct {
	Year    *int `bson:"year,omitempty"`
	Month   *int `bson:"month,omitempty"`
	Date    *int `bson:"date,omitempty"`
	Weekday *int `bson:"weekday,omitempty"`
	Hour    *int `bson:"hour,omitempty"`
	Minute  *int `bson:"minute,omitempty"`
}

type InsertReminderParams struct {
	UserID     string              `bson:"userId"`
	Title      string              `bson:"title"`
	Content    string              `bson:"content"`
	Frequency  EnumRemindFrequency `bson:"frequency"`
	RemindTime RemindTime          `bson:"remindTime"`
}

type UpdateReminderParams struct {
	ID         string              `bson:"_id"`
	UserID     string              `bson:"userId"`
	Title      string              `bson:"title"`
	Content    string              `bson:"content"`
	Frequency  EnumRemindFrequency `bson:"frequency"`
	RemindTime RemindTime          `bson:"remindTime"`
}

type DeleteReminderParams struct {
	ID     string `bson:"_id"`
	UserID string `bson:"userId"`
}
