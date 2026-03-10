package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mislavperi/jafa/server/internal/domain/services"
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

type createExpenseRequest struct {
	Name   string   `json:"name" binding:"required"`
	Amount *float32 `json:"amount" binding:"required"`
	Cost   *float32 `json:"cost" binding:"required"`
}

func (ec *ExpenseController) CreateExpense() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req createExpenseRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		expense, err := ec.expenseService.CreateExpense(services.CreateExpenseInput{
			Name:   req.Name,
			Amount: *req.Amount,
			Cost:   *req.Cost,
		})
		if err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		ctx.JSON(http.StatusCreated, expense)
	}
}

func (ec *ExpenseController) GetAllExpenses() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		expenses, err := ec.expenseService.GetAllExpenses()
		if err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, err)
		}
		ctx.JSON(http.StatusOK, expenses)
	}
}

func (ec *ExpenseController) GetTotalSpendThisMonth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		total, err := ec.expenseService.GetTotalSpendThisMonth()
		if err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		ctx.JSON(http.StatusOK, total)
	}
}

func (ec *ExpenseController) GetDailySpend() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		v, err := strconv.Atoi(ctx.Query("months"))
		if err != nil {
			ctx.AbortWithError(http.StatusBadRequest, err)
			return
		}
		dailySpend, err := ec.expenseService.GetDailySpend(int32(v))
		if err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		ctx.JSON(http.StatusOK, dailySpend)
	}
}

func (ec *ExpenseController) GetFirstExpenseDate() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		date, err := ec.expenseService.GetFirstExpenseDate()
		if err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		ctx.JSON(http.StatusOK, firstExpenseDateResponse{FirstDate: date})
	}
}

func (ec *ExpenseController) GetDailySpendForMonth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		year, err := strconv.Atoi(ctx.Query("year"))
		if err != nil {
			ctx.AbortWithError(http.StatusBadRequest, err)
			return
		}
		month, err := strconv.Atoi(ctx.Query("month"))
		if err != nil {
			ctx.AbortWithError(http.StatusBadRequest, err)
			return
		}
		dailySpend, err := ec.expenseService.GetDailySpendForMonth(int32(year), int32(month))
		if err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		ctx.JSON(http.StatusOK, dailySpend)
	}
}

func (ec *ExpenseController) GetExpensesByMonth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		year, err := strconv.Atoi(ctx.Query("year"))
		if err != nil {
			ctx.AbortWithError(http.StatusBadRequest, err)
			return
		}
		month, err := strconv.Atoi(ctx.Query("month"))
		if err != nil {
			ctx.AbortWithError(http.StatusBadRequest, err)
			return
		}
		expenses, err := ec.expenseService.GetExpensesByMonth(int32(year), int32(month))
		if err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		ctx.JSON(http.StatusOK, expenses)
	}
}

func (ec *ExpenseController) GetExpenseById() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			ctx.AbortWithError(http.StatusBadRequest, err)
			return
		}
		expense, err := ec.expenseService.GetById(int64(id))
		if err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, err)
		}
		ctx.JSON(http.StatusOK, expense)
	}
}
