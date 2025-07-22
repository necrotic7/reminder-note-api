package services

import (
	"context"
	"log"

	"github.com/robfig/cron/v3"
)

type ScheduleService struct {
	cron *cron.Cron
}

func NewScheduleService() *ScheduleService {
	c := cron.New(cron.WithSeconds())
	s := ScheduleService{
		cron: c,
	}
	return &s
}

func (s *ScheduleService) registerJobs() {
	// 這邊可以註冊你所有的任務
	s.cron.AddFunc("0 * * * * *", func() {
		log.Println("Reminder 推播排程任務啟動")
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
