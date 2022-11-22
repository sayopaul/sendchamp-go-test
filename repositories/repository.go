package repositories

import (
	"go.uber.org/fx"
)

var Bootstrap = fx.Options(
	fx.Provide(NewUserRepository),
)
