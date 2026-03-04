package repository

import (
	"database/sql"
	"encoding/json"
	"errors"

	"pos-saas/internal/domain"
)

type AutoPostingRepository struct {
	db *sql.DB
}

func NewAutoPostingRepository(db *sql.DB) *AutoPostingRepository {
	return &AutoPostingRepository{db: db}
}

func (r *AutoPostingRepository) GetRuleByCode(tenantID, restaurantID int, code string) (*domain.AutoPostingRule, error) {
	rule, err := scanAutoPostingRule(r.db.QueryRow(`
		SELECT id, tenant_id, restaurant_id, rule_code, rule_name, COALESCE(description, ''),
		       is_active, default_accounts::text, created_at, updated_at
		FROM auto_posting_rules
		WHERE tenant_id = $1 AND restaurant_id = $2 AND rule_code = $3
	`, tenantID, restaurantID, code))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &rule, nil
}

func (r *AutoPostingRepository) ListRules(tenantID, restaurantID int) ([]domain.AutoPostingRule, error) {
	rows, err := r.db.Query(`
		SELECT id, tenant_id, restaurant_id, rule_code, rule_name, COALESCE(description, ''),
		       is_active, default_accounts::text, created_at, updated_at
		FROM auto_posting_rules
		WHERE tenant_id = $1 AND restaurant_id = $2
		ORDER BY rule_code ASC
	`, tenantID, restaurantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	rules := make([]domain.AutoPostingRule, 0)
	for rows.Next() {
		rule, err := scanAutoPostingRule(rows)
		if err != nil {
			return nil, err
		}
		rules = append(rules, rule)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return rules, nil
}

func (r *AutoPostingRepository) UpdateRule(rule *domain.AutoPostingRule) error {
	payload, err := json.Marshal(rule.DefaultAccounts)
	if err != nil {
		return err
	}

	result, err := r.db.Exec(`
		UPDATE auto_posting_rules
		SET rule_name = $1,
		    description = $2,
		    is_active = $3,
		    default_accounts = $4::jsonb,
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = $5 AND tenant_id = $6 AND restaurant_id = $7
	`,
		rule.RuleName,
		rule.Description,
		rule.IsActive,
		string(payload),
		rule.ID,
		rule.TenantID,
		rule.RestaurantID,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("auto posting rule not found")
	}

	return nil
}

func (r *AutoPostingRepository) GetDefaultAccounts(tenantID, restaurantID int, code string) (map[string]string, error) {
	defaults := cloneStringMap(domain.DefaultPostingRuleAccounts[code])

	rule, err := r.GetRuleByCode(tenantID, restaurantID, code)
	if err != nil {
		return nil, err
	}
	if rule == nil {
		return defaults, nil
	}

	for key, value := range rule.DefaultAccounts {
		defaults[key] = value
	}

	return defaults, nil
}

func (r *AutoPostingRepository) FindJournalEntryByReference(tenantID, restaurantID int, referenceType, entryType string, referenceID int) (*int, error) {
	var id int
	err := r.db.QueryRow(`
		SELECT id
		FROM journal_entries
		WHERE tenant_id = $1
		  AND restaurant_id = $2
		  AND reference_type = $3
		  AND reference_id = $4
		  AND entry_type = $5
		  AND status != $6
		ORDER BY id DESC
		LIMIT 1
	`, tenantID, restaurantID, referenceType, referenceID, entryType, domain.JournalStatusVoid).Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &id, nil
}

type autoPostingRuleScanner interface {
	Scan(dest ...interface{}) error
}

func scanAutoPostingRule(scanner autoPostingRuleScanner) (domain.AutoPostingRule, error) {
	var (
		rule            domain.AutoPostingRule
		defaultAccounts string
	)

	err := scanner.Scan(
		&rule.ID,
		&rule.TenantID,
		&rule.RestaurantID,
		&rule.RuleCode,
		&rule.RuleName,
		&rule.Description,
		&rule.IsActive,
		&defaultAccounts,
		&rule.CreatedAt,
		&rule.UpdatedAt,
	)
	if err != nil {
		return domain.AutoPostingRule{}, err
	}

	rule.DefaultAccounts = cloneStringMap(domain.DefaultPostingRuleAccounts[rule.RuleCode])
	if defaultAccounts != "" {
		var dbAccounts map[string]string
		if err := json.Unmarshal([]byte(defaultAccounts), &dbAccounts); err != nil {
			return domain.AutoPostingRule{}, err
		}
		for key, value := range dbAccounts {
			rule.DefaultAccounts[key] = value
		}
	}

	return rule, nil
}

func cloneStringMap(source map[string]string) map[string]string {
	if len(source) == 0 {
		return map[string]string{}
	}

	cloned := make(map[string]string, len(source))
	for key, value := range source {
		cloned[key] = value
	}
	return cloned
}
