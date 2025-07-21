package services

import (
	"context"
	"log"
	"time"

	"github.com/zivwu/reminder-note-api/internal/models"
	"github.com/zivwu/reminder-note-api/internal/types"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type ReminderService struct {
	DB *mongo.Client
}

func NewReminderService(db *mongo.Client) *ReminderService {
	return &ReminderService{
		DB: db,
	}
}

func (s *ReminderService) CreateReminder(ctx context.Context, req types.ReqCreateReminderBody) (err error) {
	collection := s.DB.Database("reminder-note").Collection("reminders")
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	params := models.Reminder{
		UserID:     req.UserId,
		Title:      req.Title,
		Content:    req.Content,
		Frequency:  req.Frequency,
		RemindTime: models.RemindTime(req.RemindTime),
		Deleted:    false,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	_, err = collection.InsertOne(ctx, params)
	if err != nil {
		log.Println("create reminder fail: ", err)
		return err
	}
	return nil
}
