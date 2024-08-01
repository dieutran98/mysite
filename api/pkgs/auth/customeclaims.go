package auth

import (
	"mysite/pkgs/env"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type KeyType string

const (
	CursorKey  KeyType = "cursorKey"
	AccessKey  KeyType = "accessKey"
	RefreshKey KeyType = "refreshKey"
)

type CustomClaims[T any] struct {
	jwt.RegisteredClaims
	KeyType  KeyType `json:"key_type"`
	MetaData T       `json:"meta_data,omitempty"`
}

func NewCustomClaims[T any]() *CustomClaims[T] {
	regClaims := defaultRegisterClaims()
	return &CustomClaims[T]{
		RegisteredClaims: regClaims,
	}
}

// summary all condition to validate claims
func (c *CustomClaims[T]) IsValid() bool {
	return c.isValidKey()
}

func (c *CustomClaims[T]) GetKeyType() KeyType {
	return c.KeyType
}

func (c *CustomClaims[T]) isValidKey() bool {
	return c.KeyType == CursorKey || c.KeyType == AccessKey || c.KeyType == RefreshKey
}

func (c *CustomClaims[T]) WithExpireAt(expireAt time.Time) *CustomClaims[T] {
	c.ExpiresAt = jwt.NewNumericDate(expireAt)
	return c
}

func (c *CustomClaims[T]) Clone() *CustomClaims[T] {
	return &CustomClaims[T]{
		RegisteredClaims: c.RegisteredClaims,
		KeyType:          c.KeyType,
		MetaData:         c.MetaData,
	}
}

func defaultRegisterClaims() jwt.RegisteredClaims {
	envObj := env.GetEnv().Jwt
	return jwt.RegisteredClaims{
		Issuer:   envObj.Issuer,
		IssuedAt: jwt.NewNumericDate(time.Now()),
		ID:       uuid.NewString(),
	}
}
