package requests

type ChangeUserPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}
