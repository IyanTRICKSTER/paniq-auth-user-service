package bcryptUtils

import (
	"golang.org/x/crypto/bcrypt"
	"paniq-auth-user-service/pkg/contracts"
)

type HashFunction struct{}

func (h HashFunction) Hash(payload string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(payload), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	return string(hash)
}

func (h HashFunction) HashCheck(hashed string, payload string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(payload))
	if err != nil {
		return false, err
	}
	return true, nil
}

func NewHashFunction() contracts.IHash {
	return &HashFunction{}
}
