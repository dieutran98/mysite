package httputil

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert" // Using testify for better assertions

	"mysite/utils/ptrconv"
)

var (
	errSample = errors.New("sample error")
)

func TestNewFailureRender(t *testing.T) {
	{ // internal error
		renderer := NewFailureRender(errSample)
		assert.IsType(t, ErrResponse{}, renderer, "Should return existing ErrInternal")

		// Assert error text is set correctly
		assert.Equal(t, ptrconv.String(errSample.Error()), renderer.(ErrResponse).ErrorText, "Error text should be set")
		assert.Equal(t, ErrInternal.HTTPStatusCode, renderer.(ErrResponse).HTTPStatusCode, "Status code should be InternalServerError")
		assert.Equal(t, ErrInternal.StatusText, renderer.(ErrResponse).StatusText, "Status text should be Internal error")
	}

	{ // notfound error
		renderer := NewFailureRender(errors.Wrap(ErrNotFound, "not found"))
		assert.IsType(t, ErrResponse{}, renderer, "Should return existing ErrInternal")

		// Assert error text is set correctly
		assert.Equal(t, "not found: ", *renderer.(ErrResponse).ErrorText, "Error text should be set")
		assert.Equal(t, ErrNotFound.HTTPStatusCode, renderer.(ErrResponse).HTTPStatusCode, "Status code should be InternalServerError")
		assert.Equal(t, ErrNotFound.StatusText, renderer.(ErrResponse).StatusText, "Status text should be Internal error")
	}

	{ // invalid request error
		renderer := NewFailureRender(errors.Wrap(ErrInvalidRequest, "invalid request"))
		assert.IsType(t, ErrResponse{}, renderer, "Should return existing ErrInternal")

		// Assert error text is set correctly
		assert.Equal(t, "invalid request: ", *renderer.(ErrResponse).ErrorText, "Error text should be set")
		assert.Equal(t, ErrInvalidRequest.HTTPStatusCode, renderer.(ErrResponse).HTTPStatusCode, "Status code should be InternalServerError")
		assert.Equal(t, ErrInvalidRequest.StatusText, renderer.(ErrResponse).StatusText, "Status text should be Internal error")
	}

}
