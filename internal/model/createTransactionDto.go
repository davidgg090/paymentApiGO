package model

type CreateTransactionRequest struct {
	MerchantID     uint    `json:"merchant_id"`
	CustomerID     uint    `json:"customer_id"`
	Amount         float64 `json:"amount"`
	Currency       string  `json:"currency"`
	HashCreditCard string  `json:"hash_credit_card"`
}
