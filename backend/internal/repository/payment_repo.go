package repository

import (
	"context"
	"time"

	"pos-saas/internal/domain"
)

// PaymentRepository handles payment data access.
type PaymentRepository interface {
	CreatePayment(ctx context.Context, payment *domain.Payment) (int64, error)
	GetPaymentByID(ctx context.Context, tenantID, id int64) (*domain.Payment, error)
	GetPaymentTenantID(ctx context.Context, id int64) (int64, error)
	GetPaymentByMerchantOrderID(ctx context.Context, tenantID int64, merchantOrderID string) (*domain.Payment, error)
	GetPaymentByPaymobOrderID(ctx context.Context, paymobOrderID int64, merchantOrderID string) (*domain.Payment, error)
	UpdatePaymentOrder(ctx context.Context, tenantID, id int64, paymobOrderID int64, authToken string) error
	UpdatePaymentStatus(ctx context.Context, tenantID, id int64, status string, txnID *int64, cardLastFour *string, cardBrand *string, errorMsg *string, paidAt *time.Time) error
	GetPaymentsByTenantID(ctx context.Context, tenantID int64, limit, offset int) ([]domain.Payment, int, error)
	GetPaymentsByTenantCursor(ctx context.Context, tenantID int64, cursorID int64, limit int) ([]domain.Payment, *int64, error)
	UpdateTenantSubscription(ctx context.Context, tenantID int64, status string, expiresAt time.Time, lastPaymentAt time.Time, lastPaymentAmount int, currency string, planID string) error
	GetTenantBillingInfo(ctx context.Context, tenantID int64) (*domain.TenantBillingInfo, error)
}
