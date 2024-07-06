package httputil

import (
	"mysite/dtos"
	"mysite/utils/ptrconv"
	"net/http"
)

var (
	ErrNotFound = ErrResponse{
		HTTPStatusCode: http.StatusNotFound,
		ErrorResponse: dtos.ErrorResponse{
			StatusText: ptrconv.String("Resource not found"),
		},
	}

	ErrInvalidRequest = ErrResponse{
		HTTPStatusCode: http.StatusBadRequest,
		ErrorResponse: dtos.ErrorResponse{
			StatusText: ptrconv.String("Invalid request"),
		},
	}

	ErrInternal = ErrResponse{
		HTTPStatusCode: http.StatusInternalServerError,
		ErrorResponse: dtos.ErrorResponse{
			StatusText: ptrconv.String("Internal error"),
		},
	}

	ErrUnauthorize = ErrResponse{
		HTTPStatusCode: http.StatusUnauthorized,
		ErrorResponse: dtos.ErrorResponse{
			StatusText: ptrconv.String("Unauthorize error"),
		},
	}
)
