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
	db *mongo.Client
}

func NewRemindersRepository(db *mongo.Client) *RemindersRepository {
	return &RemindersRepository{
		db: db,
	}
}

func (r *RemindersRepository) InsertReminder(ctx context.Context, params models.InsertReminderParams) (err error) {
	collection := r.db.Database("reminder-note").Collection("reminders")
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

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
	_, err = collection.InsertOne(ctx, doc)
	if err != nil {
		log.Println("insert reminder fail: ", err)
		return err
	}
	return
}

func (r *RemindersRepository) SearchUserReminders(ctx context.Context, params types.ReqGetUserRemindersQuery) ([]models.ReminderModel, error) {
	collection := r.db.Database("reminder-note").Collection("reminders")
	filter := bson.M{}
	if !utils.IsEmpty(params.UserId) {
		filter["userId"] = params.UserId
	}

	filter["deleted"] = false

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

	var results []models.ReminderModel
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}

func (r *RemindersRepository) SearchReminderNotifications(ctx context.Context, remindTime models.RemindTime) ([]models.ReminderModel, error) {
	collection := r.db.Database("reminder-note").Collection("reminders")

	filter := bson.M{
		"deleted": false,
		// Daily
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

	cursor, err := collection.Find(ctx, filter)
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
