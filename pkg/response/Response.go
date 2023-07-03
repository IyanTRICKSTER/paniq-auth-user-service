package response

import (
	"paniq-auth-user-service/pkg/contracts/statusCodes"
)

type Response struct {
	status     bool
	statusCode statusCodes.StatusCode
	message    string
	data       any
}

func (r Response) GetStatus() bool {
	return r.status
}

func (r Response) GetStatusCode() statusCodes.StatusCode {
	return r.statusCode
}

func (r Response) GetMessage() string {
	return r.message
}

func (r Response) GetData() any {
	return r.data
}

func (r *Response) SetMessage(msg string) {
	r.message = msg
}

func (r Response) ErrorIs(code statusCodes.StatusCode) bool {
	if r.statusCode == code {
		return true
	}
	return false
}

func (r Response) IsFailed() bool {
	if r.status == true {
		return false
	}
	return true
}

func (r Response) ToMapStringInterface() map[string]interface{} {
	return map[string]interface{}{
		"status":      r.status,
		"status_code": r.statusCode,
		"message":     r.message,
		"data":        r.data,
	}
}

func New[Res UserResponse | AuthResponse](res Res, status bool, statusCode statusCodes.StatusCode, message string, data any) Res {
	return Res{Response{status: status, statusCode: statusCode, message: message, data: data}}
}
