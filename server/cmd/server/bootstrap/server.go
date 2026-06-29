package bootstrap

import (
	"encoding/gob"
	"fmt"
	"os"
	"strconv"

	"github.com/mislavperi/jafa/server/internal/domain/services"
	"github.com/mislavperi/jafa/server/internal/infrastructure/psql"
	server "github.com/mislavperi/jafa/server/internal/server"
	"github.com/mislavperi/jafa/server/internal/server/controllers"
)

func init() {
	gob.Register(int64(0))
}

const defaultPort = 8080

// serverPort returns the HTTP port from the PORT environment variable,
// falling back to the default when unset.
func serverPort() (uint, error) {
	raw := os.Getenv("PORT")
	if raw == "" {
		return defaultPort, nil
	}
	port, err := strconv.ParseUint(raw, 10, 16)
	if err != nil || port == 0 {
		return 0, fmt.Errorf("invalid PORT %q: must be a number between 1 and 65535", raw)
	}
	return uint(port), nil
}

func Server() *server.Server {
	connPool, err := psql.NewDatabaseConnection()
	if err != nil {
		panic(err)
	}

	port, err := serverPort()
	if err != nil {
		panic(err)
	}

	expenseService := services.NewExpenseService(connPool)
	tagService := services.NewTagService(connPool)
	authService := services.NewAuthService(connPool)
	preferencesService := services.NewPreferencesService(connPool)
	categoryService := services.NewCategoryService(connPool)
	reportService := services.NewReportService(connPool)

	expenseController := controllers.NewExpenseController(expenseService)
	tagController := controllers.NewTagController(tagService)
	authController := controllers.NewAuthController(authService)
	preferencesController := controllers.NewPreferencesController(preferencesService)
	categoryController := controllers.NewCategoryController(categoryService)
	reportController := controllers.NewReportController(reportService)

	return server.NewServer(expenseController, tagController, authController, preferencesController, categoryController, reportController, port)
}
