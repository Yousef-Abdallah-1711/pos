package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"pos-saas/internal/domain"
)

// PlanRepository handles subscription plan persistence.
type PlanRepository interface {
	GetAllActive(ctx context.Context) ([]domain.SubscriptionPlan, error)
	GetByCode(ctx context.Context, code string) (*domain.SubscriptionPlan, error)
	Create(ctx context.Context, plan *domain.SubscriptionPlan) error
	Update(ctx context.Context, plan *domain.SubscriptionPlan) error
	Deactivate(ctx context.Context, code string) error
}

type planRepository struct {
	db             *sql.DB
	cacheMu        sync.RWMutex
	cachedPlans    []domain.SubscriptionPlan
	cacheExpiresAt time.Time
	cacheTTL       time.Duration
}

// NewPlanRepository creates a plan repository implementation.
func NewPlanRepository(db *sql.DB) PlanRepository {
	return &planRepository{
		db:       db,
		cacheTTL: 5 * time.Minute,
	}
}

func (r *planRepository) GetAllActive(ctx context.Context) ([]domain.SubscriptionPlan, error) {
	now := time.Now()

	r.cacheMu.RLock()
	if len(r.cachedPlans) > 0 && now.Before(r.cacheExpiresAt) {
		cached := clonePlans(r.cachedPlans)
		r.cacheMu.RUnlock()
		return cached, nil
	}
	r.cacheMu.RUnlock()

	rows, err := r.db.QueryContext(ctx, `
		SELECT
			id,
			name_en,
			name_ar,
			price_cents,
			currency,
			period_months,
			features::text,
			features_ar::text,
			is_active,
			display_order,
			created_at
		FROM subscription_plans
		WHERE is_active = true
		ORDER BY display_order ASC, id ASC
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to query active subscription plans: %w", err)
	}
	defer rows.Close()

	plans := make([]domain.SubscriptionPlan, 0)
	for rows.Next() {
		plan, err := scanSubscriptionPlan(rows)
		if err != nil {
			return nil, err
		}
		plans = append(plans, *plan)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate subscription plans: %w", err)
	}

	r.cacheMu.Lock()
	r.cachedPlans = clonePlans(plans)
	r.cacheExpiresAt = time.Now().Add(r.cacheTTL)
	r.cacheMu.Unlock()

	return clonePlans(plans), nil
}

func (r *planRepository) GetByCode(ctx context.Context, code string) (*domain.SubscriptionPlan, error) {
	trimmedCode := strings.TrimSpace(code)
	if trimmedCode == "" {
		return nil, fmt.Errorf("plan not found")
	}

	if plans, err := r.GetAllActive(ctx); err == nil {
		for _, plan := range plans {
			if strings.EqualFold(plan.ID, trimmedCode) {
				copy := plan
				return &copy, nil
			}
		}
	}

	row := r.db.QueryRowContext(ctx, `
		SELECT
			id,
			name_en,
			name_ar,
			price_cents,
			currency,
			period_months,
			features::text,
			features_ar::text,
			is_active,
			display_order,
			created_at
		FROM subscription_plans
		WHERE id = $1
	`, trimmedCode)

	plan, err := scanSubscriptionPlan(row)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("plan not found")
	}
	if err != nil {
		return nil, err
	}
	return plan, nil
}

func (r *planRepository) Create(ctx context.Context, plan *domain.SubscriptionPlan) error {
	if plan == nil {
		return fmt.Errorf("plan is required")
	}

	periodMonths := periodToMonths(plan.Period)
	featuresJSON, _ := json.Marshal(plan.Features)
	featuresArJSON, _ := json.Marshal(plan.FeaturesAr)

	_, err := r.db.ExecContext(ctx, `
		INSERT INTO subscription_plans (
			id, name_en, name_ar, price_cents, currency, period_months, features, features_ar, is_active, display_order
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7::jsonb, $8::jsonb, true, 0)
	`, strings.TrimSpace(plan.ID), plan.Name, plan.NameAr, plan.PriceCents, strings.ToUpper(plan.Currency), periodMonths, string(featuresJSON), string(featuresArJSON))
	if err != nil {
		return fmt.Errorf("failed to create subscription plan: %w", err)
	}

	r.invalidateCache()
	return nil
}

func (r *planRepository) Update(ctx context.Context, plan *domain.SubscriptionPlan) error {
	if plan == nil {
		return fmt.Errorf("plan is required")
	}

	periodMonths := periodToMonths(plan.Period)
	featuresJSON, _ := json.Marshal(plan.Features)
	featuresArJSON, _ := json.Marshal(plan.FeaturesAr)

	result, err := r.db.ExecContext(ctx, `
		UPDATE subscription_plans
		SET
			name_en = $1,
			name_ar = $2,
			price_cents = $3,
			currency = $4,
			period_months = $5,
			features = $6::jsonb,
			features_ar = $7::jsonb
		WHERE id = $8
	`, plan.Name, plan.NameAr, plan.PriceCents, strings.ToUpper(plan.Currency), periodMonths, string(featuresJSON), string(featuresArJSON), strings.TrimSpace(plan.ID))
	if err != nil {
		return fmt.Errorf("failed to update subscription plan: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to read affected rows: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("plan not found")
	}

	r.invalidateCache()
	return nil
}

func (r *planRepository) Deactivate(ctx context.Context, code string) error {
	result, err := r.db.ExecContext(ctx, `
		UPDATE subscription_plans
		SET is_active = false
		WHERE id = $1
	`, strings.TrimSpace(code))
	if err != nil {
		return fmt.Errorf("failed to deactivate plan: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to read affected rows: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("plan not found")
	}

	r.invalidateCache()
	return nil
}

func (r *planRepository) invalidateCache() {
	r.cacheMu.Lock()
	defer r.cacheMu.Unlock()

	r.cachedPlans = nil
	r.cacheExpiresAt = time.Time{}
}

type subscriptionPlanScanner interface {
	Scan(dest ...interface{}) error
}

func scanSubscriptionPlan(scanner subscriptionPlanScanner) (*domain.SubscriptionPlan, error) {
	var plan domain.SubscriptionPlan
	var periodMonths int
	var featuresRaw string
	var featuresArRaw string
	var isActive bool
	var displayOrder int
	var createdAt time.Time

	err := scanner.Scan(
		&plan.ID,
		&plan.Name,
		&plan.NameAr,
		&plan.PriceCents,
		&plan.Currency,
		&periodMonths,
		&featuresRaw,
		&featuresArRaw,
		&isActive,
		&displayOrder,
		&createdAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to scan subscription plan row: %w", err)
	}

	plan.Price = plan.PriceCents / 100
	plan.Period = monthsToPeriod(periodMonths)

	_ = json.Unmarshal([]byte(featuresRaw), &plan.Features)
	_ = json.Unmarshal([]byte(featuresArRaw), &plan.FeaturesAr)

	if plan.Features == nil {
		plan.Features = []string{}
	}
	if plan.FeaturesAr == nil {
		plan.FeaturesAr = []string{}
	}

	return &plan, nil
}

func monthsToPeriod(months int) string {
	switch months {
	case 12:
		return "year"
	case 1:
		return "month"
	default:
		if months <= 0 {
			return "month"
		}
		return fmt.Sprintf("%d-month", months)
	}
}

func periodToMonths(period string) int {
	normalized := strings.ToLower(strings.TrimSpace(period))
	switch normalized {
	case "year", "annual":
		return 12
	case "month", "monthly":
		return 1
	default:
		return 1
	}
}

func clonePlans(input []domain.SubscriptionPlan) []domain.SubscriptionPlan {
	if len(input) == 0 {
		return []domain.SubscriptionPlan{}
	}

	output := make([]domain.SubscriptionPlan, len(input))
	for idx, plan := range input {
		output[idx] = plan
		if plan.Features != nil {
			output[idx].Features = append([]string(nil), plan.Features...)
		}
		if plan.FeaturesAr != nil {
			output[idx].FeaturesAr = append([]string(nil), plan.FeaturesAr...)
		}
	}

	return output
}
