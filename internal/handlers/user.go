package handlers

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"net/http"
	"paymentApi/internal/model"
	"paymentApi/pkg/db"
	"strconv"
)

// CreateUser crea un nuevo usuario en la base de datos.
// @Summary Crea un nuevo usuario
// @Description Crea un nuevo usuario con username y password proporcionados.
// @Tags user
// @Accept json
// @Produce json
// @Param credentials body model.createUserDTO true "Credenciales de usuario"
// @Success 201 {object} map[string]interface{} "Usuario creado con ID y username"
// @Failure 400 {string} string "Error al decodificar solicitud o usuario/password faltantes"
// @Failure 409 {string} string "Nombre de usuario ya existe"
// @Failure 500 {string} string "Error al hashear la contraseña / Error al crear el usuario"
// @Router /user [post]
func CreateUser(w http.ResponseWriter, r *http.Request) {
	var credentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		http.Error(w, "Error al decodificar solicitud", http.StatusBadRequest)
		return
	}

	if credentials.Username == "" || credentials.Password == "" {
		http.Error(w, "Usuario y contraseña son requeridos", http.StatusBadRequest)
		return
	}

	var exists model.User
	if err := db.DB.Where("username = ?", credentials.Username).First(&exists).Error; err == nil {
		http.Error(w, "Nombre de usuario ya existe", http.StatusConflict)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(credentials.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error al hashear la contraseña", http.StatusInternalServerError)
		return
	}

	user := model.User{
		Username: credentials.Username,
		Password: string(hashedPassword),
	}

	if result := db.DB.Create(&user); result.Error != nil {
		http.Error(w, "Error al crear el usuario", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"id":       user.ID,
		"username": user.Username,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		return
	}
}

// GetUser obtiene un usuario por su ID.
// @Summary Obtiene un usuario
// @Description Obtiene un usuario por su ID, excluyendo la contraseña en la respuesta.
// @Tags user
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "ID del Usuario"
// @Success 200 {object} model.User "Usuario encontrado"
// @Failure 401 {string} string "No autorizado"
// @Failure 404 {string} string "Usuario no encontrado"
// @Failure 500 {string} string "Error al obtener el usuario"
// @Router /api/user/{id} [get]
func GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		http.Error(w, "ID de usuario inválido", http.StatusBadRequest)
		return
	}

	var user model.User
	result := db.DB.Select("ID", "Username", "CreatedAt", "UpdatedAt").First(&user, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			http.NotFound(w, r)
			return
		}
		http.Error(w, "Error al obtener el usuario", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, "Error al codificar la respuesta", http.StatusInternalServerError)
	}
}

// UpdateUser actualiza un usuario existente por su ID.
// @Summary Actualiza un usuario
// @Description Actualiza la información de un usuario existente por su ID.
// @Tags user
// @Accept json
// @Security ApiKeyAuth
// @Produce json
// @Param id path int true "ID del Usuario"
// @Param user body model.User true "Información actualizada del Usuario"
// @Success 200 {object} model.User
// @Failure 400 {object} map[string]interface{} "Error en la solicitud"
// @Failure 404 {object} map[string]interface{} "Usuario no encontrado"
// @Router /api/user/{id} [put]
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var updateUser model.User
	if err := json.NewDecoder(r.Body).Decode(&updateUser); err != nil {
		http.Error(w, "Error al procesar los datos de entrada", http.StatusBadRequest)
		return
	}

	if updateUser.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(updateUser.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Error al hashear la contraseña", http.StatusInternalServerError)
			return
		}
		updateUser.Password = string(hashedPassword)
	}

	var user model.User
	if result := db.DB.First(&user, id); result.Error != nil {
		http.Error(w, "Usuario no encontrado", http.StatusNotFound)
		return
	}

	if result := db.DB.Model(&user).Updates(updateUser); result.Error != nil {
		http.Error(w, "Error al actualizar el usuario", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(user)
	if err != nil {
		return
	}
}

// DeleteUser elimina un usuario por ID.
// @Summary Elimina un usuario
// @Description Elimina un usuario por su ID.
// @Tags user
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "ID del Usuario"
// @Success 204 "Usuario eliminado correctamente"
// @Failure 404 {object} map[string]interface{} "Usuario no encontrado"
// @Router /api/user/{id} [delete]
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if result := db.DB.Delete(&model.User{}, id); result.Error != nil {
		http.Error(w, "Usuario no encontrado", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}
