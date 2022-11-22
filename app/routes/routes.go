package routes

import (
	"go.uber.org/fx"
)

var Bootstrap = fx.Options(
	fx.Provide(NewAuthRoutes),
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
) Routes {
	return Routes{
		authRoutes,
	}
}

// Setup each of the route groupings
func (r Routes) Setup() {
	for _, route := range r {
		route.Setup()
	}
}
