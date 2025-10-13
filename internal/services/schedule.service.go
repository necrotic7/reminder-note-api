package services

import (
	"context"
	"log"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/zivwu/reminder-note-api/internal/repositories"
)

type ScheduleService struct {
	cron            *cron.Cron
	reminderService *ReminderService
}

func NewScheduleService(reminderRepo *repositories.RemindersRepository, lineBotService *LineBotService) *ScheduleService {
	c := cron.New(cron.WithSeconds())
	s := ScheduleService{
		cron:            c,
		reminderService: NewReminderService(reminderRepo, lineBotService),
	}
	return &s
}

func (s *ScheduleService) registerJobs() {
	// 這邊可以註冊你所有的任務
	s.cron.AddFunc("0 * * * * *", func() {
		log.Println("Reminder 推播排程任務啟動")

		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()
		s.reminderService.ReminderScheduler(ctx)

	})

	log.Println("🟢 Register Job Complete")
}

func (s *ScheduleService) Start(ctx context.Context) error {
	s.registerJobs()
	s.cron.Start()
	return nil
}

func (s *ScheduleService) Stop(ctx context.Context) error {
	stopCtx := s.cron.Stop()
	select {
	case <-stopCtx.Done():
		log.Println("🛑 Schedule 已停止")
	case <-ctx.Done():
		log.Println("🟡 Schedule 逾時")
	}
	return nil
}
