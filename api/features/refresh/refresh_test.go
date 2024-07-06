package refresh

import (
	"bytes"
	"context"
	"encoding/json"
	"mysite/dtos"
	"mysite/features/refresh/internal"
	"mysite/utils/httputil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"slices"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockService struct {
	RefreshTokenFunc func() (*internal.RefreshResponse, error)
}

func (m mockService) RefreshToken(ctx context.Context) (*internal.RefreshResponse, error) {
	return m.RefreshTokenFunc()
}

func newTestRouter() *chi.Mux {
	router := chi.NewRouter()
	router.Route("/api/v1", func(subr chi.Router) {
		HandlerFromMux(NewHandler(), subr)
	})
	return router
}

func TestDashboardGetStores(t *testing.T) {
	t.Parallel()
	router := newTestRouter()

	tests := []struct {
		name       string
		req        func(context.Context) (*http.Request, error)
		assert     func(*httptest.ResponseRecorder, *http.Request)
		newService func(req internal.RefreshRequest) service
	}{
		{
			name: "404",
			req: func(ctx context.Context) (*http.Request, error) {
				return http.NewRequest(http.MethodPost, "http://example.com", nil)
			},
			assert: func(w *httptest.ResponseRecorder, r *http.Request) {
				assert.Equal(t, http.StatusNotFound, w.Result().StatusCode)
			},
		},
		{
			name: "400 - Invalid request",
			req: func(ctx context.Context) (*http.Request, error) {
				_url, err := url.Parse("http://example.com/api/v1/refresh")
				require.NoError(t, err)
				var buf bytes.Buffer
				if err := json.NewEncoder(&buf).Encode(dtos.RefreshRequest{RefreshToken: "token"}); err != nil {
					return nil, errors.Wrap(err, "failed encode body")
				}

				return http.NewRequest(http.MethodPost, _url.String(), &buf)
			},
			assert: func(w *httptest.ResponseRecorder, r *http.Request) {
				assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
			},
			newService: func(req internal.RefreshRequest) service {
				return mockService{RefreshTokenFunc: func() (*internal.RefreshResponse, error) { return nil, httputil.ErrInvalidRequest }}
			},
		},
		{
			name: "400 - Invalid request, empty body",
			req: func(ctx context.Context) (*http.Request, error) {
				_url, err := url.Parse("http://example.com/api/v1/refresh")
				require.NoError(t, err)

				return http.NewRequest(http.MethodPost, _url.String(), nil)
			},
			assert: func(w *httptest.ResponseRecorder, r *http.Request) {
				assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
			},
			newService: func(req internal.RefreshRequest) service {
				return mockService{RefreshTokenFunc: func() (*internal.RefreshResponse, error) { return nil, nil }}
			},
		},
		{
			name: "500 - internal error request",
			req: func(ctx context.Context) (*http.Request, error) {
				_url, err := url.Parse("http://example.com/api/v1/refresh")
				require.NoError(t, err)
				var buf bytes.Buffer
				if err := json.NewEncoder(&buf).Encode(dtos.RefreshRequest{RefreshToken: "token"}); err != nil {
					return nil, errors.Wrap(err, "failed encode body")
				}

				return http.NewRequest(http.MethodPost, _url.String(), &buf)
			},
			assert: func(w *httptest.ResponseRecorder, r *http.Request) {
				assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
			},
			newService: func(req internal.RefreshRequest) service {
				return mockService{RefreshTokenFunc: func() (*internal.RefreshResponse, error) { return nil, httputil.ErrInternal }}
			},
		},
		{
			name: "200 - success",
			req: func(ctx context.Context) (*http.Request, error) {
				_url, err := url.Parse("http://example.com/api/v1/refresh")
				require.NoError(t, err)
				var buf bytes.Buffer
				if err := json.NewEncoder(&buf).Encode(dtos.RefreshRequest{RefreshToken: "token"}); err != nil {
					return nil, errors.Wrap(err, "failed encode body")
				}

				return http.NewRequest(http.MethodPost, _url.String(), &buf)
			},
			assert: func(w *httptest.ResponseRecorder, r *http.Request) {
				assert.Equal(t, http.StatusOK, w.Result().StatusCode)

				cookies := w.Result().Cookies()

				checkList := []string{"accessToken"}
				found := slices.ContainsFunc(cookies, func(c *http.Cookie) bool {
					if c.Value != "token" {
						return false
					}
					return slices.Contains(checkList, c.Name)
				})
				require.True(t, found)

			},
			newService: func(req internal.RefreshRequest) service {
				return mockService{RefreshTokenFunc: func() (*internal.RefreshResponse, error) {
					return &internal.RefreshResponse{
						AccessToken: "token",
					}, nil
				}}
			},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			newService = tt.newService
			ctx := context.Background()
			var err error

			w := httptest.NewRecorder()
			r, err := tt.req(ctx)
			if assert.NoError(t, err) {
				router.ServeHTTP(w, r)
				tt.assert(w, r)
			}
		})
	}
}
