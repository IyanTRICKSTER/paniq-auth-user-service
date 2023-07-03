package requests

type RegisterUserRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	NIM      string `json:"nim" binding:""`
	NIP      string `json:"nip" binding:""`
	Major    string `json:"major" binding:"required"`
}

type RegisterUserCSVRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	NIM      string `json:"nim" binding:""`
	NIP      string `json:"nip" binding:""`
	Major    string `json:"major" binding:"required"`
}
