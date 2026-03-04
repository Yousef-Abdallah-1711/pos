package repository

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"pos-saas/internal/domain"
)

type PerformanceReviewRepository struct {
	db *sql.DB
}

func NewPerformanceReviewRepository(db *sql.DB) *PerformanceReviewRepository {
	return &PerformanceReviewRepository{db: db}
}

func (r *PerformanceReviewRepository) List(tenantID, restaurantID int) ([]domain.PerformanceReview, error) {
	items, _, err := r.ListPaginated(tenantID, restaurantID, 1000, 0)
	return items, err
}

func (r *PerformanceReviewRepository) ListPaginated(tenantID, restaurantID, limit, offset int) ([]domain.PerformanceReview, int, error) {
	if limit <= 0 {
		limit = 50
	}
	if limit > 200 {
		limit = 200
	}
	if offset < 0 {
		offset = 0
	}

	var total int
	if err := r.db.QueryRow(`
		SELECT COUNT(*)
		FROM performance_reviews
		WHERE tenant_id = $1 AND restaurant_id = $2
	`, tenantID, restaurantID).Scan(&total); err != nil {
		return nil, 0, err
	}

	rows, err := r.db.Query(`
		SELECT
			id, tenant_id, restaurant_id, employee_id, reviewer_id,
			review_period_start, review_period_end, review_date,
			overall_rating,
			COALESCE(productivity, 0) AS goals_achievement,
			COALESCE(technical_skills, 0) AS technical_skills,
			COALESCE(communication, 0) AS communication,
			COALESCE(leadership, 0) AS leadership,
			COALESCE(strengths, '') AS strengths,
			COALESCE(areas_for_improvement, '') AS areas_to_improve,
			COALESCE(action_plan, '') AS action_plan,
			status,
			acknowledged_at,
			NULL::INT AS acknowledged_by,
			finalized_at AS completed_at,
			created_at, updated_at, created_by, updated_by
		FROM performance_reviews
		WHERE tenant_id = $1 AND restaurant_id = $2
		ORDER BY review_date DESC, id DESC
		LIMIT $3 OFFSET $4
	`, tenantID, restaurantID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	out := make([]domain.PerformanceReview, 0)
	for rows.Next() {
		var item domain.PerformanceReview
		var acknowledgedBy, createdBy, updatedBy sql.NullInt64
		if err := rows.Scan(
			&item.ID, &item.TenantID, &item.RestaurantID, &item.EmployeeID, &item.ReviewerID,
			&item.ReviewPeriodStart, &item.ReviewPeriodEnd, &item.ReviewDate,
			&item.OverallRating, &item.GoalsAchievement, &item.TechnicalSkills, &item.Communication, &item.Leadership,
			&item.Strengths, &item.AreasToImprove, &item.ActionPlan, &item.Status,
			&item.AcknowledgedAt, &acknowledgedBy, &item.CompletedAt,
			&item.CreatedAt, &item.UpdatedAt, &createdBy, &updatedBy,
		); err != nil {
			return nil, 0, err
		}
		if acknowledgedBy.Valid {
			v := int(acknowledgedBy.Int64)
			item.AcknowledgedBy = &v
		}
		if createdBy.Valid {
			v := int(createdBy.Int64)
			item.CreatedBy = &v
		}
		if updatedBy.Valid {
			v := int(updatedBy.Int64)
			item.UpdatedBy = &v
		}
		out = append(out, item)
	}
	return out, total, rows.Err()
}

func (r *PerformanceReviewRepository) GetByID(tenantID, restaurantID, id int) (*domain.PerformanceReview, error) {
	var item domain.PerformanceReview
	var acknowledgedBy, createdBy, updatedBy sql.NullInt64
	err := r.db.QueryRow(`
		SELECT
			id, tenant_id, restaurant_id, employee_id, reviewer_id,
			review_period_start, review_period_end, review_date,
			overall_rating,
			COALESCE(productivity, 0) AS goals_achievement,
			COALESCE(technical_skills, 0) AS technical_skills,
			COALESCE(communication, 0) AS communication,
			COALESCE(leadership, 0) AS leadership,
			COALESCE(strengths, '') AS strengths,
			COALESCE(areas_for_improvement, '') AS areas_to_improve,
			COALESCE(action_plan, '') AS action_plan,
			status,
			acknowledged_at,
			NULL::INT AS acknowledged_by,
			finalized_at AS completed_at,
			created_at, updated_at, created_by, updated_by
		FROM performance_reviews
		WHERE tenant_id = $1 AND restaurant_id = $2 AND id = $3
	`, tenantID, restaurantID, id).Scan(
		&item.ID, &item.TenantID, &item.RestaurantID, &item.EmployeeID, &item.ReviewerID,
		&item.ReviewPeriodStart, &item.ReviewPeriodEnd, &item.ReviewDate,
		&item.OverallRating, &item.GoalsAchievement, &item.TechnicalSkills, &item.Communication, &item.Leadership,
		&item.Strengths, &item.AreasToImprove, &item.ActionPlan, &item.Status,
		&item.AcknowledgedAt, &acknowledgedBy, &item.CompletedAt,
		&item.CreatedAt, &item.UpdatedAt, &createdBy, &updatedBy,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if acknowledgedBy.Valid {
		v := int(acknowledgedBy.Int64)
		item.AcknowledgedBy = &v
	}
	if createdBy.Valid {
		v := int(createdBy.Int64)
		item.CreatedBy = &v
	}
	if updatedBy.Valid {
		v := int(updatedBy.Int64)
		item.UpdatedBy = &v
	}
	return &item, nil
}

func (r *PerformanceReviewRepository) Create(review *domain.PerformanceReview) (int, error) {
	status := normalizePerformanceStatus(review.Status)
	reviewerComments := strings.TrimSpace(review.ActionPlan)
	if reviewerComments == "" {
		reviewerComments = "N/A"
	}
	review.OverallRating = calculateWeightedRating(review.GoalsAchievement, review.TechnicalSkills, review.Communication, review.Leadership)
	var id int
	err := r.db.QueryRow(`
		INSERT INTO performance_reviews (
			tenant_id, restaurant_id, employee_id, reviewer_id,
			review_period_start, review_period_end, review_date,
			review_type,
			overall_rating, quality_of_work, productivity, technical_skills, communication, leadership,
			strengths, areas_for_improvement, action_plan, reviewer_comments, status, created_by
		)
		VALUES ($1,$2,$3,$4,$5,$6,$7,'annual',$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18)
		RETURNING id
	`,
		review.TenantID, review.RestaurantID, review.EmployeeID, review.ReviewerID,
		review.ReviewPeriodStart, review.ReviewPeriodEnd, review.ReviewDate,
		review.OverallRating, review.GoalsAchievement, review.GoalsAchievement, review.TechnicalSkills, review.Communication, review.Leadership,
		review.Strengths, review.AreasToImprove, review.ActionPlan, reviewerComments, status, review.CreatedBy,
	).Scan(&id)
	return id, err
}

func (r *PerformanceReviewRepository) Update(tenantID, restaurantID int, review *domain.PerformanceReview) error {
	status := normalizePerformanceStatus(review.Status)
	reviewerComments := strings.TrimSpace(review.ActionPlan)
	if reviewerComments == "" {
		reviewerComments = "N/A"
	}
	review.OverallRating = calculateWeightedRating(review.GoalsAchievement, review.TechnicalSkills, review.Communication, review.Leadership)
	res, err := r.db.Exec(`
		UPDATE performance_reviews
		SET review_period_start = $1,
		    review_period_end = $2,
		    review_date = $3,
		    overall_rating = $4,
		    quality_of_work = $5,
		    productivity = $6,
		    technical_skills = $7,
		    communication = $8,
		    leadership = $9,
		    strengths = $10,
		    areas_for_improvement = $11,
		    action_plan = $12,
		    reviewer_comments = $13,
		    status = $14,
		    updated_by = $15,
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = $16 AND tenant_id = $17 AND restaurant_id = $18
	`, review.ReviewPeriodStart, review.ReviewPeriodEnd, review.ReviewDate,
		review.OverallRating, review.GoalsAchievement, review.GoalsAchievement, review.TechnicalSkills, review.Communication, review.Leadership,
		review.Strengths, review.AreasToImprove, review.ActionPlan, reviewerComments, status, review.UpdatedBy,
		review.ID, tenantID, restaurantID,
	)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return errors.New("performance review not found")
	}
	return nil
}

func (r *PerformanceReviewRepository) Delete(tenantID, restaurantID, id int) error {
	res, err := r.db.Exec(`
		DELETE FROM performance_reviews
		WHERE id = $1 AND tenant_id = $2 AND restaurant_id = $3
	`, id, tenantID, restaurantID)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return errors.New("performance review not found")
	}
	return nil
}

func (r *PerformanceReviewRepository) Acknowledge(tenantID, restaurantID, id, acknowledgedBy int) error {
	res, err := r.db.Exec(`
		UPDATE performance_reviews
		SET status = 'acknowledged',
		    employee_acknowledged = true,
		    acknowledged_at = CURRENT_TIMESTAMP,
		    updated_by = $1,
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = $2 AND tenant_id = $3 AND restaurant_id = $4
	`, acknowledgedBy, id, tenantID, restaurantID)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return errors.New("performance review not found")
	}
	return nil
}

func (r *PerformanceReviewRepository) Complete(tenantID, restaurantID, id int) error {
	res, err := r.db.Exec(`
		UPDATE performance_reviews
		SET status = 'finalized',
		    is_finalized = true,
		    finalized_at = CURRENT_TIMESTAMP,
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = $1 AND tenant_id = $2 AND restaurant_id = $3
	`, id, tenantID, restaurantID)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return errors.New("performance review not found")
	}
	return nil
}

func (r *PerformanceReviewRepository) GetEmployeeHistory(tenantID, restaurantID, employeeID int) ([]domain.PerformanceReview, float64, error) {
	rows, err := r.db.Query(`
		SELECT
			id, tenant_id, restaurant_id, employee_id, reviewer_id,
			review_period_start, review_period_end, review_date,
			overall_rating,
			COALESCE(productivity, 0) AS goals_achievement,
			COALESCE(technical_skills, 0) AS technical_skills,
			COALESCE(communication, 0) AS communication,
			COALESCE(leadership, 0) AS leadership,
			COALESCE(strengths, '') AS strengths,
			COALESCE(areas_for_improvement, '') AS areas_to_improve,
			COALESCE(action_plan, '') AS action_plan,
			status,
			acknowledged_at,
			NULL::INT AS acknowledged_by,
			finalized_at AS completed_at,
			created_at, updated_at, created_by, updated_by
		FROM performance_reviews
		WHERE tenant_id = $1 AND restaurant_id = $2 AND employee_id = $3
		ORDER BY review_date DESC, id DESC
	`, tenantID, restaurantID, employeeID)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	out := make([]domain.PerformanceReview, 0)
	var sum float64
	for rows.Next() {
		var item domain.PerformanceReview
		var acknowledgedBy, createdBy, updatedBy sql.NullInt64
		if err := rows.Scan(
			&item.ID, &item.TenantID, &item.RestaurantID, &item.EmployeeID, &item.ReviewerID,
			&item.ReviewPeriodStart, &item.ReviewPeriodEnd, &item.ReviewDate,
			&item.OverallRating, &item.GoalsAchievement, &item.TechnicalSkills, &item.Communication, &item.Leadership,
			&item.Strengths, &item.AreasToImprove, &item.ActionPlan, &item.Status,
			&item.AcknowledgedAt, &acknowledgedBy, &item.CompletedAt,
			&item.CreatedAt, &item.UpdatedAt, &createdBy, &updatedBy,
		); err != nil {
			return nil, 0, err
		}
		if acknowledgedBy.Valid {
			v := int(acknowledgedBy.Int64)
			item.AcknowledgedBy = &v
		}
		sum += item.OverallRating
		out = append(out, item)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}
	avg := 0.0
	if len(out) > 0 {
		avg = sum / float64(len(out))
	}
	return out, avg, nil
}

func normalizePerformanceStatus(status string) string {
	switch strings.ToLower(strings.TrimSpace(status)) {
	case "finalized", "acknowledged", "submitted", "draft", "disputed":
		return strings.ToLower(strings.TrimSpace(status))
	case "completed":
		return "finalized"
	default:
		return "draft"
	}
}

func calculateWeightedRating(goals, technical, communication, leadership float64) float64 {
	// 40% goals, 30% technical, 20% communication, 10% leadership
	return (goals * 0.40) + (technical * 0.30) + (communication * 0.20) + (leadership * 0.10)
}

func parseDateOrDefault(raw string) time.Time {
	if raw == "" {
		return time.Now().UTC()
	}
	parsed, err := time.Parse("2006-01-02", raw)
	if err != nil {
		return time.Now().UTC()
	}
	return parsed.UTC()
}
