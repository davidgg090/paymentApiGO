package middlewares

import (
	"log"
	"net/http"
	"paymentApi/internal/model"
	"paymentApi/pkg/db"
	"strings"
	"time"
)

func AuditLogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)

		var userID *uint
		var bearerToken *string
		authHeader := r.Header.Get("Authorization")
		if authHeader != "" {
			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
			bearerToken = &tokenStr

			var token model.Token
			if err := db.DB.Where("access_token = ?", tokenStr).First(&token).Error; err == nil {
				userID = token.UserID
			} else {
				log.Printf("No se pudo encontrar un token válido: %v", err)

			}
		}

		auditLog := model.AuditLog{
			UserID:       userID,
			ActivityType: r.Method,
			BearerToken:  bearerToken,
			IPAddress:    &r.RemoteAddr,
			Path:         r.URL.Path,
			Timestamp:    time.Now(),
		}

		if result := db.DB.Create(&auditLog); result.Error != nil {
			log.Printf("Error al guardar el registro de auditoría: %v", result.Error)
		}
	})
}
