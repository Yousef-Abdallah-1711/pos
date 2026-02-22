package service

// PaymobClient defines Paymob integration operations.
type PaymobClient interface {
	Authenticate() (string, error)
	CreateOrder(authToken string, amountCents int, currency string, merchantOrderID string) (int64, error)
	GetPaymentKey(authToken string, orderID int64, amountCents int, currency string, integrationID int, billingData map[string]string, expiration int) (string, error)
	RefundTransaction(authToken string, transactionID int64, amountCents int) error
	GetIframeURL(paymentKey string) string
	VerifyHMAC(payload []byte, receivedHMAC string) bool
}
