package api

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/mislavperi/jafa/server/internal/server/controllers"
)

type Server struct {
	Gin  *gin.Engine
	port uint
}

func NewServer(expenseController *controllers.ExpenseController, tagController *controllers.TagController, port uint) *Server {
	server := &Server{
		Gin:  gin.Default(),
		port: port,
	}

	server.registerRoutes(expenseController, tagController)
	return server
}

func (s *Server) Start(ctx context.Context) {
	errs := make(chan error, 1)

	go func() {
		errs <- s.Gin.Run(fmt.Sprintf(":%d", s.port))
	}()

	select {
	case <-errs:
		return
	case <-ctx.Done():
		return
	}
}

func (s *Server) registerRoutes(expenseController *controllers.ExpenseController, tagController *controllers.TagController) {
	expenseGroup := s.Gin.Group("/expense")
	expenseGroup.GET("/", expenseController.GetAllExpenses())
	expenseGroup.POST("/", expenseController.CreateExpense())
	expenseGroup.GET("/:id", expenseController.GetExpenseById())
	expenseGroup.GET("/:id/tags", tagController.GetTagsForExpense())
	expenseGroup.POST("/:id/tags", tagController.AddTagToExpense())
	expenseGroup.DELETE("/:id/tags/:tag_id", tagController.RemoveTagFromExpense())

	expenseStatsGroup := s.Gin.Group("/expense-stats")
	expenseStatsGroup.GET("/monthly-total", expenseController.GetTotalSpendThisMonth())
	expenseStatsGroup.GET("/daily-spend", expenseController.GetDailySpend())
	expenseStatsGroup.GET("/first-expense-date", expenseController.GetFirstExpenseDate())
	expenseStatsGroup.GET("/daily-spend-for-month", expenseController.GetDailySpendForMonth())
	expenseStatsGroup.GET("/expenses-by-month", expenseController.GetExpensesByMonth())

	tagGroup := s.Gin.Group("/tags")
	tagGroup.GET("/", tagController.GetAllTags())
	tagGroup.POST("/", tagController.CreateTag())
}
