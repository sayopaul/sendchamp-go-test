package routes

import (
	"github.com/sayopaul/sendchamp-go-test/app/middlewares"
	"github.com/sayopaul/sendchamp-go-test/infrastructure"
)

// Routes contains multiple routes (that implement the IRoute interface)
type Routes []IRoute

// IRoute : Route interface
type IRoute interface {
	Setup()
}

// handler is an instance of *gin.Enginer
type Route struct {
	handler infrastructure.Router
}

type JwtAuthRoute struct {
	Route
	jwtMiddleware middlewares.JWTMiddleware
}

// for routes that do not use authentication
func NewRoute(handler infrastructure.Router) Route {
	return Route{
		handler: handler,
	}
}

func NewJwtAuthRoute(route Route, jwtMiddleware middlewares.JWTMiddleware) JwtAuthRoute {
	return JwtAuthRoute{route, jwtMiddleware}
}

// NewRoutes sets up routes
func NewRoutes(
	authRoutes AuthRoutes,
) Routes {
	return Routes{
		authRoutes,
	}
}

// Setup all the route
func (r Routes) Setup() {
	for _, route := range r {
		route.Setup()
	}
}
