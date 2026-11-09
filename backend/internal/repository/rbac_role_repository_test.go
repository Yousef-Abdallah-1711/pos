package repository

import (
	"context"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestCreateWithPermissionsRollsBackIncompleteRole(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("create sql mock: %v", err)
	}
	defer db.Close()
	setTenantContext := regexp.QuoteMeta("SELECT set_config($1, $2, true), set_config($3, $4, true), set_config($5, $6, true)")
	mock.ExpectBegin()
	mock.ExpectExec(setTenantContext).WithArgs("app.tenant_id", "7", "app.restaurant_id", "11", "app.user_id", "").WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectQuery("INSERT INTO roles").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(42))
	mock.ExpectExec("INSERT INTO role_permissions").WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectRollback()

	role := &RBACRole{TenantID: 7, RestaurantID: 11, RoleName: "cashier", RoleCode: "CASHIER"}
	if _, err := NewRBACRoleRepository(db).CreateWithPermissions(context.Background(), role, 3, []RBACPermissionMutation{{ModuleID: 999, Level: "WRITE"}}); err == nil {
		t.Fatal("invalid module should fail the role creation transaction")
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet transactional expectations: %v", err)
	}
}

func TestUpdateWithPermissionsRollsBackRoleAndPermissionChanges(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("create sql mock: %v", err)
	}
	defer db.Close()
	setTenantContext := regexp.QuoteMeta("SELECT set_config($1, $2, true), set_config($3, $4, true), set_config($5, $6, true)")
	mock.ExpectBegin()
	mock.ExpectExec(setTenantContext).WithArgs("app.tenant_id", "7", "app.restaurant_id", "11", "app.user_id", "").WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec("UPDATE roles SET role_name").WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec("DELETE FROM role_permissions").WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec("INSERT INTO role_permissions").WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectRollback()

	role := &RBACRole{ID: 42, TenantID: 7, RestaurantID: 11, RoleName: "updated cashier"}
	err = NewRBACRoleRepository(db).UpdateWithPermissions(context.Background(), role, 3, []RBACPermissionMutation{{ModuleID: 999, Level: "WRITE"}}, true)
	if err == nil {
		t.Fatal("invalid module should fail the role update transaction")
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet transactional expectations: %v", err)
	}
}

func TestUpdateWithPermissionsRejectsSystemRole(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("create sql mock: %v", err)
	}
	defer db.Close()
	setTenantContext := regexp.QuoteMeta("SELECT set_config($1, $2, true), set_config($3, $4, true), set_config($5, $6, true)")
	mock.ExpectBegin()
	mock.ExpectExec(setTenantContext).WithArgs("app.tenant_id", "7", "app.restaurant_id", "11", "app.user_id", "").WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec("UPDATE roles SET role_name").WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectRollback()

	role := &RBACRole{ID: 42, TenantID: 7, RestaurantID: 11, RoleName: "system mutation"}
	err = NewRBACRoleRepository(db).UpdateWithPermissions(context.Background(), role, 3, nil, false)
	if !errors.Is(err, ErrRBACRoleNotFound) {
		t.Fatalf("system role update error = %v, want role not found sentinel", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet transactional expectations: %v", err)
	}
}

func TestDeactivateWithPermissionsRollsBackRoleWhenGrantRemovalFails(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("create sql mock: %v", err)
	}
	defer db.Close()
	setTenantContext := regexp.QuoteMeta("SELECT set_config($1, $2, true), set_config($3, $4, true), set_config($5, $6, true)")
	mock.ExpectBegin()
	mock.ExpectExec(setTenantContext).WithArgs("app.tenant_id", "7", "app.restaurant_id", "11", "app.user_id", "").WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec("UPDATE roles SET is_active = FALSE").WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec("DELETE FROM role_permissions").WillReturnError(errors.New("database unavailable"))
	mock.ExpectRollback()

	err = NewRBACRoleRepository(db).DeactivateWithPermissions(context.Background(), 7, 11, 42, 3)
	if err == nil {
		t.Fatal("permission removal failure should roll back role deactivation")
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet transactional expectations: %v", err)
	}
}

func TestLegacyRoleMutationsUseTenantScopedTransactions(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("create sql mock: %v", err)
	}
	defer db.Close()
	setTenantContext := regexp.QuoteMeta("SELECT set_config($1, $2, true), set_config($3, $4, true), set_config($5, $6, true)")
	expectTenantScope := func() {
		t.Helper()
		mock.ExpectBegin()
		mock.ExpectExec(setTenantContext).WithArgs("app.tenant_id", "7", "app.restaurant_id", "11", "app.user_id", "").WillReturnResult(sqlmock.NewResult(0, 1))
	}

	expectTenantScope()
	mock.ExpectQuery("INSERT INTO roles").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(42))
	mock.ExpectCommit()

	expectTenantScope()
	mock.ExpectExec("UPDATE roles SET role_name").WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	expectTenantScope()
	mock.ExpectExec("UPDATE roles SET is_active = FALSE").WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec("DELETE FROM role_permissions").WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	repo := NewRBACRoleRepository(db)
	role := &RBACRole{ID: 42, TenantID: 7, RestaurantID: 11, RoleName: "cashier", RoleCode: "CASHIER"}
	if _, err := repo.Create(context.Background(), role, 3); err != nil {
		t.Fatalf("legacy create: %v", err)
	}
	if err := repo.Update(context.Background(), role, 3); err != nil {
		t.Fatalf("legacy update: %v", err)
	}
	if err := repo.Deactivate(context.Background(), 7, 11, 42, 3); err != nil {
		t.Fatalf("legacy deactivate: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet transactional expectations: %v", err)
	}
}
