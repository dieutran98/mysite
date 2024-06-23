package authmock

import (
	"mysite/pkgs/auth"
	"time"
)

type authMock struct {
	auth.AuthService
	HashPasswordFunc           func() (string, error)
	ComparePasswordAndHashFunc func() (match bool, err error)
	NewClaimsFunc              func() auth.Claims
	CreateTokenFunc            func() (string, error)
	ParseTokenFunc             func() (*auth.Claims, error)
}

func NewMockService() authMock {
	return authMock{}
}

func (m authMock) HashPassword(password string) (string, error) {
	return m.HashPasswordFunc()
}

func (m authMock) ComparePasswordAndHash(password, encodedHash string) (match bool, err error) {
	return m.ComparePasswordAndHashFunc()
}

func (m authMock) NewClaims(userId int, expiredTime time.Time) auth.Claims {
	return m.NewClaimsFunc()
}

func (m authMock) CreateToken(claims auth.Claims, signingKey []byte) (string, error) {
	return m.CreateTokenFunc()
}

func (m authMock) ParseToken(tokenStr string, signingKey []byte) (*auth.Claims, error) {
	return m.ParseTokenFunc()
}
