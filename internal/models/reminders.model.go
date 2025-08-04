package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type ReminderModel struct {
	ID         bson.ObjectID       `bson:"_id,omitempty" json:"id"`
	UserID     string              `bson:"userId" json:"userId"`
	Title      string              `bson:"title" json:"title"`
	Content    string              `bson:"content" json:"content"`
	Frequency  EnumRemindFrequency `bson:"frequency" json:"frequency"`
	RemindTime RemindTime          `bson:"remindTime" json:"remindTime"`
	Deleted    bool                `bson:"deleted" json:"deleted"`
	CreatedAt  time.Time           `bson:"createdAt" json:"createdAt"`
	UpdatedAt  time.Time           `bson:"updatedAt" json:"updatedAt"`
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
	Year    *int `bson:"year,omitempty" json:"year"`
	Month   *int `bson:"month,omitempty" json:"month"`
	Date    *int `bson:"date,omitempty" json:"date"`
	Weekday *int `bson:"weekday,omitempty" json:"weekday"`
	Hour    *int `bson:"hour,omitempty" json:"hour"`
	Minute  *int `bson:"minute,omitempty" json:"minute"`
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
