package auth

import "time"

type auth struct{}

type AuthService interface {
	HashPassword(password string) (string, error)
	ComparePasswordAndHash(password, encodedHash string) (match bool, err error)
	NewClaims(userId string, expiredTime time.Time) Claims
	CreateToken(claims Claims, signingKey []byte) (string, error)
	ParseToken(tokenStr string, signingKey []byte) (*Claims, error)
}

func NewAuthService() AuthService {
	return auth{}
}
