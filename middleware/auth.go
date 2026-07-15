package middleware

import (
	"context"
	"net/http"
	"strings"

	"eduplanner/utils"
)

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			utils.SendError(
				w,
				http.StatusUnauthorized,
				"Missing authorization Header",
			)
			return
		}
		parts := strings.Split(authHeader, " ")

		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.SendError(w, http.StatusUnauthorized, "Invalid Authorization Header")
			return
		}

		tokenString := parts[1]

		claims, err := utils.ValidateToken(tokenString)

		if err != nil {
			utils.SendError(w, http.StatusUnauthorized, "Invalid or expired token")
			return
		}
		ctx := context.WithValue(r.Context(), "userID", claims.UserID)
		ctx = context.WithValue(ctx, "email", claims.Email)
		ctx = context.WithValue(ctx, "role", claims.Role)

		next.ServeHTTP(w, r.WithContext(ctx))

	})

}
