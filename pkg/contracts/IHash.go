package contracts

type IHash interface {
	Hash(payload string) string
	HashCheck(hashed string, payload string) (bool, error)
}
