package httputil

import (
	"mysite/models/model"
	"mysite/utils/ptrconv"
	"net/http"
)

var (
	ErrNotFound = ErrResponse{
		HTTPStatusCode: http.StatusNotFound,
		ErrorResponse: model.ErrorResponse{
			StatusText: ptrconv.String("Resource not found"),
		},
	}

	ErrInvalidRequest = ErrResponse{
		HTTPStatusCode: http.StatusBadRequest,
		ErrorResponse: model.ErrorResponse{
			StatusText: ptrconv.String("Invalid request"),
		},
	}

	ErrInternal = ErrResponse{
		HTTPStatusCode: http.StatusInternalServerError,
		ErrorResponse: model.ErrorResponse{
			StatusText: ptrconv.String("Internal error"),
		},
	}
)
