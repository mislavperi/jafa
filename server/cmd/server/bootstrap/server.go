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
	tagService := services.NewTagService(queries)
	authService := services.NewAuthService(queries)

	expenseController := controllers.NewExpenseController(expenseService)
	tagController := controllers.NewTagController(tagService)
	authController := controllers.NewAuthController(authService)

	return server.NewServer(expenseController, tagController, authController, 8080)
}
