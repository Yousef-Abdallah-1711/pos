package repository

import (
	"database/sql"
	"fmt"
	"testing"
	"time"

	"pos-saas/internal/domain"
)

func TestHRRepositoriesPostgresRejectForeignEmployeeReferences(t *testing.T) {
	db, fixture := openDriverScopeFixture(t)
	foreignEmployees := []struct {
		name                   string
		tenantID, restaurantID int
	}{
		{name: "sibling restaurant", tenantID: fixture.tenantA, restaurantID: fixture.restaurantA2},
		{name: "foreign tenant", tenantID: fixture.tenantB, restaurantID: fixture.restaurantB1},
	}

	for _, foreign := range foreignEmployees {
		t.Run(foreign.name, func(t *testing.T) {
			employeeID := insertHREmployee(t, db, foreign.tenantID, foreign.restaurantID)
			attendance := &domain.Attendance{
				TenantID: fixture.tenantA, RestaurantID: fixture.restaurantA1,
				EmployeeID: employeeID, AttendanceDate: time.Now(), Status: "present",
			}
			if _, err := NewAttendanceRepository(db).CreateAttendance(attendance); err == nil {
				t.Fatal("attendance creation accepted a foreign employee")
			}
			if _, err := NewAttendanceRepository(db).ClockIn(fixture.tenantA, fixture.restaurantA1, employeeID); err == nil {
				t.Fatal("clock-in accepted a foreign employee")
			}

			leave := &domain.Leave{
				TenantID: fixture.tenantA, RestaurantID: fixture.restaurantA1,
				EmployeeID: employeeID, StartDate: time.Now().AddDate(0, 1, 0),
				EndDate: time.Now().AddDate(0, 1, 1), TotalDays: 1,
				LeaveType: "annual", LeaveCategory: "paid", Status: "pending",
			}
			if _, err := NewLeaveRepository(db).CreateLeave(leave); err == nil {
				t.Fatal("leave creation accepted a foreign employee")
			}

			salary := &domain.Salary{
				TenantID: fixture.tenantA, RestaurantID: fixture.restaurantA1,
				EmployeeID: employeeID, PayPeriodStart: time.Now().AddDate(0, -1, 0),
				PayPeriodEnd: time.Now(), BaseSalary: 100, Currency: "EGP", Status: "pending",
			}
			if _, err := NewSalaryRepository(db).CreateSalary(salary); err == nil {
				t.Fatal("salary creation accepted a foreign employee")
			}
			assertNoHREmployeeSideEffects(t, db, employeeID)
		})
	}
}

func TestAttendanceRepositoryPostgresClockInDoesNotUpdateConflictingForeignScopeRow(t *testing.T) {
	db, fixture := openDriverScopeFixture(t)
	employeeID := insertHREmployee(t, db, fixture.tenantA, fixture.restaurantA2)
	var attendanceID int
	err := db.QueryRow(`
		INSERT INTO attendance (tenant_id, restaurant_id, employee_id, attendance_date, clock_in, status)
		VALUES ($1, $2, $3, CURRENT_DATE, CURRENT_TIMESTAMP - INTERVAL '2 hours', 'present')
		RETURNING id
	`, fixture.tenantA, fixture.restaurantA1, employeeID).Scan(&attendanceID)
	if err != nil {
		t.Fatalf("seed legacy mis-scoped attendance conflict row: %v", err)
	}
	var before time.Time
	if err := db.QueryRow(`SELECT clock_in FROM attendance WHERE id = $1`, attendanceID).Scan(&before); err != nil {
		t.Fatalf("read seeded clock-in: %v", err)
	}
	if _, err := NewAttendanceRepository(db).ClockIn(fixture.tenantA, fixture.restaurantA2, employeeID); err == nil {
		t.Fatal("clock-in unexpectedly updated an attendance row stored under another restaurant")
	}
	var after time.Time
	if err := db.QueryRow(`SELECT clock_in FROM attendance WHERE id = $1`, attendanceID).Scan(&after); err != nil {
		t.Fatalf("read clock-in after denied conflict update: %v", err)
	}
	if !after.Equal(before) {
		t.Fatalf("denied clock-in changed existing timestamp: before=%s after=%s", before, after)
	}
}

func TestEmployeeRepositoryPostgresRejectsForeignManagerReferences(t *testing.T) {
	db, fixture := openDriverScopeFixture(t)
	managerID := insertHREmployee(t, db, fixture.tenantA, fixture.restaurantA2)
	employeeID := insertHREmployee(t, db, fixture.tenantA, fixture.restaurantA1)
	var employeeCode string
	if err := db.QueryRow(`SELECT employee_code FROM employees WHERE id = $1`, employeeID).Scan(&employeeCode); err != nil {
		t.Fatalf("read target employee code: %v", err)
	}
	repo := NewEmployeeRepository(db)
	managerRef := managerID
	if _, err := repo.CreateEmployee(&domain.Employee{
		TenantID: fixture.tenantA, RestaurantID: fixture.restaurantA1,
		EmployeeCode: "HR-MGR-" + fmt.Sprint(time.Now().UnixNano())[:10],
		FirstName:    "Forged", LastName: "Manager", Email: fmt.Sprintf("hr-mgr-%d@example.test", time.Now().UnixNano()),
		HireDate: time.Now(), EmploymentType: "full_time", EmploymentStatus: "active",
		Position: "staff", ManagerID: &managerRef, BaseSalary: 0,
	}); err == nil {
		t.Fatal("employee creation accepted a manager from a sibling restaurant")
	}
	if err := repo.UpdateEmployee(&domain.Employee{
		ID: employeeID, TenantID: fixture.tenantA, RestaurantID: fixture.restaurantA1,
		FirstName: "Changed", LastName: "Name", Email: "changed@example.test",
		EmploymentType: "full_time", EmploymentStatus: "active", Position: "staff",
		ManagerID: &managerRef, WorkingHoursPerWeek: 40, IsActive: true,
	}); err == nil {
		t.Fatal("employee update accepted a manager from a sibling restaurant")
	}
	var firstName string
	var storedManager sql.NullInt64
	if err := db.QueryRow(`SELECT first_name, manager_id FROM employees WHERE id = $1`, employeeID).Scan(&firstName, &storedManager); err != nil {
		t.Fatalf("read employee after rejected manager update: %v", err)
	}
	if firstName != "Scope" || storedManager.Valid {
		t.Fatalf("rejected manager update changed target: first_name=%q manager=%v", firstName, storedManager)
	}
	var count int
	if err := db.QueryRow(`SELECT COUNT(*) FROM employees WHERE tenant_id = $1 AND restaurant_id = $2 AND employee_code LIKE 'HR-MGR-%'`, fixture.tenantA, fixture.restaurantA1).Scan(&count); err != nil {
		t.Fatalf("count rejected employee create: %v", err)
	}
	if count != 0 {
		t.Fatalf("foreign manager create left %d employees", count)
	}
	if employeeCode == "" {
		t.Fatal("fixture employee code unexpectedly empty")
	}
}

func insertHREmployee(t *testing.T, db *sql.DB, tenantID, restaurantID int) int {
	t.Helper()
	unique := fmt.Sprintf("%d", time.Now().UnixNano())
	var employeeID int
	err := db.QueryRow(`
		INSERT INTO employees (
			tenant_id, restaurant_id, employee_code, first_name, last_name, email,
			hire_date, employment_type, employment_status, position, base_salary
		) VALUES ($1, $2, $3, 'Scope', 'Test', $4, CURRENT_DATE, 'full_time', 'active', 'staff', 0)
		RETURNING id
	`, tenantID, restaurantID, "HR-SCOPE-"+unique, "hr-scope-"+unique+"@example.test").Scan(&employeeID)
	if err != nil {
		t.Fatalf("insert scoped test employee: %v", err)
	}
	return employeeID
}

func assertNoHREmployeeSideEffects(t *testing.T, db *sql.DB, employeeID int) {
	t.Helper()
	for _, table := range []string{"attendance", "leaves", "salaries"} {
		var count int
		query := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE employee_id = $1", table)
		if err := db.QueryRow(query, employeeID).Scan(&count); err != nil {
			t.Fatalf("count %s rows for rejected employee reference: %v", table, err)
		}
		if count != 0 {
			t.Errorf("rejected employee reference left %d rows in %s", count, table)
		}
	}
}
