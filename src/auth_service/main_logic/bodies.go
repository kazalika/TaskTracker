package auth_service

type AuthenticateBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type ProfileInfo struct {
	FirstName   string `json:"firstName,omitempty"`
	LastName    string `json:"lastName,omitempty"`
	Birthday    string `json:"birthday,omitempty"`
	Email       string `json:"email,omitempty"`
	PhoneNumber string `json:"phoneNumber,omitempty"`
}

type RegisterBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type TaskListRequest struct {
	Offset   int32 `json:"offset"`
	PageSize int32 `json:"page_size"`
}

type CreateTaskRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

type UpdateTaskRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

type LikeRequest struct {
	TaskID int32 `json:"post_id"`
}

type ViewRequest struct {
	TaskID int32 `json:"post_id"`
}

type TaskID struct {
	TaskID int32 `json:"task_id"`
}

type TaskContent struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
}
