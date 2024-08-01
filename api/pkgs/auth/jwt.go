package auth

import (
	"mysite/pkgs/env"

	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
)

type jwtHandler struct {
	claims Claims
}

type Claims interface {
	jwt.Claims
	IsValid() bool
	GetKeyType() KeyType
}

//go:generate moq -pkg pkgmock -out ../../testing/mocking/pkgmock/jwt.mock.go . JwtHandler
type JwtHandler interface {
	CreateToken() (string, error)
	ParseToken(tokenString string, claims Claims) error
	WithClaims(claims Claims) JwtHandler
}

func NewJwtHandler() JwtHandler {
	return &jwtHandler{}
}

func (j jwtHandler) WithClaims(claims Claims) JwtHandler {
	j.claims = claims
	return &j
}

func (j *jwtHandler) CreateToken() (string, error) {
	// validate claims
	if !j.claims.IsValid() {
		return "", errors.New("invalid claims")
	}

	// get jwt key
	key, err := getJwtKey(j.claims.GetKeyType())
	if err != nil {
		return "", errors.Wrap(err, "failed get jwt key")
	}

	return jwt.NewWithClaims(jwt.SigningMethodHS256, j.claims).SignedString(key)
}

func (j *jwtHandler) ParseToken(tokenString string, claims Claims) error {
	token, err := jwt.ParseWithClaims(tokenString, claims, keyParser)
	if err != nil {
		return errors.Wrap(err, "failed parse token")
	}

	if !token.Valid {
		return errors.New("invalid token")
	}
	return nil
}

func keyParser(token *jwt.Token) (interface{}, error) {
	method, ok := token.Method.(*jwt.SigningMethodHMAC)
	if !ok {
		return nil, errors.Errorf("Unexpected signing method: %v", token.Header["alg"])
	}
	if method.Alg() != jwt.SigningMethodHS256.Name {
		return nil, errors.New("the alg must be HS256")
	}

	claims, ok := token.Claims.(Claims)
	if !ok {
		return nil, errors.New("invalid claims")
	}

	return getJwtKey(claims.GetKeyType())
}

func getJwtKey(keyType KeyType) ([]byte, error) {
	envObj := env.GetEnv().Jwt
	switch keyType {
	case CursorKey:
		return []byte(envObj.CursorKey), nil
	case AccessKey:
		return []byte(envObj.AccessKey), nil
	case RefreshKey:
		return []byte(envObj.RefreshKey), nil
	default:
		return nil, errors.New("unsupported key type")
	}
}
