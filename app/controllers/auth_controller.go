package controllers

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sayopaul/sendchamp-go-test/app/middlewares"
	"github.com/sayopaul/sendchamp-go-test/common"
	"github.com/sayopaul/sendchamp-go-test/models"
	"github.com/sayopaul/sendchamp-go-test/repositories"
	"gorm.io/gorm"
)

type AuthController struct {
	userRepository repositories.UserRepository
	jwtMiddleware  middlewares.JWTMiddleware
}
type login struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func NewAuthController(userRepository repositories.UserRepository, jwtMiddleware middlewares.JWTMiddleware) AuthController {
	return AuthController{
		userRepository: userRepository,
		jwtMiddleware:  jwtMiddleware,
	}
}

func (ac AuthController) SignInHandler(c *gin.Context) {
	// validate sign in details
	var signInDetails login
	if err := c.ShouldBindJSON(&signInDetails); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	// get user with email and hashed password
	hashedPassword := common.Sha256HashString(signInDetails.Password)
	user, err := ac.userRepository.CheckSignIn(signInDetails.Email, hashedPassword)
	if err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Incorrect Email/Password",
			"data":    nil,
		})
		return
	}

	//update last login time
	_, err = ac.userRepository.Update(user.ID, models.User{LastLogin: time.Now()})

	if err != nil {
		log.Panicf("%v : %v", "Error updating login  time", err)
	}
	//generate a new jwt token
	token, expire, err := ac.jwtMiddleware.GinJWTMiddleware.TokenGenerator(user)

	if err != nil {
		log.Panicf("%v : %v", "Error Generating JWT", err)
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Sign in successful.",
		"data": gin.H{
			"token":  token,
			"expire": expire,
		},
	})

}

func (ac AuthController) SignupHandler(c *gin.Context) {
	// validate signup details
	var signupDetails login
	if err := c.ShouldBindJSON(&signupDetails); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	// check if a user with this email already exists
	_, err := ac.userRepository.GetOne("email", signupDetails.Email)
	if !errors.Is(err, gorm.ErrRecordNotFound) && err == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "email has been taken.",
			"data":    nil,
		})
		return
	}

	// hash password
	hashedPassword := common.Sha256HashString(signupDetails.Password)
	newUser := models.User{
		Email:     signupDetails.Email,
		Password:  string(hashedPassword),
		LastLogin: time.Now(),
	}
	user, err := ac.userRepository.CreateUser(newUser)
	if err != nil {
		log.Fatalf("Error creating User: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Something went wrong, Try again later",
			"data":    nil,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Sign up successful.",
		"data": gin.H{
			"user": user,
		},
	})
}
