package requests

type UpdateUserRequest struct {
	//Username string `json:"username" binding:"required"`
	//Email    string `json:"email" binding:"required,email"`
	Avatar string `json:"avatar" binding:"required"`
	NIM    string `json:"nim" binding:""`
	NIP    string `json:"nip" binding:""`
	Major  string `json:"major" binding:"required"`
}
