package middlewares

import (
	"log"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/sayopaul/sendchamp-go-test/config"
	"github.com/sayopaul/sendchamp-go-test/models"
)

type JWTMiddleware struct {
	GinJWTMiddleware *jwt.GinJWTMiddleware
}

func NewJWTMiddleware(
	configEnv config.Config,
) JWTMiddleware {
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "gin jwt",
		Key:         []byte(configEnv.JWTSecret),
		Timeout:     time.Hour * 24 * 30, // 30 days
		MaxRefresh:  time.Hour * 24 * 60, // 60 days
		IdentityKey: "id",
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if user, ok := data.(models.User); ok {
				return jwt.MapClaims{
					"id":   user.ID,
					"uuid": user.UUID,
				}
			}
			return jwt.MapClaims{}
		},

		Authorizator: func(data interface{}, c *gin.Context) bool {
			return true
		},
		TimeFunc: time.Now,
	})

	if err != nil {
		log.Panic("JWT Error:" + err.Error())
	}

	errInit := authMiddleware.MiddlewareInit()
	if errInit != nil {
		log.Panic("authMiddleware.MiddlewareInit() Error:" + errInit.Error())
	}

	return JWTMiddleware{
		GinJWTMiddleware: authMiddleware,
	}
}
