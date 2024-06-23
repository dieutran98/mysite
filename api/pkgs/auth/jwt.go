package auth

import (
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
)

type Claims struct {
	jwt.RegisteredClaims
}

func (auth) NewClaims(userId int, expiredTime time.Time) Claims {
	return Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiredTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Subject:   fmt.Sprintf("%d", userId),
		},
	}
}

func (auth) CreateToken(claims Claims, signingKey []byte) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedStr, err := token.SignedString(signingKey)
	return signedStr, err
}

func (auth) ParseToken(tokenStr string, signingKey []byte) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		return signingKey, nil
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse token")
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, errors.New("invalid claims")
	}

	return claims, nil
}

func (c Claims) GetUserId() (int, error) {
	id, err := strconv.Atoi(c.Subject)
	if err != nil {
		return 0, errors.Wrap(err, "failed to convert string to int")
	}
	return id, nil
}
