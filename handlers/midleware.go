package handlers

import (
	"context"
	"net/http"
	"strings"

	internalCtx "project/internal/context"
	"project/pkg/jwt"
)

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")

		if auth == "" {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		if !strings.HasPrefix(auth, "Bearer ") {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		token := strings.TrimPrefix(auth, "Bearer ")

		claims, err := jwt.ParseToken(token)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), internalCtx.RoleKey, claims.Role)
		ctx = context.WithValue(ctx, internalCtx.UserIDKey, claims.UserID)
		ctx = context.WithValue(ctx, internalCtx.EmailKey, claims.Email)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
