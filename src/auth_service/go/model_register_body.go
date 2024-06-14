package auth_service

type RegisterBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
