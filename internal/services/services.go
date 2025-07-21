package services

import "go.uber.org/fx"

var Module = fx.Module(
	"service",
	fx.Provide(
		NewLineBotService,
		NewLineWebhookService,
		NewReminderService,
	),
)
