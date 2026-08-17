package middleware

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"strings"

	"go-backend/pkg/utils"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const MerchantIDKey contextKey = "merchant_id"

// AuthMiddleware validates JWT tokens and injects merchant_id into request context
func AuthMiddleware(db *sql.DB) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            var tokenString string

            // 1. Try to extract token from HttpOnly cookie first
            cookie, err := r.Cookie("auth_token")
            if err == nil && cookie.Value != "" {
                tokenString = cookie.Value
            } else {
                // 2. Fall back to Authorization: Bearer <token> header
                authHeader := r.Header.Get("Authorization")
                if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
                    tokenString = strings.TrimPrefix(authHeader, "Bearer ")
                }
            }

            if tokenString == "" {
                utils.ErrorResponse(w, http.StatusUnauthorized, "Missing authentication token or cookie")
                return
            }

            // 3. Get JWT secret from environment variables
            jwtSecret := os.Getenv("JWT_SECRET")
            if jwtSecret == "" {
                jwtSecret = "go-backend_super_secret_key" // Fallback (ensure this matches your login generation secret)
            }

            // 4. Parse and validate JWT token
            token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
                // Ensure signing method is HMAC (e.g. HS256)
                if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                    return nil, jwt.ErrSignatureInvalid
                }
                return []byte(jwtSecret), nil
            })

            if err != nil || !token.Valid {
                utils.ErrorResponse(w, http.StatusUnauthorized, "Invalid or expired token")
                return
            }

            // 5. Extract claims (merchant_id)
            claims, ok := token.Claims.(jwt.MapClaims)
            if !ok {
                utils.ErrorResponse(w, http.StatusUnauthorized, "Invalid token claims")
                return
            }

            merchantID, ok := claims["merchant_id"].(string)
			fmt.Println("Thsi is themerchan ",merchantID)
            if !ok || merchantID == "" {
                utils.ErrorResponse(w, http.StatusUnauthorized, "Merchant ID missing from token claims")
                return
            }

            // 6. Inject merchant_id into request context
            ctx := context.WithValue(r.Context(), MerchantIDKey, merchantID)
            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}

// GetMerchantIDFromContext retrieves merchant ID in controllers/handlers safely
func GetMerchantIDFromContext(ctx context.Context) (string, bool) {
    id, ok := ctx.Value(MerchantIDKey).(string)
    return id, ok
}
