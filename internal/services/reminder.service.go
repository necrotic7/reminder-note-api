package services

import (
	"context"
	"fmt"
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

type ReminderService struct {
	DB *mongo.Client
}

func NewReminderService(db *mongo.Client) *ReminderService {
	return &ReminderService{
		DB: db,
	}
}

func (s *ReminderService) CreateReminderFlow(ctx context.Context, req types.ReqCreateReminderBody) (err error) {
	err = s.ValidationCreateReminderReq(req)
	if err != nil {
		log.Println("檢查創建 Reminder 參數失敗：", err)
		return
	}
	err = s.InsertReminder(ctx, req)
	if err != nil {
		log.Println("創建 Reminder 失敗：", err)
		return
	}
	return
}

func (s *ReminderService) ValidationCreateReminderReq(req types.ReqCreateReminderBody) (err error) {
	switch req.Frequency {
	case models.EnumRemindFrequencyOnce:
		if utils.IsEmpty(req.RemindTime.Year, req.RemindTime.Month, req.RemindTime.Date, req.RemindTime.Hour, req.RemindTime.Minute) {
			return fmt.Errorf("輸入時間不合法，需輸入年/月/日/時/分")
		}
		if *req.RemindTime.Year < time.Now().Year() {
			return fmt.Errorf("新創建的年份(%v)不可小於今年年份(%v)", *req.RemindTime.Year, time.Now().Year())
		}
	case models.EnumRemindFrequencyDaily:
		if utils.IsEmpty(req.RemindTime.Hour, req.RemindTime.Minute) {
			return fmt.Errorf("輸入時間不合法，需輸入時/分")
		}
	case models.EnumRemindFrequencyWeekly:
		if utils.IsEmpty(req.RemindTime.Weekday, req.RemindTime.Hour, req.RemindTime.Minute) {
			return fmt.Errorf("輸入時間不合法，需輸入星期/時/分")
		}
		if *req.RemindTime.Weekday < 1 || *req.RemindTime.Weekday > 7 {
			return fmt.Errorf("輸入的星期不合法：%d", *req.RemindTime.Weekday)
		}
	case models.EnumRemindFrequencyMonthly:
		if utils.IsEmpty(req.RemindTime.Date, req.RemindTime.Hour, req.RemindTime.Minute) {
			return fmt.Errorf("輸入時間不合法，需輸入日/時/分")
		}
	case models.EnumRemindFrequencyAnnually:
		if utils.IsEmpty(req.RemindTime.Month, req.RemindTime.Date, req.RemindTime.Hour, req.RemindTime.Minute) {
			return fmt.Errorf("輸入時間不合法，需輸入月/日/時/分")
		}
	default:
		return fmt.Errorf("不合法的 Reminder Frequency: %s", req.Frequency)
	}

	timeString := fmt.Sprintf("%d-%d-%d %d:%d:00", *req.RemindTime.Year, *req.RemindTime.Month, *req.RemindTime.Date, *req.RemindTime.Hour, *req.RemindTime.Minute)
	_, err = time.Parse("2006-1-2 15:4:5", timeString)
	if err != nil {
		return fmt.Errorf("time.Parse失敗：%w", err)
	}

	return
}

func (s *ReminderService) InsertReminder(ctx context.Context, req types.ReqCreateReminderBody) (err error) {
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
		log.Println("insert reminder fail: ", err)
		return err
	}
	return
}

func (s *ReminderService) GetUserReminders(ctx context.Context, params types.ReqGetUserRemindersQuery) ([]models.Reminder, error) {
	collection := s.DB.Database("reminder-note").Collection("reminders")
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
