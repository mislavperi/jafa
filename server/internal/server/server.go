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

func NewServer(expenseController *controllers.ExpenseController, port uint) *Server {
	server := &Server{
		Gin:  gin.Default(),
		port: port,
	}

	server.registerRoutes(expenseController)
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

func (s *Server) registerRoutes(expenseController *controllers.ExpenseController) {
	expenseGroup := s.Gin.Group("/expense")
	expenseGroup.GET("/:id", expenseController.GetExpenseById())
}
