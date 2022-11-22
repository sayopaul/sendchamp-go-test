package middlewares

import "go.uber.org/fx"

var Bootstrap = fx.Options(
	fx.Provide(NewJWTMiddleware),
)
