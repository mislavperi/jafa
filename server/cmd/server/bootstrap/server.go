package bootstrap

import (
	"encoding/gob"

	"github.com/mislavperi/jafa/server/internal/domain/services"
	"github.com/mislavperi/jafa/server/internal/infrastructure/psql"
	psqlrepositories "github.com/mislavperi/jafa/server/internal/infrastructure/psql/repositories"
	server "github.com/mislavperi/jafa/server/internal/server"
	"github.com/mislavperi/jafa/server/internal/server/controllers"
)

func init() {
	gob.Register(int64(0))
}

func Server() *server.Server {
	connPool, err := psql.NewDatabaseConnection()
	if err != nil {
		panic(err)
	}

	queries := psqlrepositories.New(connPool)

	expenseService := services.NewExpenseService(queries)
	authService := services.NewAuthService(queries)

	expenseController := controllers.NewExpenseController(expenseService)
	authController := controllers.NewAuthController(authService)

	server := server.NewServer(expenseController, authController, 8080)
	return server
}
