package usecase

import (
	"context"
	"errors"
	"testing"

	"pos-saas/internal/domain"
	"pos-saas/internal/repository"
)

var errOfferCatalogUnavailable = errors.New("synthetic offer catalog outage")

type unavailableOfferCatalog struct {
	repository.PlanRepository
}

func (unavailableOfferCatalog) GetAllActive(context.Context) ([]domain.SubscriptionPlan, error) {
	return nil, errOfferCatalogUnavailable
}

func (unavailableOfferCatalog) GetByCode(context.Context, string) (*domain.SubscriptionPlan, error) {
	return nil, errOfferCatalogUnavailable
}

func TestCurrentSubscriptionDoesNotRepriceLegacyPlanDuringCatalogOutage(t *testing.T) {
	payments := newTestPaymentRepo()
	subscriptions := newTestSubscriptionRepo()
	payments.billingByTenant[1] = &domain.TenantBillingInfo{SubscriptionPlan: "monthly"}
	subscriptions.infoByTenant[1] = &domain.TenantSubscriptionInfo{TenantID: 1, SubscriptionStatus: domain.StatusActive, CanAccess: true}
	uc := NewPaymentUseCase(payments, subscriptions, unavailableOfferCatalog{}, &testPaymobClient{}, nil)
	result, err := uc.GetCurrentSubscription(context.Background(), 1)
	if !errors.Is(err, errOfferCatalogUnavailable) || result != nil {
		t.Fatalf("catalog outage must not display a future offer as a historical plan: result=%+v err=%v", result, err)
	}
}
