package domain

import "errors"

var ErrInvalidOrderQuote = errors.New("invalid order quote")

// OrderQuoteRequest carries no customer PII or client-provided prices.
type OrderQuoteRequest struct {
	Items []CreateOrderItemRequest `json:"items"`
}

type OrderQuoteItem struct {
	ProductID  int64   `json:"product_id"`
	Quantity   int     `json:"quantity"`
	UnitPrice  float64 `json:"unit_price"`
	TotalPrice float64 `json:"total_price"`
}

type OrderQuote struct {
	Items       []OrderQuoteItem `json:"items"`
	Subtotal    float64          `json:"subtotal"`
	TaxAmount   float64          `json:"tax_amount"`
	DeliveryFee float64          `json:"delivery_fee"`
	TotalAmount float64          `json:"total_amount"`
	Currency    string           `json:"currency"`
}
