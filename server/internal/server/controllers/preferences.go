package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mislavperi/jafa/server/internal/domain/dto"
	"github.com/mislavperi/jafa/server/internal/domain/models"
	requestmodels "github.com/mislavperi/jafa/server/internal/domain/models/request"
	"github.com/mislavperi/jafa/server/internal/server/httperr"
)

type PreferencesController struct {
	preferencesService PreferencesService
}

func NewPreferencesController(preferencesService PreferencesService) *PreferencesController {
	return &PreferencesController{preferencesService: preferencesService}
}

func (pc *PreferencesController) Get() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		uid, ok := requireUser(ctx)
		if !ok {
			return
		}
		prefs, err := pc.preferencesService.Get(ctx.Request.Context(), uid)
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
		var req requestmodels.UpsertPreferencesRequest
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
		// Optional fields fall back to the user's stored values (or defaults when
		// no row exists yet) so clients that omit them don't silently reset them.
		stored, err := pc.preferencesService.Get(ctx.Request.Context(), uid)
		if err != nil {
			httperr.Internal(ctx, err)
			return
		}
		if req.WeekStart == "" {
			req.WeekStart = stored.WeekStart
		}
		if !models.ValidWeekStarts[req.WeekStart] {
			httperr.BadRequest(ctx, "invalid weekStart", nil)
			return
		}
		monthlyBudget := stored.MonthlyBudget
		if req.MonthlyBudget != nil {
			monthlyBudget = *req.MonthlyBudget
		}
		if monthlyBudget < 0 {
			httperr.BadRequest(ctx, "monthlyBudget must not be negative", nil)
			return
		}
		notifyWeeklySummary := stored.NotifyWeeklySummary
		if req.NotifyWeeklySummary != nil {
			notifyWeeklySummary = *req.NotifyWeeklySummary
		}
		notifyBudgetAlerts := stored.NotifyBudgetAlerts
		if req.NotifyBudgetAlerts != nil {
			notifyBudgetAlerts = *req.NotifyBudgetAlerts
		}
		notifyProductUpdates := stored.NotifyProductUpdates
		if req.NotifyProductUpdates != nil {
			notifyProductUpdates = *req.NotifyProductUpdates
		}
		prefs, err := pc.preferencesService.Upsert(ctx.Request.Context(), dto.UpsertPreferencesInput{
			UserID:               uid,
			AccentID:             req.AccentID,
			FontSize:             req.FontSize,
			DarkMode:             req.DarkMode,
			Currency:             req.Currency,
			WeekStart:            req.WeekStart,
			MonthlyBudget:        monthlyBudget,
			NotifyWeeklySummary:  notifyWeeklySummary,
			NotifyBudgetAlerts:   notifyBudgetAlerts,
			NotifyProductUpdates: notifyProductUpdates,
		})
		if err != nil {
			httperr.Internal(ctx, err)
			return
		}
		ctx.JSON(http.StatusOK, prefs)
	}
}
