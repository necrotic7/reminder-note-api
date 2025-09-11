package services

import (
	"context"
	"log"

	"go.uber.org/fx"
)

var Module = fx.Module(
	"services",
	fx.Provide(
		NewScheduleService,
		NewLineBotService,
		NewLineWebhookService,
		NewReminderService,
		NewUsersService,
	),
	fx.Invoke(func(lc fx.Lifecycle, s *ScheduleService) {
		lc.Append(fx.Hook{
			OnStart: func(ctx context.Context) error {
				log.Println("🟢 Starting ScheduleService")
				return s.Start(ctx)
			},
			OnStop: func(ctx context.Context) error {
				log.Println("🛑 ScheduleService shutdown")
				return s.Stop(ctx)
			},
		})
	}),
)
