package validate

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type testStruct struct {
	FieldA string `validate:"required"`
}

func TestValidateStruct(t *testing.T) {

	{ // failed validate
		testData := testStruct{
			FieldA: "",
		}
		require.Error(t, ValidateStruct(testData))
	}
	{ // success validate
		testData := testStruct{
			FieldA: "testing",
		}
		require.NoError(t, ValidateStruct(testData))
	}

}
