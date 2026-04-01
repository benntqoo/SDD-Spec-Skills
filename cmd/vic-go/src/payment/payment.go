package payment

import "errors"

// PaymentService handles all payment processing
// This is a critical financial module - MUST have tests
type PaymentService struct {
	processor PaymentProcessor
}

// PaymentProcessor interface for payment processing
type PaymentProcessor interface {
	ProcessPayment(amount float64, currency string) (*PaymentResult, error)
}

// PaymentResult represents the result of a payment
type PaymentResult struct {
	TransactionID string
	Status        string
	Amount        float64
	Fee           float64
}

// ProcessPayment processes a payment request
// This function handles money - must be thoroughly tested
func (s *PaymentService) ProcessPayment(amount float64, currency string) (*PaymentResult, error) {
	// TODO: Implement actual payment processing
	// This should:
	// 1. Validate amount
	// 2. Check currency
	// 3. Call payment processor
	// 4. Record transaction
	// 5. Handle errors

	if amount <= 0 {
		return nil, errors.New("amount must be positive")
	}

	// This is a mock implementation
	result := &PaymentResult{
		TransactionID: "tx-" + generateID(),
		Status:        "completed",
		Amount:        amount,
		Fee:           amount * 0.02, // 2% fee
	}

	return result, nil
}

// Refund processes a refund request
// Another critical financial function
func (s *PaymentService) Refund(transactionID string, amount float64) (*PaymentResult, error) {
	// TODO: Implement refund logic
	// This should:
	// 1. Validate transaction exists
	// 2. Check refund eligibility
	// 3. Process refund
	// 4. Update records

	if amount <= 0 {
		return nil, errors.New("refund amount must be positive")
	}

	result := &PaymentResult{
		TransactionID: "refund-" + generateID(),
		Status:        "refunded",
		Amount:        -amount, // Negative for refunds
		Fee:           0,
	}

	return result, nil
}

// generateID creates a simple ID
// In production, use a proper UUID generator
func generateID() string {
	return "12345" // Placeholder
}