package request

type CreateTagRequest struct {
	Name string `json:"name" binding:"required"`
}
