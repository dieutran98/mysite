package register

import (
	"bytes"
	"context"
	"encoding/json"
	"mysite/features/register/internal"
	"mysite/models/model"
	"mysite/utils/httputil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockService struct {
	service
	RegisterFunc func() error
}

func (m mockService) Register(ctx context.Context) error {
	return m.RegisterFunc()
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
		newService func(req internal.RegisterRequest) service
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
				_url, err := url.Parse("http://example.com/api/v1/register")
				require.NoError(t, err)
				var buf bytes.Buffer
				if err := json.NewEncoder(&buf).Encode(model.RegisterRequest{Password: "secret", UserName: "test@gmail.com"}); err != nil {
					return nil, errors.Wrap(err, "failed encode body")
				}

				return http.NewRequest(http.MethodPost, _url.String(), &buf)
			},
			assert: func(w *httptest.ResponseRecorder, r *http.Request) {
				assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
			},
			newService: func(req internal.RegisterRequest) service {
				return mockService{RegisterFunc: func() error { return httputil.ErrInvalidRequest }}
			},
		},
		{
			name: "400 - Invalid request, empty body",
			req: func(ctx context.Context) (*http.Request, error) {
				_url, err := url.Parse("http://example.com/api/v1/register")
				require.NoError(t, err)

				return http.NewRequest(http.MethodPost, _url.String(), nil)
			},
			assert: func(w *httptest.ResponseRecorder, r *http.Request) {
				assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
			},
			newService: func(req internal.RegisterRequest) service {
				return mockService{RegisterFunc: func() error { return nil }}
			},
		},
		{
			name: "500 - internal error request",
			req: func(ctx context.Context) (*http.Request, error) {
				_url, err := url.Parse("http://example.com/api/v1/register")
				require.NoError(t, err)
				var buf bytes.Buffer
				if err := json.NewEncoder(&buf).Encode(model.RegisterRequest{Password: "secret", UserName: "test@gmail.com"}); err != nil {
					return nil, errors.Wrap(err, "failed encode body")
				}

				return http.NewRequest(http.MethodPost, _url.String(), &buf)
			},
			assert: func(w *httptest.ResponseRecorder, r *http.Request) {
				assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
			},
			newService: func(req internal.RegisterRequest) service {
				return mockService{RegisterFunc: func() error { return httputil.ErrInternal }}
			},
		},
		{
			name: "201 - internal error request",
			req: func(ctx context.Context) (*http.Request, error) {
				_url, err := url.Parse("http://example.com/api/v1/register")
				require.NoError(t, err)
				var buf bytes.Buffer
				if err := json.NewEncoder(&buf).Encode(model.RegisterRequest{Password: "secret", UserName: "test@gmail.com"}); err != nil {
					return nil, errors.Wrap(err, "failed encode body")
				}

				return http.NewRequest(http.MethodPost, _url.String(), &buf)
			},
			assert: func(w *httptest.ResponseRecorder, r *http.Request) {
				assert.Equal(t, http.StatusCreated, w.Result().StatusCode)
			},
			newService: func(req internal.RegisterRequest) service {
				return mockService{RegisterFunc: func() error { return nil }}
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
