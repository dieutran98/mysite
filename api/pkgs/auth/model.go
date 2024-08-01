package auth

type auth struct{}

//go:generate moq -pkg pkgmock -out ../../testing/mocking/pkgmock/authservice.mock.go . AuthService
type AuthService interface {
	HashPassword(password string) (string, error)
	ComparePasswordAndHash(password, encodedHash string) (match bool, err error)
}

func NewAuthService() AuthService {
	return auth{}
}
