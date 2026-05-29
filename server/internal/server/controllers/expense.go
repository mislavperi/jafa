package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mislavperi/jafa/server/internal/domain/models"
	requestmodels "github.com/mislavperi/jafa/server/internal/domain/models/request"
	"github.com/mislavperi/jafa/server/internal/domain/services"
	"github.com/mislavperi/jafa/server/internal/server/httperr"
	"github.com/mislavperi/jafa/server/internal/server/middleware"
)

type firstExpenseDateResponse struct {
	FirstDate string `json:"firstDate"`
}

type ExpenseController struct {
	expenseService *services.ExpenseService
}

func NewExpenseController(expenseService *services.ExpenseService) *ExpenseController {
	return &ExpenseController{
		expenseService: expenseService,
	}
}

func requireUser(ctx *gin.Context) (int64, bool) {
	uid, ok := middleware.CurrentUserID(ctx)
	if !ok {
		httperr.Unauthorized(ctx, "authentication required")
		return 0, false
	}
	return uid, true
}

func (ec *ExpenseController) CreateExpense() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		uid, ok := requireUser(ctx)
		if !ok {
			return
		}
		var req requestmodels.CreateExpenseRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			httperr.BadRequest(ctx, err.Error(), err)
			return
		}
		var recurringSchedule *models.RecurringSchedule
		if req.RecurringSchedule != nil {
			recurringSchedule = &models.RecurringSchedule{
				Interval:   models.RecurrenceInterval(req.RecurringSchedule.Interval),
				DayOfMonth: req.RecurringSchedule.DayOfMonth,
				StartDate:  req.RecurringSchedule.StartDate,
			}
		}
		expense, err := ec.expenseService.CreateExpense(services.CreateExpenseInput{
			UserID:            uid,
			Name:              req.Name,
			Amount:            *req.Amount,
			Cost:              *req.Cost,
			RecurringSchedule: recurringSchedule,
		})
		if err != nil {
			httperr.Internal(ctx, err)
			return
		}
		ctx.JSON(http.StatusCreated, expense)
	}
}

func (ec *ExpenseController) GetAllExpenses() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		uid, ok := requireUser(ctx)
		if !ok {
			return
		}
		expenses, err := ec.expenseService.GetAllExpenses(uid)
		if err != nil {
			httperr.Internal(ctx, err)
			return
		}
		ctx.JSON(http.StatusOK, expenses)
	}
}

func (ec *ExpenseController) GetTotalSpendThisMonth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		uid, ok := requireUser(ctx)
		if !ok {
			return
		}
		total, err := ec.expenseService.GetTotalSpendThisMonth(uid)
		if err != nil {
			httperr.Internal(ctx, err)
			return
		}
		ctx.JSON(http.StatusOK, total)
	}
}

func (ec *ExpenseController) GetDailySpend() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		uid, ok := requireUser(ctx)
		if !ok {
			return
		}
		v, err := strconv.Atoi(ctx.Query("months"))
		if err != nil {
			httperr.BadRequest(ctx, "invalid request", err)
			return
		}
		dailySpend, err := ec.expenseService.GetDailySpend(uid, int32(v))
		if err != nil {
			httperr.Internal(ctx, err)
			return
		}
		ctx.JSON(http.StatusOK, dailySpend)
	}
}

func (ec *ExpenseController) GetFirstExpenseDate() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		uid, ok := requireUser(ctx)
		if !ok {
			return
		}
		date, err := ec.expenseService.GetFirstExpenseDate(uid)
		if err != nil {
			httperr.Internal(ctx, err)
			return
		}
		ctx.JSON(http.StatusOK, firstExpenseDateResponse{FirstDate: date})
	}
}

func (ec *ExpenseController) GetDailySpendForMonth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		uid, ok := requireUser(ctx)
		if !ok {
			return
		}
		year, err := strconv.Atoi(ctx.Query("year"))
		if err != nil {
			httperr.BadRequest(ctx, "invalid request", err)
			return
		}
		month, err := strconv.Atoi(ctx.Query("month"))
		if err != nil {
			httperr.BadRequest(ctx, "invalid request", err)
			return
		}
		dailySpend, err := ec.expenseService.GetDailySpendForMonth(uid, int32(year), int32(month))
		if err != nil {
			httperr.Internal(ctx, err)
			return
		}
		ctx.JSON(http.StatusOK, dailySpend)
	}
}

func (ec *ExpenseController) GetExpensesByMonth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		uid, ok := requireUser(ctx)
		if !ok {
			return
		}
		year, err := strconv.Atoi(ctx.Query("year"))
		if err != nil {
			httperr.BadRequest(ctx, "invalid request", err)
			return
		}
		month, err := strconv.Atoi(ctx.Query("month"))
		if err != nil {
			httperr.BadRequest(ctx, "invalid request", err)
			return
		}
		expenses, err := ec.expenseService.GetExpensesByMonth(uid, int32(year), int32(month))
		if err != nil {
			httperr.Internal(ctx, err)
			return
		}
		ctx.JSON(http.StatusOK, expenses)
	}
}

func (ec *ExpenseController) UpdateExpense() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		uid, ok := requireUser(ctx)
		if !ok {
			return
		}
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			httperr.BadRequest(ctx, err.Error(), err)
			return
		}
		var req requestmodels.UpdateExpenseRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			httperr.BadRequest(ctx, err.Error(), err)
			return
		}
		var recurringSchedule *models.RecurringSchedule
		if req.RecurringSchedule != nil {
			recurringSchedule = &models.RecurringSchedule{
				Interval:   models.RecurrenceInterval(req.RecurringSchedule.Interval),
				DayOfMonth: req.RecurringSchedule.DayOfMonth,
				StartDate:  req.RecurringSchedule.StartDate,
			}
		}
		expense, err := ec.expenseService.UpdateExpense(services.UpdateExpenseInput{
			ID:                int64(id),
			UserID:            uid,
			Name:              req.Name,
			Amount:            *req.Amount,
			Cost:              *req.Cost,
			RecurringSchedule: recurringSchedule,
		})
		if err != nil {
			httperr.Internal(ctx, err)
			return
		}
		ctx.JSON(http.StatusOK, expense)
	}
}

func (ec *ExpenseController) DeleteExpense() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		uid, ok := requireUser(ctx)
		if !ok {
			return
		}
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			httperr.BadRequest(ctx, err.Error(), err)
			return
		}
		if err := ec.expenseService.DeleteExpense(uid, int64(id)); err != nil {
			httperr.Internal(ctx, err)
			return
		}
		ctx.Status(http.StatusNoContent)
	}
}

func (ec *ExpenseController) GetExpenseById() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		uid, ok := requireUser(ctx)
		if !ok {
			return
		}
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			httperr.BadRequest(ctx, "invalid request", err)
			return
		}
		expense, err := ec.expenseService.GetById(uid, int64(id))
		if err != nil {
			httperr.Internal(ctx, err)
			return
		}
		ctx.JSON(http.StatusOK, expense)
	}
}
