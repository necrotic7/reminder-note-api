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

type RemindersRepository struct {
	db         *mongo.Database
	collection *mongo.Collection
}

func NewRemindersRepository(db *mongo.Database) *RemindersRepository {
	return &RemindersRepository{
		db:         db,
		collection: db.Collection("reminders"),
	}
}

func (r *RemindersRepository) InsertReminder(ctx context.Context, params *models.InsertReminderParams) (err error) {
	doc := models.ReminderModel{
		UserID:     params.UserID,
		Title:      params.Title,
		Content:    params.Content,
		Frequency:  params.Frequency,
		RemindTime: models.RemindTime(params.RemindTime),
		Deleted:    false,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	_, err = r.collection.InsertOne(ctx, doc)
	if err != nil {
		log.Println("insert reminder fail: ", err)
		return err
	}
	return
}

func (r *RemindersRepository) UpdateReminder(ctx context.Context, params *models.UpdateReminderParams) (err error) {
	// 轉換 ID 為 ObjectID
	objID, err := bson.ObjectIDFromHex(params.ID)
	if err != nil {
		log.Println("update reminders fail: ", err)
		return err
	}

	filter := bson.M{
		"_id":    objID,
		"userId": params.UserID,
	}

	updateParams := bson.M{
		"updatedAt": time.Now(),
	}

	if !utils.IsEmpty(params.Frequency) {
		updateParams["frequency"] = params.Frequency
	}

	if !utils.IsEmpty(params.Title) {
		updateParams["title"] = params.Title
	}

	if !utils.IsEmpty(params.Content) {
		updateParams["content"] = params.Content
	}

	if !utils.IsEmpty(params.RemindTime) {
		updateParams["remindTime"] = params.RemindTime
	}

	doc := bson.M{
		"$set": updateParams,
	}

	_, err = r.collection.UpdateOne(ctx, filter, doc)
	if err != nil {
		log.Println("update reminder fail: ", err)
		return
	}
	return
}

func (r *RemindersRepository) DeleteReminder(ctx context.Context, params *models.DeleteReminderParams) (err error) {
	// 轉換 ID 為 ObjectID
	objID, err := bson.ObjectIDFromHex(params.ID)
	if err != nil {
		log.Println("delete reminders fail: ", err)
		return err
	}

	filter := bson.M{
		"_id":    objID,
		"userId": params.UserID,
	}

	doc := bson.M{
		"$set": bson.M{
			"updatedAt": time.Now(),
			"deleted":   true,
		},
	}

	_, err = r.collection.UpdateOne(ctx, filter, doc)
	if err != nil {
		log.Println("delete reminder fail: ", err)
		return
	}
	return
}

func (r *RemindersRepository) SearchUserReminders(ctx context.Context, params types.SearchUserRemindersParams) ([]models.ReminderModel, error) {
	filter := bson.M{}
	// TODO 搜尋提醒時間
	if !utils.IsEmpty(params.UserId) {
		filter["userId"] = params.UserId
	}

	if !utils.IsEmpty(params.CreateStartTime) {
		filter["createdAt"] = bson.M{
			"$gt": time.Unix(params.CreateStartTime, 0),
		}
	}

	if !utils.IsEmpty(params.CreateEndTime) {
		filter["createdAt"] = bson.M{
			"$lt": time.Unix(params.CreateEndTime, 0),
		}
	}

	if !utils.IsEmpty(params.Frequency) {
		filter["frequency"] = params.Frequency
	}

	filter["deleted"] = false

	page := 1
	if params.Page != nil {
		page = *params.Page
	}
	options := options.Find()

	offset := (page - 1) * consts.PageSize
	options.SetLimit(int64(consts.PageSize))
	options.SetSkip(int64(offset))
	options.SetSort(bson.M{"createdAt": -1})

	cursor, err := r.collection.Find(ctx, filter, options)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []models.ReminderModel
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}

func (r *RemindersRepository) SearchReminderNotifications(ctx context.Context, remindTime models.RemindTime) ([]models.ReminderModel, error) {

	filter := bson.M{
		"deleted":           false,
		"remindTime.hour":   *remindTime.Hour,
		"remindTime.minute": *remindTime.Minute,
		"$expr": bson.M{
			"$switch": bson.M{
				"branches": bson.A{
					// Once
					bson.M{
						"case": bson.M{"$eq": bson.A{"$frequency", models.EnumRemindFrequencyOnce}},
						"then": bson.M{"$and": bson.A{
							bson.M{"$eq": bson.A{"$remindTime.year", *remindTime.Year}},
							bson.M{"$eq": bson.A{"$remindTime.month", *remindTime.Month}},
							bson.M{"$eq": bson.A{"$remindTime.date", *remindTime.Date}},
						}},
					},
					// Daily
					bson.M{
						"case": bson.M{"$eq": bson.A{"$frequency", models.EnumRemindFrequencyDaily}},
						"then": true,
					},
					// Weekly
					bson.M{
						"case": bson.M{"$eq": bson.A{"$frequency", models.EnumRemindFrequencyWeekly}},
						"then": bson.M{"$and": bson.A{
							bson.M{"$eq": bson.A{"$remindTime.weekday", *remindTime.Weekday}},
						}},
					},
					// Monthly
					bson.M{
						"case": bson.M{"$eq": bson.A{"$frequency", models.EnumRemindFrequencyMonthly}},
						"then": bson.M{"$and": bson.A{
							bson.M{"$eq": bson.A{"$remindTime.date", *remindTime.Date}},
						}},
					},
					// Annually
					bson.M{
						"case": bson.M{"$eq": bson.A{"$frequency", models.EnumRemindFrequencyAnnually}},
						"then": bson.M{"$and": bson.A{
							bson.M{"$eq": bson.A{"$remindTime.monthly", *remindTime.Month}},
							bson.M{"$eq": bson.A{"$remindTime.date", *remindTime.Date}},
						}},
					},
				},
				"default": false,
			},
		},
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []models.ReminderModel
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}
