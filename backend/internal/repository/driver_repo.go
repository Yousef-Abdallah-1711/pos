package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"pos-saas/internal/domain"
)

var ErrOrderAlreadyAssigned = errors.New("order already assigned to a driver")

// DriverRepository defines operations for driver data access
type DriverRepository interface {
	// Driver Management
	Create(ctx context.Context, driver *domain.Driver) error
	GetByID(ctx context.Context, id, tenantID, restaurantID int) (*domain.Driver, error)
	GetByEmail(ctx context.Context, email string, tenantID, restaurantID int) (*domain.Driver, error)
	GetByTenantAndRestaurant(ctx context.Context, tenantID, restaurantID int) ([]domain.Driver, error)
	Update(ctx context.Context, driver *domain.Driver, tenantID, restaurantID int) error
	Delete(ctx context.Context, id, tenantID, restaurantID int) error

	// Status & Availability
	UpdateStatus(ctx context.Context, id, tenantID, restaurantID int, status string) error
	UpdateAvailability(ctx context.Context, id, tenantID, restaurantID int, status string) error
	GetAvailableDrivers(ctx context.Context, tenantID, restaurantID int) ([]domain.Driver, error)
	UpdateActiveOrders(ctx context.Context, id, tenantID, restaurantID, count int) error

	// Location
	UpdateLocation(ctx context.Context, driverID, tenantID, restaurantID int, lat, lon float64) error
	GetLastLocation(ctx context.Context, driverID, tenantID, restaurantID int) (*domain.DriverLocation, error)
	RecordLocation(ctx context.Context, location *domain.DriverLocation, tenantID, restaurantID int) error
	GetLocationHistory(ctx context.Context, driverID, tenantID, restaurantID, limit int) ([]domain.DriverLocation, error)

	// Assignments
	AssignOrder(ctx context.Context, assignment *domain.DriverAssignment, tenantID, restaurantID int) error
	GetAssignmentByID(ctx context.Context, id, tenantID, restaurantID int) (*domain.DriverAssignment, error)
	GetActiveAssignment(ctx context.Context, driverID, tenantID, restaurantID int) (*domain.DriverAssignment, error)
	GetAssignmentByOrderID(ctx context.Context, orderID, tenantID, restaurantID int) (*domain.DriverAssignment, error)
	UpdateAssignmentStatus(ctx context.Context, assignmentID, tenantID, restaurantID int, status string) error
	RateDriver(ctx context.Context, assignmentID, tenantID, restaurantID, rating int, comment string) error
	GetAssignmentHistory(ctx context.Context, driverID, tenantID, restaurantID, limit int) ([]domain.DriverAssignment, error)
	CompleteAssignment(ctx context.Context, assignmentID, tenantID, restaurantID int) error

	// Statistics
	GetDriverStats(ctx context.Context, driverID, tenantID, restaurantID int) (*domain.DriverStats, error)
	UpdateDeliveryMetrics(ctx context.Context, driverID, tenantID, restaurantID int, completed bool) error
}

type driverRepository struct {
	db *sql.DB
}

// NewDriverRepository creates a new driver repository
func NewDriverRepository(db *sql.DB) DriverRepository {
	return &driverRepository{db: db}
}

// Create creates a new driver
func (r *driverRepository) Create(ctx context.Context, driver *domain.Driver) error {
	query := `
		INSERT INTO drivers (
			tenant_id, restaurant_id, first_name, last_name, email, phone_number,
			license_number, vehicle_type, vehicle_number, status, availability_status
		) SELECT $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
		WHERE EXISTS (SELECT 1 FROM restaurants WHERE id = $2 AND tenant_id = $1)
		RETURNING id, created_at, updated_at
	`

	return r.db.QueryRowContext(
		ctx, query,
		driver.TenantID, driver.RestaurantID, driver.FirstName, driver.LastName,
		driver.Email, driver.PhoneNumber, driver.LicenseNumber, driver.VehicleType,
		driver.VehicleNumber, driver.Status, driver.AvailabilityStatus,
	).Scan(&driver.ID, &driver.CreatedAt, &driver.UpdatedAt)
}

// GetByID retrieves a driver by ID
func (r *driverRepository) GetByID(ctx context.Context, id, tenantID, restaurantID int) (*domain.Driver, error) {
	driver := &domain.Driver{}
	query := `
		SELECT id, tenant_id, restaurant_id, first_name, last_name, email, phone_number,
		       license_number, COALESCE(vehicle_type, ''), COALESCE(vehicle_number, ''),
		       status, availability_status,
		       current_latitude, current_longitude, total_deliveries, completed_deliveries,
		       cancelled_deliveries, average_rating, active_orders_count, date_of_birth,
		       COALESCE(address, ''), COALESCE(city, ''), COALESCE(state, ''),
		       COALESCE(zip_code, ''), is_verified, verification_date, joined_date,
		       last_active, COALESCE(notes, ''), created_at, updated_at
		FROM drivers WHERE id = $1 AND tenant_id = $2 AND restaurant_id = $3
	`

	err := r.db.QueryRowContext(ctx, query, id, tenantID, restaurantID).Scan(
		&driver.ID, &driver.TenantID, &driver.RestaurantID, &driver.FirstName, &driver.LastName,
		&driver.Email, &driver.PhoneNumber, &driver.LicenseNumber, &driver.VehicleType,
		&driver.VehicleNumber, &driver.Status, &driver.AvailabilityStatus,
		&driver.CurrentLatitude, &driver.CurrentLongitude, &driver.TotalDeliveries,
		&driver.CompletedDeliveries, &driver.CancelledDeliveries, &driver.AverageRating,
		&driver.ActiveOrdersCount, &driver.DateOfBirth, &driver.Address, &driver.City,
		&driver.State, &driver.ZipCode, &driver.IsVerified, &driver.VerificationDate,
		&driver.JoinedDate, &driver.LastActive, &driver.Notes, &driver.CreatedAt, &driver.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("driver not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get driver: %w", err)
	}

	return driver, nil
}

// GetByEmail retrieves a driver by email
func (r *driverRepository) GetByEmail(ctx context.Context, email string, tenantID, restaurantID int) (*domain.Driver, error) {
	driver := &domain.Driver{}
	query := `
		SELECT id, tenant_id, restaurant_id, first_name, last_name, email, phone_number,
		       license_number, COALESCE(vehicle_type, ''), COALESCE(vehicle_number, ''),
		       status, availability_status,
		       current_latitude, current_longitude, total_deliveries, completed_deliveries,
		       cancelled_deliveries, average_rating, active_orders_count, date_of_birth,
		       COALESCE(address, ''), COALESCE(city, ''), COALESCE(state, ''),
		       COALESCE(zip_code, ''), is_verified, verification_date, joined_date,
		       last_active, COALESCE(notes, ''), created_at, updated_at
		FROM drivers WHERE email = $1 AND tenant_id = $2 AND restaurant_id = $3 LIMIT 1
	`

	err := r.db.QueryRowContext(ctx, query, email, tenantID, restaurantID).Scan(
		&driver.ID, &driver.TenantID, &driver.RestaurantID, &driver.FirstName, &driver.LastName,
		&driver.Email, &driver.PhoneNumber, &driver.LicenseNumber, &driver.VehicleType,
		&driver.VehicleNumber, &driver.Status, &driver.AvailabilityStatus,
		&driver.CurrentLatitude, &driver.CurrentLongitude, &driver.TotalDeliveries,
		&driver.CompletedDeliveries, &driver.CancelledDeliveries, &driver.AverageRating,
		&driver.ActiveOrdersCount, &driver.DateOfBirth, &driver.Address, &driver.City,
		&driver.State, &driver.ZipCode, &driver.IsVerified, &driver.VerificationDate,
		&driver.JoinedDate, &driver.LastActive, &driver.Notes, &driver.CreatedAt, &driver.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("driver not found")
	}
	return driver, err
}

// GetByTenantAndRestaurant retrieves all drivers for a restaurant
func (r *driverRepository) GetByTenantAndRestaurant(ctx context.Context, tenantID, restaurantID int) ([]domain.Driver, error) {
	query := `
		SELECT id, tenant_id, restaurant_id, first_name, last_name, email, phone_number,
		       license_number, COALESCE(vehicle_type, ''), COALESCE(vehicle_number, ''),
		       status, availability_status,
		       current_latitude, current_longitude, total_deliveries, completed_deliveries,
		       cancelled_deliveries, average_rating, active_orders_count, date_of_birth,
		       COALESCE(address, ''), COALESCE(city, ''), COALESCE(state, ''),
		       COALESCE(zip_code, ''), is_verified, verification_date, joined_date,
		       last_active, COALESCE(notes, ''), created_at, updated_at
		FROM drivers
		WHERE tenant_id = $1 AND restaurant_id = $2
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, tenantID, restaurantID)
	if err != nil {
		return nil, fmt.Errorf("failed to query drivers: %w", err)
	}
	defer rows.Close()

	var drivers []domain.Driver
	for rows.Next() {
		driver := domain.Driver{}
		err := rows.Scan(
			&driver.ID, &driver.TenantID, &driver.RestaurantID, &driver.FirstName, &driver.LastName,
			&driver.Email, &driver.PhoneNumber, &driver.LicenseNumber, &driver.VehicleType,
			&driver.VehicleNumber, &driver.Status, &driver.AvailabilityStatus,
			&driver.CurrentLatitude, &driver.CurrentLongitude, &driver.TotalDeliveries,
			&driver.CompletedDeliveries, &driver.CancelledDeliveries, &driver.AverageRating,
			&driver.ActiveOrdersCount, &driver.DateOfBirth, &driver.Address, &driver.City,
			&driver.State, &driver.ZipCode, &driver.IsVerified, &driver.VerificationDate,
			&driver.JoinedDate, &driver.LastActive, &driver.Notes, &driver.CreatedAt, &driver.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan driver: %w", err)
		}
		drivers = append(drivers, driver)
	}

	return drivers, rows.Err()
}

// Update updates an existing driver
func (r *driverRepository) Update(ctx context.Context, driver *domain.Driver, tenantID, restaurantID int) error {
	query := `
		UPDATE drivers SET
			first_name = $1, last_name = $2, email = $3, phone_number = $4,
			license_number = $5, vehicle_type = $6, vehicle_number = $7,
			status = $8, availability_status = $9, date_of_birth = $10,
			address = $11, city = $12, state = $13, zip_code = $14,
			notes = $15, updated_at = NOW()
		WHERE id = $16 AND tenant_id = $17 AND restaurant_id = $18
	`

	result, err := r.db.ExecContext(
		ctx, query,
		driver.FirstName, driver.LastName, driver.Email, driver.PhoneNumber,
		driver.LicenseNumber, driver.VehicleType, driver.VehicleNumber,
		driver.Status, driver.AvailabilityStatus, driver.DateOfBirth,
		driver.Address, driver.City, driver.State, driver.ZipCode,
		driver.Notes, driver.ID, tenantID, restaurantID,
	)

	if err != nil {
		return fmt.Errorf("failed to update driver: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rows == 0 {
		return errors.New("driver not found")
	}

	return nil
}

// Delete deletes a driver
func (r *driverRepository) Delete(ctx context.Context, id, tenantID, restaurantID int) error {
	query := `DELETE FROM drivers WHERE id = $1 AND tenant_id = $2 AND restaurant_id = $3`

	result, err := r.db.ExecContext(ctx, query, id, tenantID, restaurantID)
	if err != nil {
		return fmt.Errorf("failed to delete driver: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rows == 0 {
		return errors.New("driver not found")
	}

	return nil
}

// UpdateStatus updates driver status
func (r *driverRepository) UpdateStatus(ctx context.Context, id, tenantID, restaurantID int, status string) error {
	query := `UPDATE drivers SET status = $1, updated_at = NOW() WHERE id = $2 AND tenant_id = $3 AND restaurant_id = $4`
	result, err := r.db.ExecContext(ctx, query, status, id, tenantID, restaurantID)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("get driver status affected rows: %w", err)
	}
	if rows == 0 {
		return errors.New("driver not found")
	}
	return nil
}

// UpdateAvailability updates driver availability status
func (r *driverRepository) UpdateAvailability(ctx context.Context, id, tenantID, restaurantID int, status string) error {
	query := `UPDATE drivers SET availability_status = $1, last_active = NOW(), updated_at = NOW() WHERE id = $2 AND tenant_id = $3 AND restaurant_id = $4`
	_, err := r.db.ExecContext(ctx, query, status, id, tenantID, restaurantID)
	return err
}

// GetAvailableDrivers retrieves all available drivers
func (r *driverRepository) GetAvailableDrivers(ctx context.Context, tenantID, restaurantID int) ([]domain.Driver, error) {
	query := `
		SELECT id, tenant_id, restaurant_id, first_name, last_name, email, phone_number,
		       license_number, COALESCE(vehicle_type, ''), COALESCE(vehicle_number, ''),
		       status, availability_status,
		       current_latitude, current_longitude, total_deliveries, completed_deliveries,
		       cancelled_deliveries, average_rating, active_orders_count, date_of_birth,
		       COALESCE(address, ''), COALESCE(city, ''), COALESCE(state, ''),
		       COALESCE(zip_code, ''), is_verified, verification_date, joined_date,
		       last_active, COALESCE(notes, ''), created_at, updated_at
		FROM drivers
		WHERE tenant_id = $1 AND restaurant_id = $2 AND status = 'active' AND availability_status = 'available'
		ORDER BY active_orders_count ASC, average_rating DESC
	`

	rows, err := r.db.QueryContext(ctx, query, tenantID, restaurantID)
	if err != nil {
		return nil, fmt.Errorf("failed to query available drivers: %w", err)
	}
	defer rows.Close()

	var drivers []domain.Driver
	for rows.Next() {
		driver := domain.Driver{}
		err := rows.Scan(
			&driver.ID, &driver.TenantID, &driver.RestaurantID, &driver.FirstName, &driver.LastName,
			&driver.Email, &driver.PhoneNumber, &driver.LicenseNumber, &driver.VehicleType,
			&driver.VehicleNumber, &driver.Status, &driver.AvailabilityStatus,
			&driver.CurrentLatitude, &driver.CurrentLongitude, &driver.TotalDeliveries,
			&driver.CompletedDeliveries, &driver.CancelledDeliveries, &driver.AverageRating,
			&driver.ActiveOrdersCount, &driver.DateOfBirth, &driver.Address, &driver.City,
			&driver.State, &driver.ZipCode, &driver.IsVerified, &driver.VerificationDate,
			&driver.JoinedDate, &driver.LastActive, &driver.Notes, &driver.CreatedAt, &driver.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan driver: %w", err)
		}
		drivers = append(drivers, driver)
	}

	return drivers, rows.Err()
}

// UpdateActiveOrders updates the count of active orders for a driver
func (r *driverRepository) UpdateActiveOrders(ctx context.Context, id, tenantID, restaurantID, count int) error {
	query := `UPDATE drivers SET active_orders_count = $1, updated_at = NOW() WHERE id = $2 AND tenant_id = $3 AND restaurant_id = $4`
	_, err := r.db.ExecContext(ctx, query, count, id, tenantID, restaurantID)
	return err
}

// UpdateLocation updates driver current location
func (r *driverRepository) UpdateLocation(ctx context.Context, driverID, tenantID, restaurantID int, lat, lon float64) error {
	query := `UPDATE drivers SET current_latitude = $1, current_longitude = $2, last_active = NOW(), updated_at = NOW() WHERE id = $3 AND tenant_id = $4 AND restaurant_id = $5`
	_, err := r.db.ExecContext(ctx, query, lat, lon, driverID, tenantID, restaurantID)
	return err
}

// GetLastLocation retrieves the last known location of a driver
func (r *driverRepository) GetLastLocation(ctx context.Context, driverID, tenantID, restaurantID int) (*domain.DriverLocation, error) {
	location := &domain.DriverLocation{}
	query := `
		SELECT l.id, l.driver_id, l.order_id, l.latitude, l.longitude, l.accuracy, l.recorded_at
		FROM driver_location_history l
		JOIN drivers d ON d.id = l.driver_id
		WHERE l.driver_id = $1 AND d.tenant_id = $2 AND d.restaurant_id = $3
		  AND (l.order_id IS NULL OR EXISTS (
			SELECT 1 FROM orders o WHERE o.id = l.order_id AND o.tenant_id = $2 AND o.restaurant_id = $3
		  ))
		ORDER BY l.recorded_at DESC
		LIMIT 1
	`

	err := r.db.QueryRowContext(ctx, query, driverID, tenantID, restaurantID).Scan(
		&location.ID, &location.DriverID, &location.OrderID, &location.Latitude, &location.Longitude, &location.Accuracy, &location.RecordedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("no location found")
	}

	return location, err
}

// RecordLocation records a new location entry
func (r *driverRepository) RecordLocation(ctx context.Context, location *domain.DriverLocation, tenantID, restaurantID int) error {
	query := `
		INSERT INTO driver_location_history (driver_id, order_id, latitude, longitude, accuracy)
		SELECT $1, $2, $3, $4, $5
		FROM drivers d
		WHERE d.id = $1 AND d.tenant_id = $6 AND d.restaurant_id = $7
		  AND ($2::INT IS NULL OR EXISTS (
			SELECT 1 FROM orders o WHERE o.id = $2 AND o.tenant_id = $6 AND o.restaurant_id = $7
		  ))
		RETURNING id, recorded_at
	`

	err := r.db.QueryRowContext(ctx, query, location.DriverID, location.OrderID, location.Latitude, location.Longitude, location.Accuracy, tenantID, restaurantID).Scan(
		&location.ID, &location.RecordedAt,
	)
	if err == sql.ErrNoRows {
		return errors.New("driver or order not found in tenant restaurant")
	}
	return err
}

// GetLocationHistory retrieves location history for a driver
func (r *driverRepository) GetLocationHistory(ctx context.Context, driverID, tenantID, restaurantID, limit int) ([]domain.DriverLocation, error) {
	query := `
		SELECT l.id, l.driver_id, l.order_id, l.latitude, l.longitude, l.accuracy, l.recorded_at
		FROM driver_location_history l
		JOIN drivers d ON d.id = l.driver_id
		WHERE l.driver_id = $1 AND d.tenant_id = $2 AND d.restaurant_id = $3
		  AND (l.order_id IS NULL OR EXISTS (
			SELECT 1 FROM orders o WHERE o.id = l.order_id AND o.tenant_id = $2 AND o.restaurant_id = $3
		  ))
		ORDER BY l.recorded_at DESC
		LIMIT $4
	`

	rows, err := r.db.QueryContext(ctx, query, driverID, tenantID, restaurantID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query location history: %w", err)
	}
	defer rows.Close()

	var locations []domain.DriverLocation
	for rows.Next() {
		loc := domain.DriverLocation{}
		err := rows.Scan(
			&loc.ID, &loc.DriverID, &loc.OrderID, &loc.Latitude, &loc.Longitude, &loc.Accuracy, &loc.RecordedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan location: %w", err)
		}
		locations = append(locations, loc)
	}

	return locations, rows.Err()
}

// AssignOrder creates a new order assignment
func (r *driverRepository) AssignOrder(ctx context.Context, assignment *domain.DriverAssignment, tenantID, restaurantID int) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin driver assignment: %w", err)
	}
	defer tx.Rollback()

	// Lock the order so concurrent requests cannot assign it to different drivers.
	var orderID int
	err = tx.QueryRowContext(ctx, `
		SELECT o.id
		FROM orders o
		JOIN drivers d ON d.id = $4
		WHERE o.id = $1 AND o.tenant_id = $2 AND o.restaurant_id = $3
		  AND d.tenant_id = o.tenant_id AND d.restaurant_id = o.restaurant_id
		FOR UPDATE OF o, d
	`, assignment.OrderID, tenantID, restaurantID, assignment.DriverID).Scan(&orderID)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("order or driver not found in tenant restaurant")
		}
		return fmt.Errorf("verify driver assignment ownership: %w", err)
	}

	var existing bool
	if err := tx.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM driver_assignments WHERE order_id = $1)`, orderID).Scan(&existing); err != nil {
		return fmt.Errorf("check existing driver assignment: %w", err)
	}
	if existing {
		return ErrOrderAlreadyAssigned
	}

	if err := tx.QueryRowContext(ctx, `
		INSERT INTO driver_assignments (order_id, driver_id, assigned_at, status, assignment_notes)
		VALUES ($1, $2, NOW(), $3, $4)
		RETURNING id, created_at, updated_at
	`, orderID, assignment.DriverID, assignment.Status, assignment.AssignmentNotes).Scan(
		&assignment.ID, &assignment.CreatedAt, &assignment.UpdatedAt,
	); err != nil {
		return fmt.Errorf("insert driver assignment: %w", err)
	}

	if _, err := tx.ExecContext(ctx, `
		UPDATE drivers
		SET availability_status = 'busy', active_orders_count = active_orders_count + 1,
		    last_active = NOW(), updated_at = NOW()
		WHERE id = $1 AND tenant_id = $2 AND restaurant_id = $3
	`, assignment.DriverID, tenantID, restaurantID); err != nil {
		return fmt.Errorf("update driver assignment state: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit driver assignment: %w", err)
	}
	return nil
}

// GetAssignmentByID retrieves assignment by ID
func (r *driverRepository) GetAssignmentByID(ctx context.Context, id, tenantID, restaurantID int) (*domain.DriverAssignment, error) {
	assignment := &domain.DriverAssignment{}
	query := `
		SELECT da.id, da.order_id, da.driver_id, da.assigned_at, da.accepted_at, da.started_at, da.completed_at,
		       da.status, da.assignment_notes, da.rating, da.rating_comment, da.created_at, da.updated_at
		FROM driver_assignments da
		JOIN drivers d ON d.id = da.driver_id
		JOIN orders o ON o.id = da.order_id
		WHERE da.id = $1 AND d.tenant_id = $2 AND d.restaurant_id = $3
		  AND o.tenant_id = $2 AND o.restaurant_id = $3
	`

	err := r.db.QueryRowContext(ctx, query, id, tenantID, restaurantID).Scan(
		&assignment.ID, &assignment.OrderID, &assignment.DriverID, &assignment.AssignedAt,
		&assignment.AcceptedAt, &assignment.StartedAt, &assignment.CompletedAt,
		&assignment.Status, &assignment.AssignmentNotes, &assignment.Rating,
		&assignment.RatingComment, &assignment.CreatedAt, &assignment.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("assignment not found")
	}

	return assignment, err
}

// GetActiveAssignment retrieves the active assignment for a driver
func (r *driverRepository) GetActiveAssignment(ctx context.Context, driverID, tenantID, restaurantID int) (*domain.DriverAssignment, error) {
	assignment := &domain.DriverAssignment{}
	query := `
		SELECT da.id, da.order_id, da.driver_id, da.assigned_at, da.accepted_at, da.started_at, da.completed_at,
		       da.status, da.assignment_notes, da.rating, da.rating_comment, da.created_at, da.updated_at
		FROM driver_assignments da
		JOIN drivers d ON d.id = da.driver_id
		JOIN orders o ON o.id = da.order_id
		WHERE da.driver_id = $1 AND d.tenant_id = $2 AND d.restaurant_id = $3
		  AND o.tenant_id = $2 AND o.restaurant_id = $3
		  AND da.status IN ('pending', 'accepted', 'in_progress')
		ORDER BY da.assigned_at DESC
		LIMIT 1
	`

	err := r.db.QueryRowContext(ctx, query, driverID, tenantID, restaurantID).Scan(
		&assignment.ID, &assignment.OrderID, &assignment.DriverID, &assignment.AssignedAt,
		&assignment.AcceptedAt, &assignment.StartedAt, &assignment.CompletedAt,
		&assignment.Status, &assignment.AssignmentNotes, &assignment.Rating,
		&assignment.RatingComment, &assignment.CreatedAt, &assignment.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil // No active assignment is not an error
	}

	return assignment, err
}

// GetAssignmentByOrderID retrieves assignment by order ID
func (r *driverRepository) GetAssignmentByOrderID(ctx context.Context, orderID, tenantID, restaurantID int) (*domain.DriverAssignment, error) {
	assignment := &domain.DriverAssignment{}
	query := `
		SELECT da.id, da.order_id, da.driver_id, da.assigned_at, da.accepted_at, da.started_at, da.completed_at,
		       da.status, da.assignment_notes, da.rating, da.rating_comment, da.created_at, da.updated_at
		FROM driver_assignments da
		JOIN orders o ON o.id = da.order_id
		JOIN drivers d ON d.id = da.driver_id
		WHERE da.order_id = $1 AND o.tenant_id = $2 AND o.restaurant_id = $3
		  AND d.tenant_id = $2 AND d.restaurant_id = $3
		LIMIT 1
	`

	err := r.db.QueryRowContext(ctx, query, orderID, tenantID, restaurantID).Scan(
		&assignment.ID, &assignment.OrderID, &assignment.DriverID, &assignment.AssignedAt,
		&assignment.AcceptedAt, &assignment.StartedAt, &assignment.CompletedAt,
		&assignment.Status, &assignment.AssignmentNotes, &assignment.Rating,
		&assignment.RatingComment, &assignment.CreatedAt, &assignment.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return assignment, err
}

// UpdateAssignmentStatus updates assignment status
func (r *driverRepository) UpdateAssignmentStatus(ctx context.Context, assignmentID, tenantID, restaurantID int, status string) error {
	now := time.Now()
	query := `UPDATE driver_assignments SET status = $1, updated_at = $2`
	args := []interface{}{status, now}

	switch status {
	case "accepted":
		query += fmt.Sprintf(", accepted_at = $%d", len(args)+1)
		args = append(args, now)
	case "in_progress":
		query += fmt.Sprintf(", started_at = $%d", len(args)+1)
		args = append(args, now)
	case "completed":
		query += fmt.Sprintf(", completed_at = $%d", len(args)+1)
		args = append(args, now)
	}

	query += fmt.Sprintf(` WHERE id = $%d
		AND EXISTS (
			SELECT 1 FROM drivers d JOIN orders o ON o.id = driver_assignments.order_id
			WHERE d.id = driver_assignments.driver_id
			  AND d.tenant_id = $%d AND d.restaurant_id = $%d
			  AND o.tenant_id = $%d AND o.restaurant_id = $%d
		)`, len(args)+1, len(args)+2, len(args)+3, len(args)+2, len(args)+3)
	args = append(args, assignmentID)
	args = append(args, tenantID, restaurantID)

	result, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("update assignment status: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("get assignment status affected rows: %w", err)
	}
	if rows == 0 {
		return errors.New("assignment not found")
	}
	return nil
}

// RateDriver rates a driver for delivery
func (r *driverRepository) RateDriver(ctx context.Context, assignmentID, tenantID, restaurantID, rating int, comment string) error {
	query := `
		UPDATE driver_assignments SET rating = $1, rating_comment = $2, updated_at = NOW()
		WHERE id = $3 AND EXISTS (
			SELECT 1 FROM drivers d JOIN orders o ON o.id = driver_assignments.order_id
			WHERE d.id = driver_assignments.driver_id
			  AND d.tenant_id = $4 AND d.restaurant_id = $5
			  AND o.tenant_id = $4 AND o.restaurant_id = $5
		)
	`

	result, err := r.db.ExecContext(ctx, query, rating, comment, assignmentID, tenantID, restaurantID)
	if err != nil {
		return fmt.Errorf("rate driver assignment: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("get rating affected rows: %w", err)
	}
	if rows == 0 {
		return errors.New("assignment not found")
	}
	return nil
}

// GetAssignmentHistory retrieves assignment history for a driver
func (r *driverRepository) GetAssignmentHistory(ctx context.Context, driverID, tenantID, restaurantID, limit int) ([]domain.DriverAssignment, error) {
	query := `
		SELECT da.id, da.order_id, da.driver_id, da.assigned_at, da.accepted_at, da.started_at, da.completed_at,
		       da.status, da.assignment_notes, da.rating, da.rating_comment, da.created_at, da.updated_at
		FROM driver_assignments da
		JOIN drivers d ON d.id = da.driver_id
		JOIN orders o ON o.id = da.order_id
		WHERE da.driver_id = $1 AND d.tenant_id = $2 AND d.restaurant_id = $3
		  AND o.tenant_id = $2 AND o.restaurant_id = $3
		ORDER BY da.assigned_at DESC
		LIMIT $4
	`

	rows, err := r.db.QueryContext(ctx, query, driverID, tenantID, restaurantID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query assignments: %w", err)
	}
	defer rows.Close()

	var assignments []domain.DriverAssignment
	for rows.Next() {
		assignment := domain.DriverAssignment{}
		err := rows.Scan(
			&assignment.ID, &assignment.OrderID, &assignment.DriverID, &assignment.AssignedAt,
			&assignment.AcceptedAt, &assignment.StartedAt, &assignment.CompletedAt,
			&assignment.Status, &assignment.AssignmentNotes, &assignment.Rating,
			&assignment.RatingComment, &assignment.CreatedAt, &assignment.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan assignment: %w", err)
		}
		assignments = append(assignments, assignment)
	}

	return assignments, rows.Err()
}

// CompleteAssignment marks an assignment as completed
func (r *driverRepository) CompleteAssignment(ctx context.Context, assignmentID, tenantID, restaurantID int) error {
	return r.UpdateAssignmentStatus(ctx, assignmentID, tenantID, restaurantID, "completed")
}

// GetDriverStats retrieves driver performance statistics
func (r *driverRepository) GetDriverStats(ctx context.Context, driverID, tenantID, restaurantID int) (*domain.DriverStats, error) {
	stats := &domain.DriverStats{}
	query := `
		SELECT id, first_name, last_name, total_deliveries, completed_deliveries,
		       cancelled_deliveries, average_rating
		FROM drivers
		WHERE id = $1 AND tenant_id = $2 AND restaurant_id = $3
	`

	err := r.db.QueryRowContext(ctx, query, driverID, tenantID, restaurantID).Scan(
		&stats.DriverID, &stats.FirstName, &stats.LastName, &stats.TotalDeliveries,
		&stats.CompletedDeliveries, &stats.CancelledDeliveries, &stats.AverageRating,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("driver not found")
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get driver stats: %w", err)
	}

	// Calculate completion rate
	if stats.TotalDeliveries > 0 {
		stats.CompletionRate = float64(stats.CompletedDeliveries) / float64(stats.TotalDeliveries) * 100
	}

	return stats, nil
}

// UpdateDeliveryMetrics updates driver delivery metrics
func (r *driverRepository) UpdateDeliveryMetrics(ctx context.Context, driverID, tenantID, restaurantID int, completed bool) error {
	var query string
	if completed {
		query = `
			UPDATE drivers SET
				total_deliveries = total_deliveries + 1,
				completed_deliveries = completed_deliveries + 1,
				updated_at = NOW()
			WHERE id = $1 AND tenant_id = $2 AND restaurant_id = $3
		`
	} else {
		query = `
			UPDATE drivers SET
				total_deliveries = total_deliveries + 1,
				cancelled_deliveries = cancelled_deliveries + 1,
				updated_at = NOW()
			WHERE id = $1 AND tenant_id = $2 AND restaurant_id = $3
		`
	}

	_, err := r.db.ExecContext(ctx, query, driverID, tenantID, restaurantID)
	return err
}
