package auth_service

type AuthenticateBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
