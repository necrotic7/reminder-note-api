package services

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/zivwu/reminder-note-api/internal/models"
	"github.com/zivwu/reminder-note-api/internal/repositories"
	"github.com/zivwu/reminder-note-api/internal/types"
	"github.com/zivwu/reminder-note-api/internal/utils"
)

type ReminderService struct {
	ReminderRepo   *repositories.RemindersRepository
	LineBotService *LineBotService
}

func NewReminderService(reminderRepo *repositories.RemindersRepository, lineBotService *LineBotService) *ReminderService {
	return &ReminderService{
		ReminderRepo:   reminderRepo,
		LineBotService: lineBotService,
	}
}

func (s *ReminderService) CreateReminderFlow(ctx context.Context, req *types.ReqCreateReminderBody) (err error) {
	err = s.ValidationCreateReminderReq(req)
	if err != nil {
		log.Println("檢查創建 Reminder 參數失敗：", err)
		return
	}
	params := models.InsertReminderParams{
		UserID:     req.UserID,
		Title:      req.Title,
		Content:    req.Content,
		Frequency:  req.Frequency,
		RemindTime: models.RemindTime(req.RemindTime),
	}
	err = s.ReminderRepo.InsertReminder(ctx, &params)
	if err != nil {
		log.Println("創建 Reminder 失敗：", err)
		return
	}
	return
}

func (s *ReminderService) ValidationCreateReminderReq(req *types.ReqCreateReminderBody) (err error) {
	now := time.Now()
	timeString := ""
	switch req.Frequency {
	case models.EnumRemindFrequencyOnce:
		if utils.IsEmpty(req.RemindTime.Year, req.RemindTime.Month, req.RemindTime.Date, req.RemindTime.Hour, req.RemindTime.Minute) {
			return fmt.Errorf("輸入時間不合法，需輸入年/月/日/時/分")
		}
		timeString = fmt.Sprintf("%d-%d-%d %d:%d:00", *req.RemindTime.Year, *req.RemindTime.Month, *req.RemindTime.Date, *req.RemindTime.Hour, *req.RemindTime.Minute)

	case models.EnumRemindFrequencyDaily:
		if utils.IsEmpty(req.RemindTime.Hour, req.RemindTime.Minute) {
			return fmt.Errorf("輸入時間不合法，需輸入時/分")
		}
		timeString = fmt.Sprintf("%d-%d-%d %d:%d:00", now.Year(), int(now.Month()), now.Day(), *req.RemindTime.Hour, *req.RemindTime.Minute)

	case models.EnumRemindFrequencyWeekly:
		if utils.IsEmpty(req.RemindTime.Weekday, req.RemindTime.Hour, req.RemindTime.Minute) {
			return fmt.Errorf("輸入時間不合法，需輸入星期/時/分")
		}
		if *req.RemindTime.Weekday < 0 || *req.RemindTime.Weekday > 6 {
			return fmt.Errorf("輸入的星期不合法：%d", *req.RemindTime.Weekday)
		}
		timeString = fmt.Sprintf("%d-%d-%d %d:%d:00", now.Year(), int(now.Month()), now.Day(), *req.RemindTime.Hour, *req.RemindTime.Minute)

	case models.EnumRemindFrequencyMonthly:
		if utils.IsEmpty(req.RemindTime.Date, req.RemindTime.Hour, req.RemindTime.Minute) {
			return fmt.Errorf("輸入時間不合法，需輸入日/時/分")
		}
		timeString = fmt.Sprintf("%d-%d-%d %d:%d:00", now.Year(), int(now.Month()), *req.RemindTime.Date, *req.RemindTime.Hour, *req.RemindTime.Minute)

	case models.EnumRemindFrequencyAnnually:
		if utils.IsEmpty(req.RemindTime.Month, req.RemindTime.Date, req.RemindTime.Hour, req.RemindTime.Minute) {
			return fmt.Errorf("輸入時間不合法，需輸入月/日/時/分")
		}
		timeString = fmt.Sprintf("%d-%d-%d %d:%d:00", now.Year(), *req.RemindTime.Month, *req.RemindTime.Date, *req.RemindTime.Hour, *req.RemindTime.Minute)

	default:
		return fmt.Errorf("不合法的 Reminder Frequency: %s", req.Frequency)
	}

	resultTime, err := time.Parse("2006-1-2 15:4:5", timeString)
	if err != nil {
		return fmt.Errorf("time.Parse失敗：%w", err)
	}

	if req.Frequency == models.EnumRemindFrequencyOnce && resultTime.Before(time.Now()) {
		return fmt.Errorf("創建單次提醒時，提醒時間(%s)不可小於現在", resultTime.String())
	}

	return
}

func (s *ReminderService) GetUserReminders(ctx context.Context, req types.ReqGetUserRemindersQuery) (resp types.RespGetUserRemindersBody, err error) {
	if utils.IsEmpty(req.UserId) {
		return resp, fmt.Errorf("missing userId")
	}

	result, counts, err := s.ReminderRepo.SearchUserReminders(ctx, types.SearchUserRemindersParams(req))
	if err != nil {
		return resp, fmt.Errorf("search User Reminders fail: %w", err)
	}
	resp.Counts = counts
	resp.Records = result
	return
}

func (s *ReminderService) ReminderScheduler(ctx context.Context) {
	now := time.Now()
	year := now.Year()
	month := now.Month()
	date := now.Day()
	weekday := now.Weekday()
	hour := now.Hour()
	minute := now.Minute()
	remindTime := models.RemindTime{
		Year:    &year,
		Month:   (*int)(&month),
		Date:    &date,
		Weekday: (*int)(&weekday),
		Hour:    &hour,
		Minute:  &minute,
	}
	log.Printf("[ReminderScheduler] Reminder撈取時間範圍：%v/%v/%v (%v) %v:%v \n", year, month, date, weekday, hour, minute)
	// 撈出要通知的對象
	result, err := s.ReminderRepo.SearchReminderNotifications(ctx, remindTime)
	if err != nil {
		log.Println("reminder推播排程執行失敗：", err)
		return
	}

	log.Printf("共有 %v 筆Reminder推播需要發送\n", len(result))

	if utils.IsEmpty(result) {
		return
	}

	for _, r := range result {
		messages := []linebot.SendingMessage{
			linebot.NewTextMessage(fmt.Sprintf("提醒事項：%v\n%v", r.Title, r.Content)),
		}
		s.LineBotService.PushNotifyMessage(&types.PushMessageParams{
			UserId:   r.UserID,
			Messages: messages,
		})
	}
}

func (s *ReminderService) UpdateReminderFlow(ctx context.Context, req *types.ReqUpdateReminderBody) (err error) {
	payload, err := utils.StructConvert[types.ReqCreateReminderBody](req)
	if err != nil {
		log.Println("update reminder struct convert fail:", err)
		return
	}
	err = s.ValidationCreateReminderReq(payload)
	if err != nil {
		log.Println("檢查更新 Reminder 參數失敗：", err)
		return
	}
	params := models.UpdateReminderParams{
		ID:         req.ID,
		UserID:     req.UserID,
		Title:      req.Title,
		Content:    req.Content,
		Frequency:  req.Frequency,
		RemindTime: models.RemindTime(req.RemindTime),
	}
	err = s.ReminderRepo.UpdateReminder(ctx, &params)
	if err != nil {
		log.Println("更新 Reminder 失敗：", err)
		return
	}
	return
}

func (s *ReminderService) DeleteReminderFlow(ctx context.Context, req *types.ReqDeleteReminderBody) (err error) {
	params := models.DeleteReminderParams{
		ID:     req.ID,
		UserID: req.UserID,
	}
	err = s.ReminderRepo.DeleteReminder(ctx, &params)
	if err != nil {
		log.Println("刪除 Reminder 失敗：", err)
		return
	}
	return
}
