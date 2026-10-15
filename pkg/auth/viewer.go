package auth

import (
	"context"
	"net/http"
	"strings"

	"github.com/family-cpa/shared/pkg/gqlerr"
)

const (
	HeaderID      = "X-Actor-Id"
	HeaderScopes  = "X-Actor-Scopes"
	HeaderSession = "X-Session-Id"
)

type contextKey struct{}

type Viewer struct {
	ID        string
	Scopes    []string
	SessionID string
}

func FromContext(ctx context.Context) (*Viewer, error) {
	v, _ := ctx.Value(contextKey{}).(*Viewer)
	if v == nil {
		return nil, gqlerr.Unauthorized(ctx)
	}

	if len(v.ID) == 0 {
		return nil, gqlerr.Unauthorized(ctx)
	}

	return v, nil
}

func Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		viewer := &Viewer{
			ID: r.Header.Get(HeaderID),
			Scopes: strings.FieldsFunc(r.Header.Get(HeaderScopes), func(r rune) bool {
				return r == ' ' || r == ','
			}),
			SessionID: r.Header.Get(HeaderSession),
		}
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), contextKey{}, viewer)))
	})
}
