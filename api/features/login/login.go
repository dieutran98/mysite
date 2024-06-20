package login

import (
	"context"
	"log/slog"
	"mysite/features/login/internal"
	"mysite/models/model"
	"mysite/pkgs/logger"
	"mysite/utils/httputil"
	"net/http"

	"github.com/go-chi/render"
	"github.com/pkg/errors"
)

type api struct {
}

type service interface {
	Login(ctx context.Context) (*internal.LoginResponse, error)
}

var newService = func(req internal.LoginRequest) service {
	return internal.NewService(req)
}

func NewHandler() *api {
	return &api{}
}

func (a api) Login(w http.ResponseWriter, r *http.Request) {
	var body model.LoginJSONRequestBody

	if err := httputil.ParseBody(r, &body); err != nil {
		if err := render.Render(w, r, httputil.NewFailureRender(errors.Wrap(err, "failed to parse body"))); err != nil {
			slog.Error("failed to render", logger.AttrError(err))
		}
		return
	}

	// login business logic
	params, err := internal.NewParams(body)
	if err != nil {
		if err := render.Render(w, r, httputil.NewFailureRender(errors.Wrap(err, "failed to create params"))); err != nil {
			slog.Error("failed to render", logger.AttrError(err))
		}
		return
	}

	resp, err := newService(*params).Login(r.Context())
	if err != nil {
		if err := render.Render(w, r, httputil.NewFailureRender(errors.Wrap(err, "failed to create body"))); err != nil {
			slog.Error("failed to render", logger.AttrError(err))
		}
		return
	}

	http.SetCookie(w, httputil.SetCookie("accessToken", resp.AccessToken))
	http.SetCookie(w, httputil.SetCookie("refreshToken", resp.RefreshToken))
}
