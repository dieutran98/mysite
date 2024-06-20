package httputil

import (
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
)

func SetCookie(key, value string) *http.Cookie {
	return &http.Cookie{
		Name:     key,
		Value:    value,
		Secure:   true,
		HttpOnly: true,
	}
}

func ParseBody[T comparable](r *http.Request, body *T) error {
	if r.Body == nil || r.Body == http.NoBody {
		return errors.Wrap(ErrInvalidRequest, "empty body")
	}

	if err := json.NewDecoder(r.Body).Decode(body); err != nil {
		return errors.Wrap(err, "failed to decode body")
	}

	return nil
}
