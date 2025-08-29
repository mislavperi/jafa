package bootstrap

import (
	"github.com/mislavperi/jafa/internal/api"
	"github.com/mislavperi/jafa/internal/api/controllers"
	"github.com/mislavperi/jafa/internal/domain/services"
	"github.com/mislavperi/jafa/internal/infrastructure/psql"
	psqlrepositories "github.com/mislavperi/jafa/internal/infrastructure/psql/repositories"
)

func Api() *api.Api {
	connPool, err := psql.NewDatabaseConnection()
	if err != nil {
		panic(err)
	}

	queries := psqlrepositories.New(connPool)
	expenseService := services.NewExpenseService(queries)

	expenseController := controllers.NewExpenseController(expenseService)

	api := api.NewApi(expenseController, 8080)
	return api
}
