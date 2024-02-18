package model

type TransactionUpdateDTO struct {
	Amount   float64 `json:"amount,omitempty"`
	Currency string  `json:"currency,omitempty"`
	State    string  `json:"state,omitempty"`
}
