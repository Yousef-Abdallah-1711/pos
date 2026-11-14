package usecase

import (
	"context"
	"fmt"
	"strings"

	"pos-saas/internal/commerce"
	"pos-saas/internal/domain"
)

// currentSubscriptionPlan never substitutes a future price for a known paid plan.
func (uc *PaymentUseCase) currentSubscriptionPlan(ctx context.Context, planID string, plans map[string]domain.SubscriptionPlan) (domain.SubscriptionPlan, error) {
	if plan, ok := plans[planID]; ok {
		return plan, nil
	}
	if strings.EqualFold(planID, "free") || uc.planRepo == nil {
		return uc.defaultPlan(), nil
	}
	historical, err := uc.planRepo.GetByCode(ctx, planID)
	if err != nil {
		return domain.SubscriptionPlan{}, fmt.Errorf("load historical subscription plan: %w", err)
	}
	if historical == nil {
		return domain.SubscriptionPlan{}, fmt.Errorf("historical subscription plan is unavailable")
	}
	return *historical, nil
}

func (uc *PaymentUseCase) defaultPlan() domain.SubscriptionPlan {
	monthlyPlan, ok := uc.plans["monthly"]
	if ok {
		return monthlyPlan
	}
	for _, plan := range uc.plans {
		return plan
	}
	return buildDefaultPlans()["monthly"]
}

func buildDefaultPlans() map[string]domain.SubscriptionPlan {
	offer := commerce.Current()
	features := []string{"Cashier POS", "Kitchen display", "Cash ordering", "Theme configuration"}
	featuresAr := []string{"نقاط البيع", "شاشة المطبخ", "الطلبات النقدية", "إعداد القالب"}
	return map[string]domain.SubscriptionPlan{
		"monthly": {ID: offer.MonthlyPlanID, Name: "SERVIO Monthly", NameAr: "سيرفيو شهري", Price: offer.MonthlyPriceCents / 100, PriceCents: offer.MonthlyPriceCents, Currency: offer.Currency, Period: "month", Features: features, FeaturesAr: featuresAr},
		"annual":  {ID: offer.AnnualPlanID, Name: "SERVIO Annual", NameAr: "سيرفيو سنوي", Price: offer.AnnualPriceCents / 100, PriceCents: offer.AnnualPriceCents, Currency: offer.Currency, Period: "year", Features: features, FeaturesAr: featuresAr},
	}
}
