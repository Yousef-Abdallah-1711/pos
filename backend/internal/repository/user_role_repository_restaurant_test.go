package repository

import (
	"context"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestGetUserRolesForRestaurantScopesTenantRestaurantAndUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("create SQL mock: %v", err)
	}
	defer db.Close()

	query := `(?s)SELECT DISTINCT r\.id.*FROM user_roles ur.*JOIN users u ON u\.id = ur\.user_id AND u\.tenant_id = ur\.tenant_id.*JOIN roles r ON r\.id = ur\.role_id AND r\.tenant_id = ur\.tenant_id.*ur\.user_id = \$1 AND ur\.tenant_id = \$2 AND r\.restaurant_id = \$3.*u\.restaurant_id IS NULL OR u\.restaurant_id = \$3`
	if _, err := regexp.Compile(query); err != nil {
		t.Fatalf("invalid SQL matcher: %v", err)
	}
	mock.ExpectQuery(query).
		WithArgs(int64(11), int64(22), int64(33)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "tenant_id", "role_name", "description", "role_code", "is_system_role", "created_at", "updated_at"}))

	roles, err := NewUserRoleRepository(db).GetUserRolesForRestaurant(context.Background(), 11, 22, 33)
	if err != nil {
		t.Fatalf("GetUserRolesForRestaurant returned error: %v", err)
	}
	if len(roles) != 0 {
		t.Fatalf("expected no roles, got %d", len(roles))
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet SQL expectations: %v", err)
	}
}
