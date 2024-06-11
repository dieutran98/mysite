package env

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadEnv(t *testing.T) {
	testCases := []struct {
		name        string
		expectValue error
	}{
		{
			name:        "success",
			expectValue: nil,
		},
		{
			name:        "failed",
			expectValue: errors.New("some error"),
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			require.True(t, errors.Is(ReadEnv(&viperMock{expectErr: tt.expectValue}), tt.expectValue))
		})
	}
}

func TestGetEnv(t *testing.T) {
	internalEnv = appEnv{
		Database: database{
			Database: "test",
		},
	}
	require.Equal(t, "test", GetEnv().Database.Database)

}
