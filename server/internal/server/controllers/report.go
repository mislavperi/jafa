package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mislavperi/jafa/server/internal/server/httperr"
)

type ReportController struct {
	reportService ReportService
}

func NewReportController(reportService ReportService) *ReportController {
	return &ReportController{reportService: reportService}
}

func (rc *ReportController) CategoryBreakdown() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		uid, ok := requireUser(ctx)
		if !ok {
			return
		}
		breakdown, err := rc.reportService.CategoryBreakdown(ctx.Request.Context(), uid)
		if err != nil {
			httperr.Internal(ctx, err)
			return
		}
		ctx.JSON(http.StatusOK, breakdown)
	}
}

func (rc *ReportController) MonthlySpend() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		uid, ok := requireUser(ctx)
		if !ok {
			return
		}
		monthly, err := rc.reportService.MonthlySpend(ctx.Request.Context(), uid)
		if err != nil {
			httperr.Internal(ctx, err)
			return
		}
		ctx.JSON(http.StatusOK, monthly)
	}
}
