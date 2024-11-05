// internal/middleware/auth.go
package middleware

import (
    "context"
    "log"
    "net/http"
    "strings"
    
    "github.com/golang-jwt/jwt/v4"
    "tech-test/backend/internal/domain"
    "tech-test/backend/internal/utils"
)

type ContextKey string

const UserIDKey ContextKey = "userID"

func AuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        log.Printf("Processing request: %s %s", r.Method, r.URL.Path)
        
        authHeader := r.Header.Get("Authorization")
        log.Printf("Auth header: %s", authHeader)
        
        if authHeader == "" {
            log.Printf("No auth header found")
            utils.RespondWithError(w, domain.NewAPIError(
                http.StatusUnauthorized,
                domain.ErrCodeAuthentication,
                "No authorization header",
                nil,
            ))
            return
        }

        tokenParts := strings.Split(authHeader, " ")
        if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
            log.Printf("Invalid auth header format: %s", authHeader)
            utils.RespondWithError(w, domain.NewAPIError(
                http.StatusUnauthorized,
                domain.ErrCodeAuthentication,
                "Invalid authorization header format",
                nil,
            ))
            return
        }

        token := tokenParts[1]
        claims, err := utils.ValidateToken(token, utils.JWTSecret, jwt.SigningMethodHS256)
        if err != nil {
            log.Printf("Token validation failed: %v", err)
            utils.RespondWithError(w, domain.NewAPIError(
                http.StatusUnauthorized,
                domain.ErrCodeAuthentication,
                "Invalid token",
                err,
            ))
            return
        }

        log.Printf("Token validated successfully for user ID: %d", claims.UserID)
        ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

func GetUserIDFromContext(ctx context.Context) (uint, bool) {
    userID, ok := ctx.Value(UserIDKey).(uint)
    return userID, ok
}
