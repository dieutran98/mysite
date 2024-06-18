package register

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/render"
	"github.com/pkg/errors"

	"mysite/features/register/internal"
	"mysite/models/model"
	"mysite/utils/httputil"
)

type api struct{}

type service interface {
	Register(ctx context.Context) error
}

var newService = func(req internal.RegisterRequest) service {
	return internal.NewService(req)
}

func NewHandler() *api {
	return &api{}
}

func (a *api) Register(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil || r.Body == http.NoBody {
		render.Render(w, r, httputil.NewFailureRender(errors.Wrap(httputil.ErrInvalidRequest, "empty body")))
		return
	}

	var body model.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		render.Render(w, r, httputil.NewFailureRender(errors.Wrap(err, "failed to decode body")))
		return
	}

	params, err := internal.NewParams(body)
	if err != nil {
		render.Render(w, r, httputil.NewFailureRender(errors.Wrap(err, "failed to mapping body to params")))
		return
	}

	if err := newService(*params).Register(r.Context()); err != nil {
		render.Render(w, r, httputil.NewFailureRender(errors.Wrap(err, "failed register user")))
		return
	}
	w.WriteHeader(http.StatusCreated)
}
