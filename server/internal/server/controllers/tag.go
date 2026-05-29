package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mislavperi/jafa/server/internal/server/httperr"
	"github.com/mislavperi/jafa/server/internal/domain/services"

	requestmodels "github.com/mislavperi/jafa/server/internal/domain/models/request"
)

type TagController struct {
	tagService *services.TagService
}

func NewTagController(tagService *services.TagService) *TagController {
	return &TagController{tagService: tagService}
}

func (tc *TagController) GetAllTags() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		uid, ok := requireUser(ctx)
		if !ok {
			return
		}
		tags, err := tc.tagService.GetAllTags(uid)
		if err != nil {
			httperr.Internal(ctx, err)
			return
		}
		ctx.JSON(http.StatusOK, tags)
	}
}

func (tc *TagController) CreateTag() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		uid, ok := requireUser(ctx)
		if !ok {
			return
		}
		var req requestmodels.CreateTagRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			httperr.BadRequest(ctx, err.Error(), err)
			return
		}
		tag, err := tc.tagService.CreateTag(uid, req.Name)
		if err != nil {
			httperr.Internal(ctx, err)
			return
		}
		ctx.JSON(http.StatusCreated, tag)
	}
}

func (tc *TagController) GetTagsForExpense() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		uid, ok := requireUser(ctx)
		if !ok {
			return
		}
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			httperr.BadRequest(ctx, "invalid expense id", nil)
			return
		}
		tags, err := tc.tagService.GetTagsForExpense(uid, int64(id))
		if err != nil {
			httperr.Internal(ctx, err)
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
		uid, ok := requireUser(ctx)
		if !ok {
			return
		}
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			httperr.BadRequest(ctx, "invalid expense id", nil)
			return
		}
		var req addTagToExpenseRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			httperr.BadRequest(ctx, err.Error(), err)
			return
		}
		if err := tc.tagService.AddTagToExpense(uid, int64(id), req.TagID); err != nil {
			httperr.Internal(ctx, err)
			return
		}
		ctx.Status(http.StatusNoContent)
	}
}

func (tc *TagController) RemoveTagFromExpense() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		uid, ok := requireUser(ctx)
		if !ok {
			return
		}
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			httperr.BadRequest(ctx, "invalid expense id", nil)
			return
		}
		tagID, err := strconv.Atoi(ctx.Param("tag_id"))
		if err != nil {
			httperr.BadRequest(ctx, "invalid tag id", nil)
			return
		}
		if err := tc.tagService.RemoveTagFromExpense(uid, int64(id), int64(tagID)); err != nil {
			httperr.Internal(ctx, err)
			return
		}
		ctx.Status(http.StatusNoContent)
	}
}
