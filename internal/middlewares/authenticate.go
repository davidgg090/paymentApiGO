package middlewares

import (
	"net/http"
	"paymentApi/internal/config"
	"paymentApi/internal/model"
	"paymentApi/pkg/db"
	"strings"
	"time"
)

func Authenticate(next http.Handler) http.Handler {
	config.LoadConfig()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := strings.Split(r.Header.Get("Authorization"), "Bearer ")
		http.Header{}.Set("Content-Type", "application/json")
		if len(authHeader) != 2 {
			http.Error(w, "Token de autorización no proporcionado", http.StatusUnauthorized)
			return
		}

		tokenString := authHeader[1]

		var tokenRecord model.Token
		if err := db.DB.Where("access_token = ?", tokenString).First(&tokenRecord).Error; err != nil {
			http.Error(w, "Token inválido", http.StatusUnauthorized)
			return
		}

		if time.Now().After(tokenRecord.ExpiresAt) {
			http.Error(w, "Token expirado", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
