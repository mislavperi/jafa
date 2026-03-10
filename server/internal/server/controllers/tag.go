package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mislavperi/jafa/server/internal/domain/services"
)

type TagController struct {
	tagService *services.TagService
}

func NewTagController(tagService *services.TagService) *TagController {
	return &TagController{tagService: tagService}
}

func (tc *TagController) GetAllTags() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tags, err := tc.tagService.GetAllTags()
		if err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		ctx.JSON(http.StatusOK, tags)
	}
}

type createTagRequest struct {
	Name string `json:"name" binding:"required"`
}

func (tc *TagController) CreateTag() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req createTagRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		tag, err := tc.tagService.CreateTag(req.Name)
		if err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		ctx.JSON(http.StatusCreated, tag)
	}
}

func (tc *TagController) GetTagsForExpense() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid expense id"})
			return
		}
		tags, err := tc.tagService.GetTagsForExpense(int64(id))
		if err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		ctx.JSON(http.StatusOK, tags)
	}
}

type addTagToExpenseRequest struct {
	TagID int64 `json:"tag_id" binding:"required"`
}

func (tc *TagController) AddTagToExpense() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid expense id"})
			return
		}
		var req addTagToExpenseRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := tc.tagService.AddTagToExpense(int64(id), req.TagID); err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		ctx.Status(http.StatusNoContent)
	}
}

func (tc *TagController) RemoveTagFromExpense() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid expense id"})
			return
		}
		tagID, err := strconv.Atoi(ctx.Param("tag_id"))
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid tag id"})
			return
		}
		if err := tc.tagService.RemoveTagFromExpense(int64(id), int64(tagID)); err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		ctx.Status(http.StatusNoContent)
	}
}
