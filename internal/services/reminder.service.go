package services

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/zivwu/reminder-note-api/internal/models"
	"github.com/zivwu/reminder-note-api/internal/repositories"
	"github.com/zivwu/reminder-note-api/internal/types"
	"github.com/zivwu/reminder-note-api/internal/utils"
)

type ReminderService struct {
	ReminderRepo *repositories.ReminderRepository
}

func NewReminderService(reminderRepo *repositories.ReminderRepository) *ReminderService {
	return &ReminderService{
		ReminderRepo: reminderRepo,
	}
}

func (s *ReminderService) CreateReminderFlow(ctx context.Context, req types.ReqCreateReminderBody) (err error) {
	err = s.ValidationCreateReminderReq(req)
	if err != nil {
		log.Println("檢查創建 Reminder 參數失敗：", err)
		return
	}
	err = s.ReminderRepo.InsertReminder(ctx, req)
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

func (s *ReminderService) GetUserReminders(ctx context.Context, req types.ReqGetUserRemindersQuery) ([]models.Reminder, error) {
	if utils.IsEmpty(req.UserId) {
		return nil, fmt.Errorf("missing userId")
	}
	result, err := s.ReminderRepo.SearchUserReminders(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("search User Reminders fail: %w", err)
	}
	return result, nil
}
