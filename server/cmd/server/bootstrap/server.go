package bootstrap

import (
	"github.com/mislavperi/jafa/server/internal/domain/services"
	"github.com/mislavperi/jafa/server/internal/infrastructure/psql"
	psqlrepositories "github.com/mislavperi/jafa/server/internal/infrastructure/psql/repositories"
	server "github.com/mislavperi/jafa/server/internal/server"
	"github.com/mislavperi/jafa/server/internal/server/controllers"
)

func Server() *server.Server {
	connPool, err := psql.NewDatabaseConnection()
	if err != nil {
		panic(err)
	}

	queries := psqlrepositories.New(connPool)

	expenseService := services.NewExpenseService(queries)
	tagService := services.NewTagService(queries)

	expenseController := controllers.NewExpenseController(expenseService)
	tagController := controllers.NewTagController(tagService)

	return server.NewServer(expenseController, tagController, 8080)
}
