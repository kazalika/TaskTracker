package auth_service

type TaskListRequest struct {
	Offset   int32 `json:"offset"`
	PageSize int32 `json:"page_size"`
}

type CreateTaskBody struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

type UpdateTaskBody struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

type LikeRequest struct {
	TaskID string `json:"post_id"`
}

type ViewRequest struct {
	TaskID string `json:"post_id"`
}
