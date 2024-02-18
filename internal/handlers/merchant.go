package handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"paymentApi/internal/model"
	"paymentApi/pkg/db"
)

// CreateMerchant crea un nuevo merchant
// @Summary Crea un nuevo merchant
// @Description Agrega un nuevo merchant a la base de datos
// @Tags merchant
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param merchant body model.Merchant true "Información del merchant"
// @Success 201 {object} model.Merchant "Merchant creado"
// @Failure 400 {object} map[string]interface{} "Entrada inválida"
// @Failure 500 {object} map[string]interface{} "Error interno del servidor"
// @Router /api/merchant [post]
func CreateMerchant(w http.ResponseWriter, r *http.Request) {
	var merchat model.Merchant

	if err := json.NewDecoder(r.Body).Decode(&merchat); err != nil {
		log.Printf("Error al decodificar el cuerpo de la solicitud: %v", err)
		http.Error(w, `{"error": "Error al procesar los datos de entrada"}`, http.StatusBadRequest)
		return
	}

	result := db.DB.Create(&merchat)
	if result.Error != nil {
		log.Printf("Error al crear el merchant en la base de datos: %v", result.Error)
		http.Error(w, `{"error": "Error al crear el merchant"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(merchat); err != nil {
		log.Printf("Error al codificar la respuesta JSON: %v", err)
	}
}

// GetMerchant obtiene un merchant por ID
// @Summary Obtiene un merchant
// @Description Obtiene un merchant por su ID
// @Tags merchant
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param   id     path      int  true  "ID del Merchant"
// @Success 200  {object}  model.Merchant
// @Failure 404 {object} map[string]interface{} "Merchant no encontrado"
// @Router /api/merchant/{id} [get]
func GetMerchant(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	var merchat model.Merchant
	result := db.DB.First(&merchat, id)
	if result.Error != nil {
		log.Printf("Error al obtener el merchant con ID %s: %v", id, result.Error)
		http.Error(w, `{"error": "Merchant no encontrado"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(merchat); err != nil {
		log.Printf("Error al codificar la respuesta JSON para el merchant con ID %s: %v", id, err)
	}
}

// UpdateMerchant actualiza un merchant por ID
// @Summary Actualiza un merchant
// @Description Actualiza la información de un merchant existente por su ID
// @Tags merchant
// @Accept  json
// @Produce  json
// @Param   id     path      int  true  "ID del Merchant"
// @Param   merchant  body      model.Merchant true  "Información actualizada del Merchant"
// @Success 200  {object}  model.Merchant
// @Failure 400 {object} map[string]interface{} "Entrada inválida"
// @Failure 404 {object} map[string]interface{} "Merchant no encontrado"
// @Router /api/merchant/{id} [put]
func UpdateMerchant(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var updateInfo model.Merchant
	if err := json.NewDecoder(r.Body).Decode(&updateInfo); err != nil {
		log.Printf("Error al decodificar la solicitud de actualización: %v", err)
		http.Error(w, `{"error": "Entrada inválida"}`, http.StatusBadRequest)
		return
	}

	var merchant model.Merchant
	result := db.DB.First(&merchant, id)
	if result.Error != nil {
		log.Printf("Merchant con ID %s no encontrado: %v", id, result.Error)
		http.Error(w, `{"error": "Merchant no encontrado"}`, http.StatusNotFound)
		return
	}

	result = db.DB.Model(&merchant).Updates(updateInfo)
	if result.Error != nil {
		log.Printf("Error al actualizar el merchant con ID %s: %v", id, result.Error)
		http.Error(w, `{"error": "Error al actualizar el merchant"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(merchant)
	if err != nil {
		return
	}
}

// DeleteMerchant elimina un merchant por ID
// @Summary Elimina un merchant
// @Description Elimina un merchant por su ID
// @Tags merchant
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param   id     path      int  true  "ID del Merchant"
// @Success 204 "Merchant eliminado exitosamente"
// @Failure 404 {object} map[string]interface{} "Merchant no encontrado"
// @Router /api/merchant/{id} [delete]
func DeleteMerchant(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	result := db.DB.Delete(&model.Merchant{}, id)
	if result.Error != nil {
		log.Printf("Error al intentar eliminar el merchant con ID %s: %v", id, result.Error)
		http.Error(w, `{"error": "Error interno del servidor"}`, http.StatusInternalServerError)
		return
	}
	if result.RowsAffected == 0 {
		log.Printf("Intento de eliminar un merchant no encontrado con ID %s", id)
		http.Error(w, `{"error": "Merchant no encontrado"}`, http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
