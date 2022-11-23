package services

import "go.uber.org/fx"

var Bootstrap = fx.Options(
	fx.Provide(NewQueueService),
	fx.Provide(NewSocketService),
)
