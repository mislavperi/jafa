package bootstrap

import (
	"github.com/mislavperi/jafa/server/internal/domain/services"
	"github.com/mislavperi/jafa/server/internal/infrastructure/gintemplrenderer"
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

	expenseController := controllers.NewExpenseController(expenseService)

	server := server.NewServer(expenseController, 8080)

	ginHtmlRenderer := server.Gin.HTMLRender
	server.Gin.HTMLRender = &gintemplrenderer.HTMLTemplRenderer{FallbackHtmlRenderer: ginHtmlRenderer}
	server.Gin.SetTrustedProxies(nil)

	return server
}
