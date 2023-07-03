package bcryptUtils

import (
	"github.com/stretchr/testify/mock"
)

type HashFunctionMock struct {
	mock.Mock
}

func (h *HashFunctionMock) Hash(payload string) string {
	return h.Mock.Called(payload).Get(0).(string)
}

func (h *HashFunctionMock) HashCheck(hashed string, payload string) (bool, error) {
	args := h.Mock.Called(hashed, payload)
	return args.Get(0).(bool), args.Get(1).(error)
}
