package repositories

import (
	"context"
	"log"
	"time"

	"github.com/zivwu/reminder-note-api/internal/consts"
	"github.com/zivwu/reminder-note-api/internal/models"
	"github.com/zivwu/reminder-note-api/internal/types"
	"github.com/zivwu/reminder-note-api/internal/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type ReminderRepository struct {
	DB *mongo.Client
}

func NewReminderRepository(db *mongo.Client) *ReminderRepository {
	return &ReminderRepository{
		DB: db,
	}
}

func (r *ReminderRepository) InsertReminder(ctx context.Context, req types.ReqCreateReminderBody) (err error) {
	collection := r.DB.Database("reminder-note").Collection("reminders")
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
		log.Println("insert reminder fail: ", err)
		return err
	}
	return
}

func (r *ReminderRepository) SearchUserReminders(ctx context.Context, params types.ReqGetUserRemindersQuery) ([]models.Reminder, error) {
	collection := r.DB.Database("reminder-note").Collection("reminders")
	filter := bson.M{}
	if !utils.IsEmpty(params.UserId) {
		filter["userId"] = params.UserId
	}

	page := 1
	if params.Page != nil {
		page = *params.Page
	}
	options := options.Find()
	if params.Page == nil {
		*params.Page = 1
	}

	offset := (page - 1) * consts.PageSize
	options.SetLimit(int64(consts.PageSize))
	options.SetSkip(int64(offset))
	options.SetSort(bson.M{"createdAt": -1})

	cursor, err := collection.Find(ctx, filter, options)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []models.Reminder
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}
