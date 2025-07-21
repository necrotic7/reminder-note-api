package types

import "github.com/zivwu/reminder-note-api/internal/models"

type ReqCreateReminderBody struct {
	UserId     string                     `json:"userId" binding:"required"`
	Frequency  models.EnumRemindFrequency `json:"frequency" binding:"required"`
	Title      string                     `json:"title" binding:"required"`
	Content    string                     `json:"content"`
	RemindTime RemindTime                 `json:"remindTime" binding:"required"`
}

type RemindTime struct {
	Year    *int `json:"year"`
	Month   *int `json:"month"`
	Date    *int `json:"date"`
	Weekday *int `json:"weekday"`
	Hour    *int `json:"hour"`
	Minute  *int `json:"minute"`
}
