package bankutils

import (
	"errors"
	"fmt"
	"regexp"
	"time"
)

func VerifyHashCreditCard(hashTransaction, customerHash string) bool {
	return hashTransaction == customerHash
}

func SendPaymentDetails(amount float64, cardInfo map[string]string, merchantID string) (map[string]string, error) {
	fmt.Printf("Processing $%.2f payment from card %s to merchant %s\n", amount, cardInfo["card_number"], merchantID)
	return map[string]string{"status": "success", "transaction_id": "12345ABC"}, nil
}

func ProcessPayment(cardInfo map[string]string, amount float64, merchantID string) (map[string]string, error) {
	if !ValidateCard(cardInfo["card_number"], cardInfo["expiration_date"], cardInfo["cvv"]) {
		return nil, errors.New("invalid card details")
	}

	result, err := SendPaymentDetails(amount, cardInfo, merchantID)
	if err != nil {
		return nil, err
	}

	if result["status"] == "success" {
		fmt.Printf("Payment successful. Transaction ID: %s\n", result["transaction_id"])
		return result, nil
	}

	return nil, errors.New("payment failed")
}

func ValidateCard(cardNumber, expirationDate, cvv string) bool {
	if !validateCardNumber(cardNumber) {
		return false
	}
	if !validateExpirationDate(expirationDate) {
		return false
	}
	if !validateCVV(cvv) {
		return false
	}
	return true
}

func validateCardNumber(number string) bool {
	re := regexp.MustCompile(`^\d{13,19}$`)
	return re.MatchString(number)
}

func validateExpirationDate(date string) bool {
	t, err := time.Parse("01/06", date)
	if err != nil {
		return false
	}

	t = t.AddDate(0, 1, 0)

	return t.After(time.Now())
}

func validateCVV(cvv string) bool {
	re := regexp.MustCompile(`^\d{3,4}$`)
	return re.MatchString(cvv)
}
