package api

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/mislavperi/jafa/internal/api/controllers"
)

type Api struct {
	Gin  *gin.Engine
	port uint
}

func NewApi(expenseController *controllers.ExpenseController, port uint) *Api {
	api := &Api{
		Gin:  gin.Default(),
		port: port,
	}

	api.registerRoutes(expenseController)
	return api
}

func (a *Api) Start(ctx context.Context) {
	errs := make(chan error, 1)

	go func() {
		errs <- a.Gin.Run(fmt.Sprintf(":%d", a.port))
	}()

	select {
	case <-errs:
		return
	case <-ctx.Done():
		return
	}
}

func (a *Api) registerRoutes(expenseController *controllers.ExpenseController) {
	expenseGroup := a.Gin.Group("/api/expense")
	expenseGroup.GET("/:id", expenseController.GetExpenseById())
}
