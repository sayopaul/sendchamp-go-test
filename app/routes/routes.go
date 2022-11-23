package routes

import (
	"go.uber.org/fx"
)

var Bootstrap = fx.Options(
	fx.Provide(NewAuthRoutes),
	fx.Provide(NewTaskRoutes),
	fx.Provide(NewRoutes),
)

// Routes contains multiple routes (that implement the IRoute interface)
type Routes []IRoute

// IRoute : Route interface
type IRoute interface {
	Setup()
}

// NewRoutes sets up routes
func NewRoutes(
	authRoutes AuthRoutes,
	taskRoutes TaskRoutes,
) Routes {
	return Routes{
		authRoutes,
		taskRoutes,
	}
}

// Setup each of the route groupings
func (r Routes) Setup() {
	for _, route := range r {
		route.Setup()
	}
}
