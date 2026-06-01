package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mislavperi/jafa/server/internal/domain/models"
	"github.com/mislavperi/jafa/server/internal/domain/services"
	"github.com/mislavperi/jafa/server/internal/server/httperr"
)

type PreferencesController struct {
	preferencesService *services.PreferencesService
}

func NewPreferencesController(preferencesService *services.PreferencesService) *PreferencesController {
	return &PreferencesController{preferencesService: preferencesService}
}

type upsertPreferencesRequest struct {
	AccentID string `json:"accentId" binding:"required"`
	FontSize string `json:"fontSize" binding:"required"`
	DarkMode bool   `json:"darkMode"`
	Currency string `json:"currency"`
}

func (pc *PreferencesController) Get() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		uid, ok := requireUser(ctx)
		if !ok {
			return
		}
		prefs, err := pc.preferencesService.Get(uid)
		if err != nil {
			httperr.Internal(ctx, err)
			return
		}
		ctx.JSON(http.StatusOK, prefs)
	}
}

func (pc *PreferencesController) Upsert() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		uid, ok := requireUser(ctx)
		if !ok {
			return
		}
		var req upsertPreferencesRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			httperr.BadRequest(ctx, err.Error(), err)
			return
		}
		if req.Currency == "" {
			req.Currency = models.DefaultCurrency
		}
		if !models.ValidFontSizes[req.FontSize] {
			httperr.BadRequest(ctx, "invalid fontSize", nil)
			return
		}
		if !models.ValidCurrencies[req.Currency] {
			httperr.BadRequest(ctx, "invalid currency", nil)
			return
		}
		prefs, err := pc.preferencesService.Upsert(models.UpsertPreferencesInput{
			UserID:   uid,
			AccentID: req.AccentID,
			FontSize: req.FontSize,
			DarkMode: req.DarkMode,
			Currency: req.Currency,
		})
		if err != nil {
			httperr.Internal(ctx, err)
			return
		}
		ctx.JSON(http.StatusOK, prefs)
	}
}
