package http

var (
	ErrNotFound = &ErrResponse{
		HTTPStatusCode: 404,
		StatusText:     "Resource not found",
	}

	ErrInvalidRequest = &ErrResponse{
		HTTPStatusCode: 400,
		StatusText:     "Invalid request",
	}

	ErrInternal = &ErrResponse{
		HTTPStatusCode: 500,
		StatusText:     "Internal error",
	}
)
