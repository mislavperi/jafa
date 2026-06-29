package request

type CreateTagRequest struct {
	Name string `json:"name" binding:"required"`
}

type AddTagToExpenseRequest struct {
	TagID int64 `json:"tag_id" binding:"required"`
}
