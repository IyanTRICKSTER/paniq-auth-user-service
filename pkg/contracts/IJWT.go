package contracts

type IJWTService interface {
	GenerateToken(userID uint, lifeSpan int, secretKey string) (string, error)
	ValidateToken(token string, secretKey string) bool
	ExtractPayloadFromToken(token string, secretKey string) (map[string]interface{}, error)
}
