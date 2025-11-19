package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mislavperi/jafa/server/internal/domain/services"
	"github.com/mislavperi/jafa/web/templates"
)

type ExpenseController struct {
	expenseService *services.ExpenseService
}

func NewExpenseController(expenseService *services.ExpenseService) *ExpenseController {
	return &ExpenseController{
		expenseService: expenseService,
	}
}

func (ec *ExpenseController) GetExpenseById() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			ctx.AbortWithError(http.StatusBadRequest, err)
			return
		}
		_, err = ec.expenseService.GetById(int64(id))
		if err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, err)
		}
		ctx.HTML(http.StatusOK, "", templates.Hello("my name"))
	}
}
