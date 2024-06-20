package refresh

import (
	"context"
	"log/slog"
	"mysite/features/refresh/internal"
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
	RefreshToken(ctx context.Context) (*internal.RefreshResponse, error)
}

var newService = func(req internal.RefreshRequest) service {
	return internal.NewService(req)
}

func NewHandler() *api {
	return &api{}
}

func (a api) Refresh(w http.ResponseWriter, r *http.Request) {
	var body model.RefreshJSONRequestBody
	if err := httputil.ParseBody(r, &body); err != nil {
		if err := render.Render(w, r, httputil.NewFailureRender(errors.Wrap(err, "failed to parse body"))); err != nil {
			slog.Error("failed to render", logger.AttrError(err))
		}
		return
	}

	// refresh business
	params, err := internal.NewParams(body)
	if err != nil {
		if err := render.Render(w, r, httputil.NewFailureRender(errors.Wrap(err, "failed to parse params"))); err != nil {
			slog.Error("failed to render", logger.AttrError(err))
		}
		return
	}

	resp, err := newService(*params).RefreshToken(r.Context())
	if err != nil {
		if err := render.Render(w, r, httputil.NewFailureRender(errors.Wrap(err, "failed to refresh token"))); err != nil {
			slog.Error("failed to render", logger.AttrError(err))
		}
		return
	}

	http.SetCookie(w, httputil.SetCookie("accessToken", resp.AccessToken))
}
