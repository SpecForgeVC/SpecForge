package middleware

import (
	"context"

	"github.com/scott/specforge/internal/domain"
)

type contextKey string

const principalKey contextKey = "principal"

// PrincipalFromContext extracts the principal from the context.
func PrincipalFromContext(ctx context.Context) (*domain.Principal, bool) {
	p, ok := ctx.Value(principalKey).(*domain.Principal)
	return p, ok
}

// ContextWithPrincipal injects the principal into the context.
func ContextWithPrincipal(ctx context.Context, p *domain.Principal) context.Context {
	return context.WithValue(ctx, principalKey, p)
}
