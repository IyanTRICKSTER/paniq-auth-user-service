package statusCodes

type StatusCode int

const (
	Error                    StatusCode = -1
	ModelNotFound            StatusCode = -2
	ErrDuplicatedModel       StatusCode = -3
	ErrExtractJWTToken       StatusCode = -4
	ErrForbidden             StatusCode = -5
	ErrResetPasswordNotMatch StatusCode = -6

	Success StatusCode = 1
)
