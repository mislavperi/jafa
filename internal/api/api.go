package api

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
)

type API struct {
	Gin  *gin.Engine
	port uint
}

func NewAPI() *API {
	return &API{
		Gin: gin.Default(),
		port: 8080,
	}
}

func (a *API) Start(ctx context.Context) {
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
