package middleware

import (
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"ph.certs.com/clm_client/auth"
	"strings"
)

// JWTMiddleware checks for the validity of the jwt token in the Authorization header
// extracts the userEmail, and adds it to the request context.
func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return auth.SecretKey, nil
		})
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if !token.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			ctx := r.Context()
			userEmail := claims["email"].(string)
			ctx = ContextWithUser(ctx, userEmail)
			req := r.WithContext(ctx)
			next.ServeHTTP(w, req)
		}

	})
}
