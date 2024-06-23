package auth

import (
	"fmt"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/require"
)

func TestJWT(t *testing.T) {
	svc := auth{}

	{ // success create token and parse it
		userId := 1
		signKey := "secret-key"
		tokenStr, err := svc.CreateToken(svc.NewClaims(userId, time.Now().Add(15*time.Minute)), []byte(signKey))
		require.NoError(t, err)
		require.NotEmpty(t, tokenStr)

		claims, err := svc.ParseToken(tokenStr, []byte(signKey))
		require.NoError(t, err)
		require.NotNil(t, claims)
		require.Equal(t, fmt.Sprintf("%d", userId), claims.Subject)
	}
}

func TestGetUserId(t *testing.T) {
	claim := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject: "1",
		},
	}

	id, err := claim.GetUserId()
	require.NoError(t, err)
	require.Equal(t, 1, id)
}
