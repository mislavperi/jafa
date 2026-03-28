package controllers

import (
	"errors"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/mislavperi/jafa/server/internal/domain/models"
	requestmodels "github.com/mislavperi/jafa/server/internal/domain/models/request"
	"github.com/mislavperi/jafa/server/internal/domain/services"
	"github.com/mislavperi/jafa/server/internal/server/middleware"
	customerrors "github.com/mislavperi/jafa/server/utils/errors"
)

type AuthController struct {
	authService *services.AuthService
}

func NewAuthController(authService *services.AuthService) *AuthController {
	return &AuthController{authService: authService}
}

func saveUserToSession(session sessions.Session, user models.User) {
	session.Set(middleware.SessionUserIDKey, user.Id)
	session.Set("username", user.Username)
	session.Set("first_name", user.FirstName)
	session.Set("last_name", user.LastName)
	session.Set("email", user.Email)
	session.Set("avatar_url", user.AvatarUrl)
}

func (ac *AuthController) Login() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req requestmodels.LoginRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "username and password are required"})
			return
		}
		user, err := ac.authService.Login(req.Username, req.Password)
		if errors.Is(err, customerrors.ErrInvalidCredentials) {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid username or password"})
			return
		}
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}
		session := sessions.Default(ctx)
		saveUserToSession(session, user)
		if err := session.Save(); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create session"})
			return
		}
		ctx.JSON(http.StatusOK, user)
	}
}

func (ac *AuthController) Logout() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		session := sessions.Default(ctx)
		session.Clear()
		session.Options(sessions.Options{MaxAge: -1})
		if err := session.Save(); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to clear session"})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": "logged out"})
	}
}

func (ac *AuthController) Me() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		session := sessions.Default(ctx)
		userID := session.Get(middleware.SessionUserIDKey)
		if userID == nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "not authenticated"})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"id":         userID,
			"username":   session.Get("username"),
			"first_name": session.Get("first_name"),
			"last_name":  session.Get("last_name"),
			"email":      session.Get("email"),
			"avatar_url": session.Get("avatar_url"),
		})
	}
}

func (ac *AuthController) Register() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req requestmodels.RegisterRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "username and password are required"})
			return
		}
		user, err := ac.authService.Register(requestmodels.RegisterRequest{
			Username:  req.Username,
			Password:  req.Password,
			FirstName: req.FirstName,
			LastName:  req.LastName,
			Email:     req.Email,
		})
		if errors.Is(err, customerrors.ErrUsernameTaken) {
			ctx.JSON(http.StatusConflict, gin.H{"error": "username is already taken"})
			return
		}
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}
		session := sessions.Default(ctx)
		saveUserToSession(session, user)
		if err := session.Save(); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create session"})
			return
		}
		ctx.JSON(http.StatusCreated, user)
	}
}
