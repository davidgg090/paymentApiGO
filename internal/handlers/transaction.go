package handlers

import (
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"log"
	"net/http"
	"paymentApi/internal/model"
	bankutils "paymentApi/internal/utils"
	"paymentApi/pkg/db"
)

// CreateTransaction crea una nueva transacción
// @Summary Crea una nueva transacción
// @Description Crea una nueva transacción con los datos proporcionados
// @Tags transaction
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param transaction body model.CreateTransactionRequest true "Datos de la Transacción"
// @Success 201 {object} model.Transaction
// @Router /api/transaction [post]
func CreateTransaction(w http.ResponseWriter, r *http.Request) {
	var req model.CreateTransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error al decodificar el cuerpo de la solicitud: %v", err)
		http.Error(w, `{"error": "Entrada inválida"}`, http.StatusBadRequest)
		return
	}

	transaction := model.Transaction{
		MerchantID:     req.MerchantID,
		CustomerID:     req.CustomerID,
		Amount:         req.Amount,
		Currency:       req.Currency,
		HashCreditCard: req.HashCreditCard,
		State:          "pending",
		Token:          uuid.New().String(),
	}

	result := db.DB.Create(&transaction)
	if result.Error != nil {
		log.Printf("Error al crear la transacción en la base de datos: %v", result.Error)
		http.Error(w, `{"error": "Error al crear la transacción"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(transaction); err != nil {
		log.Printf("Error al codificar la respuesta JSON: %v", err)
	}
}

// GetTransaction obtiene una transaccion por ID
// @Summary Obtiene una transacción
// @Description Obtiene una transaccion por su ID
// @Tags transaction
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param   id     path      int  true  "ID del Transaction"
// @Success 200  {object}  model.Transaction "Transaccion encontrada"
// @Failure 404 {object} map[string]interface{} "Transaccion no encontrada"
// @Router /api/transaction/{id} [get]
func GetTransaction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	var transaction model.Transaction
	result := db.DB.First(&transaction, id)
	if result.Error != nil {
		log.Printf("Error al obtener la transaccion con ID %s: %v", id, result.Error)
		http.Error(w, `{"error": "Transaccion no encontrado"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(transaction); err != nil {
		log.Printf("Error al codificar la respuesta JSON para el customer con ID %s: %v", id, err)
	}
}

// UpdateTransaction actualiza una transacción por ID
// @Summary Actualiza una transacción
// @Description Actualiza una transacción existente; solo se permite si el estado es "pending"
// @Tags transaction
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "ID de la Transacción"
// @Param transaction body model.TransactionUpdateDTO true "Datos de la Transacción a actualizar"
// @Success 200 "Transacción actualizada"
// @Failure 400 "Entrada inválida"
// @Failure 404 "Transacción no encontrada"
// @Router /api/transaction/{id} [put]
func UpdateTransaction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var updateInfo model.TransactionUpdateDTO
	if err := json.NewDecoder(r.Body).Decode(&updateInfo); err != nil {
		log.Printf("Error al decodificar la solicitud de actualización: %v", err)
		http.Error(w, `{"error": "Entrada inválida"}`, http.StatusBadRequest)
		return
	}

	var transaction model.Transaction
	result := db.DB.First(&transaction, id)
	if result.Error != nil {
		log.Printf("Transacción con ID %s no encontrada: %v", id, result.Error)
		http.Error(w, `{"error": "Transacción no encontrada"}`, http.StatusNotFound)
		return
	}

	if transaction.State != "pending" {
		log.Printf("Intento de actualizar una transacción no pendiente con ID %s", id)
		http.Error(w, `{"error": "Solo se pueden actualizar transacciones en estado 'pending'"}`, http.StatusBadRequest)
		return
	}

	result = db.DB.Model(&transaction).Updates(updateInfo)
	if result.Error != nil {
		log.Printf("Error al actualizar la transacción con ID %s: %v", id, result.Error)
		http.Error(w, `{"error": "Error al actualizar la transacción"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func GetTransactionCustomer(customerID uint) (*model.Customer, error) {
	var customer model.Customer
	result := db.DB.First(&customer, customerID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, result.Error
		}
		return nil, result.Error
	}
	return &customer, nil
}

func GetMerchantByID(merchantID uint) (*model.Merchant, error) {
	var merchant model.Merchant
	result := db.DB.First(&merchant, merchantID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, result.Error
		}
		return nil, result.Error
	}
	return &merchant, nil
}

func updateTransactionState(transactionID uint, newState string) error {
	var transaction model.Transaction
	result := db.DB.First(&transaction, transactionID)
	if result.Error != nil {
		return result.Error
	}

	transaction.State = newState
	result = db.DB.Save(&transaction)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func updateTransactionAndMerchant(transaction *model.Transaction, merchant *model.Merchant) error {
	tx := db.DB.Begin()

	if err := tx.Error; err != nil {
		return err
	}

	if err := tx.Save(transaction).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Save(merchant).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

// ProcessTransaction procesa una transacción basada en su token y tipo.
// @Summary Procesa una transacción
// @Description Procesa una transacción dada, capturándola o reembolsándola basado en el tipo de transacción especificado.
// @Tags transaction
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param token query string true "Token de la transacción"
// @Param type query string true "Tipo de transacción (capture, refund)"
// @Success 200 {object} model.Transaction "Transacción procesada"
// @Failure 400 {object} map[string]interface{} "Error de solicitud - Transacción ya procesada, tipo de transacción inválido, tarjeta de crédito inválida, etc."
// @Failure 404 {object} map[string]interface{} "Transacción no encontrada"
// @Failure 500 {object} map[string]interface{} "Error interno del servidor"
// @Router /api/transaction/process [post]
func ProcessTransaction(w http.ResponseWriter, r *http.Request) {
	transactionToken := r.URL.Query().Get("token")
	typeTransaction := r.URL.Query().Get("type") // "capture" o "refund"

	var transaction model.Transaction
	result := db.DB.Where("token = ?", transactionToken).First(&transaction)
	if result.Error != nil {
		http.Error(w, `{"error": "Transaction not found"}`, http.StatusNotFound)
		return
	}
	if transaction.State != "pending" {
		http.Error(w, `{"error": "Transaction already processed"}`, http.StatusBadRequest)
		return
	}

	customer, err := GetTransactionCustomer(transaction.CustomerID)
	if err != nil || !customer.IsActive {
		err := updateTransactionState(transaction.ID, "failed")
		if err != nil {
			return
		}
		http.Error(w, `{"error": "Customer not found or not active"}`, http.StatusBadRequest)
		return
	}

	merchant, err := GetMerchantByID(transaction.MerchantID)
	if err != nil || !merchant.IsActive {
		err := updateTransactionState(transaction.ID, "failed")
		if err != nil {
			return
		}
		http.Error(w, `{"error": "Merchant not found or not active"}`, http.StatusBadRequest)
		return
	}

	if bankutils.VerifyHashCreditCard(transaction.HashCreditCard, customer.HashCreditCard) {
		switch typeTransaction {
		case "capture":
			transaction.State = "success"
			merchant.AmountAccount += transaction.Amount
		case "refund":
			transaction.State = "refunded"
			merchant.AmountAccount -= transaction.Amount
		default:
			http.Error(w, `{"error": "Invalid transaction type"}`, http.StatusBadRequest)
			return
		}
	} else {
		http.Error(w, `{"error": "Invalid credit card"}`, http.StatusBadRequest)
		return
	}

	if err := updateTransactionAndMerchant(&transaction, merchant); err != nil {
		http.Error(w, `{"error": "Failed to update transaction or merchant"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(transaction)
	if err != nil {
		return
	}
}

func transactionByToken(token string) (*model.Transaction, error) {
	var transaction model.Transaction
	result := db.DB.Where("token = ?", token).First(&transaction)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, result.Error
		}
		return nil, result.Error
	}
	return &transaction, nil
}

// GetTransactionByToken busca y retorna una transacción por su token.
// @Summary Obtiene una transacción por token
// @Description Obtiene una transacción por su token único
// @Tags transaction
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param   token   path      string  true  "Token de la Transacción"
// @Success 200  {object}  model.Transaction
// @Failure 404  {object}  map[string]interface{} "Transaction not found"
// @Failure 500  {object}  map[string]interface{} "Internal server error"
// @Router /api/transaction/token/{token} [get]
func GetTransactionByToken(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	token := vars["token"]
	var transaction model.Transaction

	result := db.DB.Where(`"token" = ?`, token).First(&transaction)
	if result.Error != nil {
		log.Printf("Error al obtener la transacción con token %s: %v", token, result.Error)
		http.Error(w, `{"error": "Transacción no encontrada"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(transaction); err != nil {
		log.Printf("Error al codificar la respuesta JSON para la transacción con token %s: %v", token, err)
		http.Error(w, "Error al procesar la respuesta", http.StatusInternalServerError)
	}
}
