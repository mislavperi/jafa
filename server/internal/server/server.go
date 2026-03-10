package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/mislavperi/jafa/server/internal/server/controllers"
	"github.com/mislavperi/jafa/server/internal/server/middleware"
)

type Server struct {
	Gin  *gin.Engine
	port uint
}

func NewServer(expenseController *controllers.ExpenseController, authController *controllers.AuthController, port uint) *Server {
	server := &Server{
		Gin:  gin.Default(),
		port: port,
	}

	server.registerRoutes(expenseController, authController)
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

func sessionSecret() []byte {
	secret := os.Getenv("SESSION_SECRET")
	if secret == "" {
		log.Fatal("SESSION_SECRET environment variable is not set")
	}
	return []byte(secret)
}

func (s *Server) registerRoutes(expenseController *controllers.ExpenseController, authController *controllers.AuthController) {
	store := cookie.NewStore(sessionSecret())
	store.Options(sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: true,
		Secure:   os.Getenv("APP_ENV") == "production",
		SameSite: http.SameSiteLaxMode,
	})
	s.Gin.Use(sessions.Sessions("jafa_session", store))

	authGroup := s.Gin.Group("/auth")
	authGroup.POST("/login", authController.Login())
	authGroup.POST("/logout", authController.Logout())
	authGroup.POST("/register", authController.Register())
	authGroup.GET("/me", authController.Me())

	protected := s.Gin.Group("/")
	protected.Use(middleware.RequireAuth())

	expenseGroup := protected.Group("/expense")
	expenseGroup.GET("/", expenseController.GetAllExpenses())
	expenseGroup.GET("/:id", expenseController.GetExpenseById())

	expenseStatsGroup := protected.Group("/expense-stats")
	expenseStatsGroup.GET("/monthly-total", expenseController.GetTotalSpendThisMonth())
	expenseStatsGroup.GET("/daily-spend", expenseController.GetDailySpend())
	expenseStatsGroup.GET("/first-expense-date", expenseController.GetFirstExpenseDate())
	expenseStatsGroup.GET("/daily-spend-for-month", expenseController.GetDailySpendForMonth())
	expenseStatsGroup.GET("/expenses-by-month", expenseController.GetExpensesByMonth())
}
