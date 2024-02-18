package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"paymentApi/internal/config"
	"paymentApi/internal/model"
	"paymentApi/pkg/db"
	"time"
)

var jwtKey = []byte(config.JWTKey)

// Login autentica a un usuario y genera un JWT si las credenciales son válidas.
// @Summary Inicia sesión de usuario
// @Description Autentica a un usuario basado en `username` y `password` y genera un JWT.
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body model.loginCredentialsDTO true "Credenciales de inicio de sesión"
// @Success 200 {object} map[string]string "Token generado exitosamente"
// @Failure 400 {string} string "Credenciales inválidas"
// @Failure 401 {string} string "Usuario no encontrado o contraseña incorrecta"
// @Router /login [post]
func Login(w http.ResponseWriter, r *http.Request) {
	var credentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		http.Error(w, "Credenciales inválidas", http.StatusBadRequest)
		return
	}

	var user model.User
	result := db.DB.Where("username = ?", credentials.Username).First(&user)
	if result.Error != nil {
		http.Error(w, "Usuario no encontrado", http.StatusUnauthorized)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password)); err != nil {
		http.Error(w, "Contraseña incorrecta", http.StatusUnauthorized)
		return
	}

	expirationTime := time.Now().Add(1 * time.Hour)
	claims := &jwt.StandardClaims{
		Subject:   fmt.Sprintf("%d", user.ID),
		ExpiresAt: expirationTime.Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		http.Error(w, "Error al generar el token", http.StatusInternalServerError)
		return
	}

	db.DB.Create(&model.Token{AccessToken: tokenString, User: user, TokenType: "bearer",
		CreatedAt: time.Now(), ExpiresAt: expirationTime, UpdatedAt: time.Now()})

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(map[string]string{
		"access_token": tokenString,
		"token_type":   "bearer",
	})
	if err != nil {
		return
	}
}
