package types

import "github.com/zivwu/reminder-note-api/internal/models"

type ReqCreateReminderBody struct {
	UserID     string                     `json:"userId" binding:"required"`
	Frequency  models.EnumRemindFrequency `json:"frequency" binding:"required"`
	Title      string                     `json:"title" binding:"required"`
	Content    string                     `json:"content"`
	RemindTime RemindTimeBody             `json:"remindTime" binding:"required"`
}

type RemindTimeBody struct {
	Year    *int `json:"year"`
	Month   *int `json:"month"`
	Date    *int `json:"date"`
	Weekday *int `json:"weekday"`
	Hour    *int `json:"hour"`
	Minute  *int `json:"minute"`
}

type ReqGetUserRemindersQuery struct {
	UserId string `form:"userId" binding:"required"`
	Page   *int   `form:"page"`
}

type SearchUserRemindersParams struct {
	UserId     string
	Page       *int
	RemindTime RemindTimeBody
}
