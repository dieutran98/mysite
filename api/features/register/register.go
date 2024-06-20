package register

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/go-chi/render"
	"github.com/pkg/errors"

	"mysite/features/register/internal"
	"mysite/models/model"
	"mysite/pkgs/logger"
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
	var body model.RegisterRequest
	if err := httputil.ParseBody(r, &body); err != nil {
		if err := render.Render(w, r, httputil.NewFailureRender(errors.Wrap(err, "failed to parse body"))); err != nil {
			slog.Error("failed to render", logger.AttrError(err))
		}
		return
	}

	params, err := internal.NewParams(body)
	if err != nil {
		if err := render.Render(w, r, httputil.NewFailureRender(errors.Wrap(err, "failed to mapping body to params"))); err != nil {
			slog.Error("failed to render", logger.AttrError(err))
		}
		return
	}

	if err := newService(*params).Register(r.Context()); err != nil {
		if err := render.Render(w, r, httputil.NewFailureRender(errors.Wrap(err, "failed register user"))); err != nil {
			slog.Error("failed to render", logger.AttrError(err))
		}
		return
	}
	w.WriteHeader(http.StatusCreated)
}
