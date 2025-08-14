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
	Hour    *int `json:"hour" binding:"required"`
	Minute  *int `json:"minute" binding:"required"`
}

type ReqGetUserRemindersQuery struct {
	UserId          string                     `form:"userId" binding:"required"`
	Page            *int                       `form:"page"`
	PageSize        *int                       `form:"pageSize"`
	CreateStartTime int64                      `form:"createStartTime"`
	CreateEndTime   int64                      `form:"createEndTime"`
	Title           string                     `form:"title"`
	Content         string                     `form:"content"`
	Frequency       models.EnumRemindFrequency `form:"frequency"`
}

type RespGetUserRemindersBody struct {
	Counts  int64                  `json:"counts"`
	Records []models.ReminderModel `json:"records"`
}

type SearchUserRemindersParams struct {
	UserId          string
	Page            *int
	PageSize        *int
	CreateStartTime int64
	CreateEndTime   int64
	Title           string
	Content         string
	Frequency       models.EnumRemindFrequency
}

type ReqUpdateReminderBody struct {
	ID         string                     `json:"id" binding:"required"`
	UserID     string                     `json:"userId" binding:"required"`
	Frequency  models.EnumRemindFrequency `json:"frequency" binding:"required"`
	Title      string                     `json:"title" binding:"required"`
	Content    string                     `json:"content"`
	RemindTime RemindTimeBody             `json:"remindTime" binding:"required"`
}

type ReqDeleteReminderBody struct {
	ID     string `json:"id" binding:"required"`
	UserID string `json:"userId" binding:"required"`
}
