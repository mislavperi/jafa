package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mislavperi/jafa/server/internal/server/httperr"
)

type CategoryController struct {
	categoryService CategoryService
}

func NewCategoryController(categoryService CategoryService) *CategoryController {
	return &CategoryController{categoryService: categoryService}
}

func (cc *CategoryController) List() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		categories, err := cc.categoryService.List(ctx.Request.Context())
		if err != nil {
			httperr.Internal(ctx, err)
			return
		}
		ctx.JSON(http.StatusOK, categories)
	}
}
