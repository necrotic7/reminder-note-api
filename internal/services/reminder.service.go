package services

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/zivwu/reminder-note-api/internal/models"
	"github.com/zivwu/reminder-note-api/internal/types"
	"github.com/zivwu/reminder-note-api/internal/utils"
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
		if utils.IsNill(req.RemindTime.Year, req.RemindTime.Month, req.RemindTime.Date, req.RemindTime.Hour, req.RemindTime.Minute) {
			return fmt.Errorf("輸入時間不合法，需輸入年/月/日/時/分")
		}
		if *req.RemindTime.Year < time.Now().Year() {
			return fmt.Errorf("新創建的年份(%v)不可小於今年年份(%v)", *req.RemindTime.Year, time.Now().Year())
		}
	case models.EnumRemindFrequencyDaily:
		if utils.IsNill(req.RemindTime.Hour, req.RemindTime.Minute) {
			return fmt.Errorf("輸入時間不合法，需輸入時/分")
		}
	case models.EnumRemindFrequencyWeekly:
		if utils.IsNill(req.RemindTime.Weekday, req.RemindTime.Hour, req.RemindTime.Minute) {
			return fmt.Errorf("輸入時間不合法，需輸入星期/時/分")
		}
		if *req.RemindTime.Weekday < 1 || *req.RemindTime.Weekday > 7 {
			return fmt.Errorf("輸入的星期不合法：%d", *req.RemindTime.Weekday)
		}
	case models.EnumRemindFrequencyMonthly:
		if utils.IsNill(req.RemindTime.Date, req.RemindTime.Hour, req.RemindTime.Minute) {
			return fmt.Errorf("輸入時間不合法，需輸入日/時/分")
		}
	case models.EnumRemindFrequencyAnnually:
		if utils.IsNill(req.RemindTime.Month, req.RemindTime.Date, req.RemindTime.Hour, req.RemindTime.Minute) {
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
