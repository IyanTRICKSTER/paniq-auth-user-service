package requests

import "errors"

type ResetUserPasswordRequest struct {
	ResetToken string `json:"reset_token" binding:"required"`
	Password   string `json:"password" binding:"required,min=8"`
	CPassword  string `json:"c_password" binding:"required,min=8"`
}

func (credential *ResetUserPasswordRequest) ValidatePassword() (bool, error) {
	if credential.Password != credential.CPassword {
		return false, errors.New("provided passwords are not match")
	}
	return true, nil
}
