package http

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/pkg/errors"
)

type ErrResponse struct {
	Err            error `json:"-"` // low-level runtime error
	HTTPStatusCode int   `json:"-"` // http response status code

	StatusText string `json:"status"`          // user-level status message
	AppCode    int64  `json:"code,omitempty"`  // application-specific error code
	ErrorText  string `json:"error,omitempty"` // application-level error message, for debugging
}

func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

func (e *ErrResponse) String() string {
	return e.ErrorText
}

func (e ErrResponse) WithErrorText(err error) render.Renderer {
	e.ErrorText = err.Error()
	return &e
}

func FailureRender(err error) render.Renderer {
	var errRender render.Renderer
	if errors.As(err, errRender) {
		return errRender
	}
	return ErrInvalidRequest.WithErrorText(err)
}
