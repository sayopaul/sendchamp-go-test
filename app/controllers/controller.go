package controllers

import "go.uber.org/fx"

var Bootstrap = fx.Options(
	fx.Provide(NewAuthController),
)
