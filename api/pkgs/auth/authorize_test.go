package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashPassword(t *testing.T) {
	{ // hash success
		hash, err := HashPassword("secret")
		assert.NoError(t, err)
		assert.NotEmpty(t, hash)
	}

}

func TestComparePasswordAndHash(t *testing.T) {
	{ // ComparePasswordAndHash success
		hash, err := HashPassword("secret")
		assert.NoError(t, err)
		assert.NotEmpty(t, hash)

		match, err := ComparePasswordAndHash("secret", hash)
		assert.NoError(t, err)
		assert.True(t, match)
	}

	{ // ComparePasswordAndHash failed
		hash, err := HashPassword("secret")
		assert.NoError(t, err)
		assert.NotEmpty(t, hash)

		match, err := ComparePasswordAndHash("secret1", hash)
		assert.NoError(t, err)
		assert.False(t, match)
	}
}
