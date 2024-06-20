package auth

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestJWT(t *testing.T) {
	svc := auth{}

	{ // success create token and parse it
		userId := "user-id"
		signKey := "secret-key"
		tokenStr, err := svc.CreateToken(svc.NewClaims(userId, time.Now().Add(15*time.Minute)), []byte(signKey))
		require.NoError(t, err)
		require.NotEmpty(t, tokenStr)

		claims, err := svc.ParseToken(tokenStr, []byte(signKey))
		require.NoError(t, err)
		require.NotNil(t, claims)
		require.Equal(t, userId, claims.Subject)
	}
}
