package middleware

import (
	"context"
	"database/sql"
	"net/http"
	"strings"

	"go-backend/pkg/utils"
)

func APIKeyMiddleware(db *sql.DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			// Get API key
			apiKey := r.Header.Get("X-API-Key")

			if apiKey == "" {
				authHeader := r.Header.Get("Authorization")

				if strings.HasPrefix(authHeader, "Bearer ") {
					apiKey = strings.TrimSpace(
						strings.TrimPrefix(authHeader, "Bearer "),
					)
				}
			}

			if apiKey == "" {
				utils.ErrorResponse(
					w,
					http.StatusUnauthorized,
					"API key required",
				)
				return
			}

			// Find API key directly
			var merchantID string
			var isActive bool

			err := db.QueryRow(`
				SELECT merchant_id, is_active
				FROM api_keys
				WHERE key_prefix = $1
			`, apiKey).Scan(
				&merchantID,
				&isActive,
			)

			if err == sql.ErrNoRows {
				utils.ErrorResponse(
					w,
					http.StatusUnauthorized,
					"Invalid API key",
				)
				return
			}

			if err != nil {
				utils.ErrorResponse(
					w,
					http.StatusInternalServerError,
					"Internal server error",
				)
				return
			}

			if !isActive {
				utils.ErrorResponse(
					w,
					http.StatusUnauthorized,
					"API key is deactivated",
				)
				return
			}

			// Put merchant ID into context
			ctx := context.WithValue(
				r.Context(),
				MerchantIDKey,
				merchantID,
			)

			next.ServeHTTP(
				w,
				r.WithContext(ctx),
			)
		})
	}
}
