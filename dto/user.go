package dto

type UserRequest struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserResponse struct {
	Email    string `json:"email"`
	Username string `json:"username"`
}
