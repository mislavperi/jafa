package controllers

import (
	"errors"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/mislavperi/jafa/server/internal/domain/apperr"
	"github.com/mislavperi/jafa/server/internal/domain/models"
	requestmodels "github.com/mislavperi/jafa/server/internal/domain/models/request"
	"github.com/mislavperi/jafa/server/internal/server/httperr"
	"github.com/mislavperi/jafa/server/internal/server/middleware"
)

type AuthController struct {
	authService AuthService
}

func NewAuthController(authService AuthService) *AuthController {
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
			httperr.BadRequest(ctx, "username and password are required", err)
			return
		}
		user, err := ac.authService.Login(ctx.Request.Context(), req.Username, req.Password)
		if errors.Is(err, apperr.ErrInvalidCredentials) {
			httperr.Unauthorized(ctx, "invalid username or password")
			return
		}
		if err != nil {
			httperr.Internal(ctx, err)
			return
		}
		session := sessions.Default(ctx)
		saveUserToSession(session, user)
		if err := session.Save(); err != nil {
			httperr.Respond(ctx, http.StatusInternalServerError, "failed to create session", err)
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
			httperr.Respond(ctx, http.StatusInternalServerError, "failed to clear session", err)
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
			httperr.Unauthorized(ctx, "not authenticated")
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

// DeleteAccount soft-deletes the authenticated user (the account can no longer
// log in; its data is kept in the database for recovery), then clears the
// session.
func (ac *AuthController) DeleteAccount() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		session := sessions.Default(ctx)
		userID, ok := session.Get(middleware.SessionUserIDKey).(int64)
		if !ok {
			httperr.Unauthorized(ctx, "not authenticated")
			return
		}
		if err := ac.authService.DeleteAccount(ctx.Request.Context(), userID); err != nil {
			if errors.Is(err, apperr.ErrUserNotFound) {
				httperr.NotFound(ctx, "user not found")
				return
			}
			httperr.Internal(ctx, err)
			return
		}
		session.Clear()
		session.Options(sessions.Options{MaxAge: -1})
		if err := session.Save(); err != nil {
			httperr.Respond(ctx, http.StatusInternalServerError, "failed to clear session", err)
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": "account deleted"})
	}
}

func (ac *AuthController) Register() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req requestmodels.RegisterRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			httperr.BadRequest(ctx, "username and password are required", err)
			return
		}
		user, err := ac.authService.Register(ctx.Request.Context(), requestmodels.RegisterRequest{
			Username:  req.Username,
			Password:  req.Password,
			FirstName: req.FirstName,
			LastName:  req.LastName,
			Email:     req.Email,
		})
		if errors.Is(err, apperr.ErrUsernameTaken) {
			httperr.Conflict(ctx, "username is already taken")
			return
		}
		if err != nil {
			httperr.Internal(ctx, err)
			return
		}
		session := sessions.Default(ctx)
		saveUserToSession(session, user)
		if err := session.Save(); err != nil {
			httperr.Respond(ctx, http.StatusInternalServerError, "failed to create session", err)
			return
		}
		ctx.JSON(http.StatusCreated, user)
	}
}
