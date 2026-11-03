package repository

import (
	"fmt"
	"testing"
	"time"

	"pos-saas/internal/domain"
)

func TestHRChildRepositoriesRejectForeignEmployeesOnCreateAndUpdate(t *testing.T) {
	db, fixture := openDriverScopeFixture(t)
	ownEmployeeID := insertHREmployee(t, db, fixture.tenantA, fixture.restaurantA1)
	foreignEmployees := []struct {
		name string
		id   int
	}{
		{name: "sibling restaurant", id: insertHREmployee(t, db, fixture.tenantA, fixture.restaurantA2)},
		{name: "foreign tenant", id: insertHREmployee(t, db, fixture.tenantB, fixture.restaurantB1)},
	}
	unique := fmt.Sprintf("%d", time.Now().UnixNano())
	now := time.Now().UTC()
	type childResource struct {
		name   string
		table  string
		create func(employeeID int, suffix string) (int, error)
		update func(id, employeeID int) error
	}
	resources := []childResource{
		{
			name: "employment contracts", table: "employment_contracts",
			create: func(employeeID int, suffix string) (int, error) {
				return NewContractRepository(db).Create(&domain.EmploymentContract{
					TenantID: fixture.tenantA, RestaurantID: fixture.restaurantA1, EmployeeID: employeeID,
					ContractType: "permanent", ContractNumber: "HR-CONTRACT-" + unique + suffix,
					StartDate: now, Currency: "EGP", Status: "active",
				})
			},
			update: func(id, employeeID int) error {
				return NewContractRepository(db).Update(fixture.tenantA, fixture.restaurantA1, &domain.EmploymentContract{
					ID: id, TenantID: fixture.tenantA, RestaurantID: fixture.restaurantA1, EmployeeID: employeeID,
					ContractType: "permanent", ContractNumber: "HR-CONTRACT-ORIGINAL-" + unique,
					StartDate: now, Currency: "EGP", Status: "active",
				})
			},
		},
		{
			name: "training records", table: "training_records",
			create: func(employeeID int, _ string) (int, error) {
				return NewTrainingRepository(db).Create(&domain.TrainingRecord{
					TenantID: fixture.tenantA, RestaurantID: fixture.restaurantA1, EmployeeID: employeeID,
					TrainingTitle: "Scope test", TrainingType: "mandatory", StartDate: now, Status: "scheduled",
				})
			},
			update: func(id, employeeID int) error {
				return NewTrainingRepository(db).Update(fixture.tenantA, fixture.restaurantA1, &domain.TrainingRecord{
					ID: id, TenantID: fixture.tenantA, RestaurantID: fixture.restaurantA1, EmployeeID: employeeID,
					TrainingTitle: "Changed scope test", TrainingType: "mandatory", StartDate: now, Status: "scheduled",
				})
			},
		},
		{
			name: "insurance tracking", table: "insurance_tracking",
			create: func(employeeID int, suffix string) (int, error) {
				return NewInsuranceRepository(db).Create(&domain.InsuranceTracking{
					TenantID: fixture.tenantA, RestaurantID: fixture.restaurantA1, EmployeeID: employeeID,
					PolicyType: "medical", ProviderName: "Scope test", PolicyNumber: "HR-POLICY-" + unique + suffix,
					CoverageStart: now, CoverageEnd: now.AddDate(1, 0, 0), Status: "active",
				})
			},
			update: func(id, employeeID int) error {
				return NewInsuranceRepository(db).Update(fixture.tenantA, fixture.restaurantA1, &domain.InsuranceTracking{
					ID: id, TenantID: fixture.tenantA, RestaurantID: fixture.restaurantA1, EmployeeID: employeeID,
					PolicyType: "medical", ProviderName: "Changed", PolicyNumber: "HR-POLICY-ORIGINAL-" + unique,
					CoverageStart: now, CoverageEnd: now.AddDate(1, 0, 0), Status: "active",
				})
			},
		},
		{
			name: "resignations", table: "resignations",
			create: func(employeeID int, _ string) (int, error) {
				return NewResignationRepository(db).Create(&domain.Resignation{
					TenantID: fixture.tenantA, RestaurantID: fixture.restaurantA1, EmployeeID: employeeID,
					ResignationDate: now, LastWorkingDate: now.AddDate(0, 0, 30), Reason: "Scope test",
				})
			},
			update: func(id, employeeID int) error {
				return NewResignationRepository(db).Update(fixture.tenantA, fixture.restaurantA1, &domain.Resignation{
					ID: id, TenantID: fixture.tenantA, RestaurantID: fixture.restaurantA1, EmployeeID: employeeID,
					ResignationDate: now, LastWorkingDate: now.AddDate(0, 0, 30), Reason: "Changed scope test", Status: "submitted",
				})
			},
		},
	}

	for _, resource := range resources {
		t.Run(resource.name, func(t *testing.T) {
			ownID, err := resource.create(ownEmployeeID, "-original")
			if err != nil {
				t.Fatalf("create in-scope %s fixture: %v", resource.name, err)
			}
			for index, foreign := range foreignEmployees {
				t.Run(foreign.name, func(t *testing.T) {
					if _, err := resource.create(foreign.id, fmt.Sprintf("-foreign-%d", index)); err == nil {
						t.Fatalf("%s creation accepted an employee from a foreign scope", resource.name)
					}
					var createdCount int
					query := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE tenant_id = $1 AND restaurant_id = $2 AND employee_id = $3", resource.table)
					if err := db.QueryRow(query, fixture.tenantA, fixture.restaurantA1, foreign.id).Scan(&createdCount); err != nil {
						t.Fatalf("count rejected %s create rows: %v", resource.name, err)
					}
					if createdCount != 0 {
						t.Fatalf("rejected %s create left %d rows", resource.name, createdCount)
					}
					if err := resource.update(ownID, foreign.id); err == nil {
						t.Fatalf("%s update accepted an employee from a foreign scope", resource.name)
					}
					var storedEmployeeID int
					if err := db.QueryRow(fmt.Sprintf("SELECT employee_id FROM %s WHERE id = $1", resource.table), ownID).Scan(&storedEmployeeID); err != nil {
						t.Fatalf("read %s after rejected update: %v", resource.name, err)
					}
					if storedEmployeeID != ownEmployeeID {
						t.Fatalf("rejected %s update changed employee_id to %d", resource.name, storedEmployeeID)
					}
				})
			}
		})
	}
}
