package env

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

type configureMock struct {
	expectErr error
}

func (c configureMock) setConfigFile() error {
	return c.expectErr
}

func (c configureMock) mappingStruct() error {
	return c.expectErr
}

func (c configureMock) setDefault() error {
	return c.expectErr
}

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
			newConfigure = func() configure {
				return configureMock{
					expectErr: tt.expectValue,
				}
			}

			require.True(t, errors.Is(ReadEnv(), tt.expectValue))
		})
	}
}

func TestGetEnv(t *testing.T) {
	internalEnv = AppEnv{
		Database: database{
			Database: "test",
		},
	}
	require.Equal(t, "test", GetEnv().Database.Database)

}
