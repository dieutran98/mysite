package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashPassword(t *testing.T) {
	svc := auth{}
	{ // hash success
		hash, err := svc.HashPassword("secret")
		assert.NoError(t, err)
		assert.NotEmpty(t, hash)
	}

}

func TestComparePasswordAndHash(t *testing.T) {
	svc := auth{}
	{ // ComparePasswordAndHash success
		hash, err := svc.HashPassword("secret")
		assert.NoError(t, err)
		assert.NotEmpty(t, hash)

		match, err := svc.ComparePasswordAndHash("secret", hash)
		assert.NoError(t, err)
		assert.True(t, match)
	}

	{ // ComparePasswordAndHash failed
		hash, err := svc.HashPassword("secret")
		assert.NoError(t, err)
		assert.NotEmpty(t, hash)

		match, err := svc.ComparePasswordAndHash("secret1", hash)
		assert.NoError(t, err)
		assert.False(t, match)
	}
}
