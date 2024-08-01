package auth

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type myClaims struct {
	Foo string `json:"foo"`
}

func TestJwt(t *testing.T) {
	metaData := myClaims{Foo: "bar"}
	claimsA := NewCustomClaims[myClaims]().WithExpireAt(time.Now().Add(time.Hour))
	claimsA.MetaData = metaData
	claimsA.KeyType = CursorKey

	tokenStr, err := NewJwtHandler().WithClaims(claimsA).CreateToken()
	require.NoError(t, err)

	fmt.Println(string(tokenStr))

	var claimsB CustomClaims[myClaims]
	err = NewJwtHandler().ParseToken(tokenStr, &claimsB)
	require.NoError(t, err)
	require.Equal(t, claimsA, &claimsB)

}
