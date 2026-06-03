package authentication

import (
	"context"
	"net/http"
)

func AuthMiddleware(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Missing token",
					http.StatusUnauthorized)
				return
			}

			id, err := ParseToken(secret, authHeader)
			if err != nil {
				http.Error(w, "Invalid token",
					http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), "id", id)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetUserFromContext(ctx context.Context) string {
	id, _ := ctx.Value("id").(string)

	return id
}
