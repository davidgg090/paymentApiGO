package handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"paymentApi/internal/model"
	"paymentApi/pkg/db"
)

// CreateCustomer crea un nuevo customer
// @Summary Crea un nuevo customer
// @Description Agrega un nuevo customer a la base de datos con la información proporcionada en el cuerpo de la solicitud.
// @Tags customer
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param customer body model.Customer true "Información del customer"
// @Success 201 {object} model.Customer "Customer creado exitosamente"
// @Failure 400 {object} map[string]interface{} "Error al procesar los datos de entrada"
// @Failure 500 {object} map[string]interface{} "Error interno del servidor al crear el customer"
// @Router /api/customer [post]
func CreateCustomer(w http.ResponseWriter, r *http.Request) {
	var customer model.Customer

	if err := json.NewDecoder(r.Body).Decode(&customer); err != nil {
		log.Printf("Error al decodificar el cuerpo de la solicitud: %v", err)
		http.Error(w, `{"error": "Error al procesar los datos de entrada"}`, http.StatusBadRequest)
		return
	}

	result := db.DB.Create(&customer)
	if result.Error != nil {
		log.Printf("Error al crear el customer en la base de datos: %v", result.Error)
		http.Error(w, `{"error": "Error al crear el customer"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(customer); err != nil {
		log.Printf("Error al codificar la respuesta JSON: %v", err)
	}
}

// GetCustomer obtiene un customer por ID
// @Summary Obtiene un customer
// @Description Obtiene un customer por su ID
// @Tags customer
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param   id     path      int  true  "ID del Customer"
// @Success 200  {object}  model.Customer "Customer encontrado"
// @Failure 404 {object} map[string]interface{} "Customer no encontrado"
// @Router /api/customer/{id} [get]
func GetCustomer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	var customer model.Customer
	result := db.DB.First(&customer, id)
	if result.Error != nil {
		log.Printf("Error al obtener el customer con ID %s: %v", id, result.Error)
		http.Error(w, `{"error": "Customer no encontrado"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(customer); err != nil {
		log.Printf("Error al codificar la respuesta JSON para el customer con ID %s: %v", id, err)
	}
}

// UpdateCustomer actualiza un customer por ID
// @Summary Actualiza un customer
// @Description Actualiza la información de un customer existente por su ID
// @Tags customer
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "ID del Customer"
// @Param customer body model.Customer true "Información actualizada del Customer"
// @Success 200 {object} model.Customer "Customer actualizado"
// @Failure 400 {object} map[string]interface{} "Entrada inválida"
// @Failure 404 {object} map[string]interface{} "Customer no encontrado"
// @Router /api/customer/{id} [put]
func UpdateCustomer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var updateInfo model.Customer
	if err := json.NewDecoder(r.Body).Decode(&updateInfo); err != nil {
		log.Printf("Error al decodificar la solicitud de actualización: %v", err)
		http.Error(w, `{"error": "Entrada inválida"}`, http.StatusBadRequest)
		return
	}

	var customer model.Customer
	result := db.DB.First(&customer, id)
	if result.Error != nil {
		log.Printf("Error al buscar el customer con ID %s: %v", id, result.Error)
		http.Error(w, `{"error": "Customer no encontrado"}`, http.StatusNotFound)
		return
	}

	result = db.DB.Model(&customer).Updates(updateInfo)
	if result.Error != nil {
		log.Printf("Error al actualizar el customer con ID %s: %v", id, result.Error)
		http.Error(w, `{"error": "Error al actualizar el customer"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(customer); err != nil {
		log.Printf("Error al codificar la respuesta JSON para el customer actualizado: %v", err)
	}
}

// DeleteCustomer elimina un customer por ID
// @Summary Elimina un customer
// @Description Elimina un customer por su ID
// @Tags customer
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "ID del Customer"
// @Success 204 "Customer eliminado exitosamente"
// @Failure 404 {object} map[string]interface{} "Customer no encontrado"
// @Failure 500 {object} map[string]interface{} "Error interno del servidor"
// @Router /api/customer/{id} [delete]
func DeleteCustomer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	result := db.DB.Delete(&model.Customer{}, id)
	if result.Error != nil {
		log.Printf("Error al intentar eliminar el customer con ID %s: %v", id, result.Error)
		http.Error(w, `{"error": "Error interno del servidor"}`, http.StatusInternalServerError)
		return
	}
	if result.RowsAffected == 0 {
		log.Printf("Intento de eliminar un customer no encontrado con ID %s", id)
		http.Error(w, `{"error": "Customer no encontrado"}`, http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
