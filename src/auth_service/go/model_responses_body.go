package auth_service

type TaskID struct {
	TaskID string `json:"task_id"`
}

type TaskContent struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
}
