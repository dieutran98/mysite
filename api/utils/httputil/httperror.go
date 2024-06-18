package httputil

import (
	"mysite/models/model"
	"mysite/utils/ptrconv"
	"net/http"

	"github.com/go-chi/render"
	"github.com/pkg/errors"
)

type ErrResponse struct {
	Err            error `json:"-"` // low-level runtime error
	HTTPStatusCode int   `json:"-"` // http response status code

	model.ErrorResponse
}

func (e ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

func (e ErrResponse) Error() string {
	return ptrconv.SafeString(e.ErrorText)
}

func (e ErrResponse) WithErrorText(err error) render.Renderer {
	e.ErrorText = ptrconv.String(err.Error())
	e.Err = err
	return e
}

func NewFailureRender(err error) render.Renderer {
	errRender := ErrInternal
	if errors.As(err, &errRender) {
		errRender.ErrorText = ptrconv.String(err.Error())
		return errRender
	}
	return errRender.WithErrorText(err)
}
