package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"pos-saas/internal/domain"
)

// DriverRepository defines operations for driver data access
type DriverRepository interface {
	// Driver Management
	Create(ctx context.Context, driver *domain.Driver) error
	GetByID(ctx context.Context, id int) (*domain.Driver, error)
	GetByEmail(ctx context.Context, email string) (*domain.Driver, error)
	GetByTenantAndRestaurant(ctx context.Context, tenantID, restaurantID int) ([]domain.Driver, error)
	GetByTenantAndRestaurantPaged(ctx context.Context, tenantID, restaurantID, page, pageSize int) ([]domain.Driver, int, error)
	Update(ctx context.Context, driver *domain.Driver) error
	Delete(ctx context.Context, id int) error

	// Status & Availability
	UpdateStatus(ctx context.Context, id int, status string) error
	UpdateAvailability(ctx context.Context, id int, status string) error
	GetAvailableDrivers(ctx context.Context, tenantID, restaurantID int) ([]domain.Driver, error)
	UpdateActiveOrders(ctx context.Context, id int, count int) error

	// Location
	UpdateLocation(ctx context.Context, driverID int, lat, lon float64) error
	GetLastLocation(ctx context.Context, driverID int) (*domain.DriverLocation, error)
	RecordLocation(ctx context.Context, location *domain.DriverLocation) error
	GetLocationHistory(ctx context.Context, driverID int, limit int) ([]domain.DriverLocation, error)

	// Assignments
	AssignOrder(ctx context.Context, assignment *domain.DriverAssignment) error
	GetAssignmentByID(ctx context.Context, id int) (*domain.DriverAssignment, error)
	GetActiveAssignment(ctx context.Context, driverID int) (*domain.DriverAssignment, error)
	GetAssignmentByOrderID(ctx context.Context, orderID int) (*domain.DriverAssignment, error)
	UpdateAssignmentStatus(ctx context.Context, assignmentID int, status string) error
	RateDriver(ctx context.Context, assignmentID int, rating int, comment string) error
	GetAssignmentHistory(ctx context.Context, driverID int, limit int) ([]domain.DriverAssignment, error)
	CompleteAssignment(ctx context.Context, assignmentID int) error

	// Statistics
	GetDriverStats(ctx context.Context, driverID int) (*domain.DriverStats, error)
	UpdateDeliveryMetrics(ctx context.Context, driverID int, completed bool) error
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
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
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
func (r *driverRepository) GetByID(ctx context.Context, id int) (*domain.Driver, error) {
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
		FROM drivers WHERE id = $1
	`

	err := r.db.QueryRowContext(ctx, query, id).Scan(
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
func (r *driverRepository) GetByEmail(ctx context.Context, email string) (*domain.Driver, error) {
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
		FROM drivers WHERE email = $1 LIMIT 1
	`

	err := r.db.QueryRowContext(ctx, query, email).Scan(
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
	drivers, _, err := r.GetByTenantAndRestaurantPaged(ctx, tenantID, restaurantID, 1, 1000000)
	return drivers, err
}

func (r *driverRepository) GetByTenantAndRestaurantPaged(ctx context.Context, tenantID, restaurantID, page, pageSize int) ([]domain.Driver, int, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 200 {
		pageSize = 200
	}
	offset := (page - 1) * pageSize

	const countQuery = `
		SELECT COUNT(*)
		FROM drivers
		WHERE tenant_id = $1 AND restaurant_id = $2
	`
	var total int
	if err := r.db.QueryRowContext(ctx, countQuery, tenantID, restaurantID).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("failed to count drivers: %w", err)
	}

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
		LIMIT $3 OFFSET $4
	`

	rows, err := r.db.QueryContext(ctx, query, tenantID, restaurantID, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query drivers: %w", err)
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
			return nil, 0, fmt.Errorf("failed to scan driver: %w", err)
		}
		drivers = append(drivers, driver)
	}

	return drivers, total, rows.Err()
}

// Update updates an existing driver
func (r *driverRepository) Update(ctx context.Context, driver *domain.Driver) error {
	query := `
		UPDATE drivers SET
			first_name = $1, last_name = $2, email = $3, phone_number = $4,
			license_number = $5, vehicle_type = $6, vehicle_number = $7,
			status = $8, availability_status = $9, date_of_birth = $10,
			address = $11, city = $12, state = $13, zip_code = $14,
			notes = $15, updated_at = NOW()
		WHERE id = $16
	`

	result, err := r.db.ExecContext(
		ctx, query,
		driver.FirstName, driver.LastName, driver.Email, driver.PhoneNumber,
		driver.LicenseNumber, driver.VehicleType, driver.VehicleNumber,
		driver.Status, driver.AvailabilityStatus, driver.DateOfBirth,
		driver.Address, driver.City, driver.State, driver.ZipCode,
		driver.Notes, driver.ID,
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
func (r *driverRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM drivers WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
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
func (r *driverRepository) UpdateStatus(ctx context.Context, id int, status string) error {
	query := `UPDATE drivers SET status = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, status, id)
	return err
}

// UpdateAvailability updates driver availability status
func (r *driverRepository) UpdateAvailability(ctx context.Context, id int, status string) error {
	query := `UPDATE drivers SET availability_status = $1, last_active = NOW(), updated_at = NOW() WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, status, id)
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
func (r *driverRepository) UpdateActiveOrders(ctx context.Context, id int, count int) error {
	query := `UPDATE drivers SET active_orders_count = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, count, id)
	return err
}

// UpdateLocation updates driver current location
func (r *driverRepository) UpdateLocation(ctx context.Context, driverID int, lat, lon float64) error {
	query := `UPDATE drivers SET current_latitude = $1, current_longitude = $2, last_active = NOW(), updated_at = NOW() WHERE id = $3`
	_, err := r.db.ExecContext(ctx, query, lat, lon, driverID)
	return err
}

// GetLastLocation retrieves the last known location of a driver
func (r *driverRepository) GetLastLocation(ctx context.Context, driverID int) (*domain.DriverLocation, error) {
	location := &domain.DriverLocation{}
	query := `
		SELECT id, driver_id, order_id, latitude, longitude, accuracy, recorded_at
		FROM driver_location_history
		WHERE driver_id = $1
		ORDER BY recorded_at DESC
		LIMIT 1
	`

	err := r.db.QueryRowContext(ctx, query, driverID).Scan(
		&location.ID, &location.DriverID, &location.OrderID, &location.Latitude, &location.Longitude, &location.Accuracy, &location.RecordedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("no location found")
	}

	return location, err
}

// RecordLocation records a new location entry
func (r *driverRepository) RecordLocation(ctx context.Context, location *domain.DriverLocation) error {
	query := `
		INSERT INTO driver_location_history (driver_id, order_id, latitude, longitude, accuracy)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, recorded_at
	`

	return r.db.QueryRowContext(ctx, query, location.DriverID, location.OrderID, location.Latitude, location.Longitude, location.Accuracy).Scan(
		&location.ID, &location.RecordedAt,
	)
}

// GetLocationHistory retrieves location history for a driver
func (r *driverRepository) GetLocationHistory(ctx context.Context, driverID int, limit int) ([]domain.DriverLocation, error) {
	query := `
		SELECT id, driver_id, order_id, latitude, longitude, accuracy, recorded_at
		FROM driver_location_history
		WHERE driver_id = $1
		ORDER BY recorded_at DESC
		LIMIT $2
	`

	rows, err := r.db.QueryContext(ctx, query, driverID, limit)
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
func (r *driverRepository) AssignOrder(ctx context.Context, assignment *domain.DriverAssignment) error {
	query := `
		INSERT INTO driver_assignments (order_id, driver_id, assigned_at, status, assignment_notes)
		VALUES ($1, $2, NOW(), $3, $4)
		RETURNING id, created_at, updated_at
	`

	return r.db.QueryRowContext(ctx, query, assignment.OrderID, assignment.DriverID, assignment.Status, assignment.AssignmentNotes).Scan(
		&assignment.ID, &assignment.CreatedAt, &assignment.UpdatedAt,
	)
}

// GetAssignmentByID retrieves assignment by ID
func (r *driverRepository) GetAssignmentByID(ctx context.Context, id int) (*domain.DriverAssignment, error) {
	assignment := &domain.DriverAssignment{}
	query := `
		SELECT id, order_id, driver_id, assigned_at, accepted_at, started_at, completed_at,
		       status, assignment_notes, rating, rating_comment, created_at, updated_at
		FROM driver_assignments
		WHERE id = $1
	`

	err := r.db.QueryRowContext(ctx, query, id).Scan(
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
func (r *driverRepository) GetActiveAssignment(ctx context.Context, driverID int) (*domain.DriverAssignment, error) {
	assignment := &domain.DriverAssignment{}
	query := `
		SELECT id, order_id, driver_id, assigned_at, accepted_at, started_at, completed_at,
		       status, assignment_notes, rating, rating_comment, created_at, updated_at
		FROM driver_assignments
		WHERE driver_id = $1 AND status IN ('pending', 'accepted', 'in_progress')
		ORDER BY assigned_at DESC
		LIMIT 1
	`

	err := r.db.QueryRowContext(ctx, query, driverID).Scan(
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
func (r *driverRepository) GetAssignmentByOrderID(ctx context.Context, orderID int) (*domain.DriverAssignment, error) {
	assignment := &domain.DriverAssignment{}
	query := `
		SELECT id, order_id, driver_id, assigned_at, accepted_at, started_at, completed_at,
		       status, assignment_notes, rating, rating_comment, created_at, updated_at
		FROM driver_assignments
		WHERE order_id = $1
		LIMIT 1
	`

	err := r.db.QueryRowContext(ctx, query, orderID).Scan(
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
func (r *driverRepository) UpdateAssignmentStatus(ctx context.Context, assignmentID int, status string) error {
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

	query += fmt.Sprintf(" WHERE id = $%d", len(args)+1)
	args = append(args, assignmentID)

	_, err := r.db.ExecContext(ctx, query, args...)
	return err
}

// RateDriver rates a driver for delivery
func (r *driverRepository) RateDriver(ctx context.Context, assignmentID int, rating int, comment string) error {
	query := `
		UPDATE driver_assignments SET rating = $1, rating_comment = $2, updated_at = NOW()
		WHERE id = $3
	`

	_, err := r.db.ExecContext(ctx, query, rating, comment, assignmentID)
	return err
}

// GetAssignmentHistory retrieves assignment history for a driver
func (r *driverRepository) GetAssignmentHistory(ctx context.Context, driverID int, limit int) ([]domain.DriverAssignment, error) {
	query := `
		SELECT id, order_id, driver_id, assigned_at, accepted_at, started_at, completed_at,
		       status, assignment_notes, rating, rating_comment, created_at, updated_at
		FROM driver_assignments
		WHERE driver_id = $1
		ORDER BY assigned_at DESC
		LIMIT $2
	`

	rows, err := r.db.QueryContext(ctx, query, driverID, limit)
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
func (r *driverRepository) CompleteAssignment(ctx context.Context, assignmentID int) error {
	return r.UpdateAssignmentStatus(ctx, assignmentID, "completed")
}

// GetDriverStats retrieves driver performance statistics
func (r *driverRepository) GetDriverStats(ctx context.Context, driverID int) (*domain.DriverStats, error) {
	stats := &domain.DriverStats{}
	query := `
		SELECT id, first_name, last_name, total_deliveries, completed_deliveries,
		       cancelled_deliveries, average_rating
		FROM drivers
		WHERE id = $1
	`

	err := r.db.QueryRowContext(ctx, query, driverID).Scan(
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
func (r *driverRepository) UpdateDeliveryMetrics(ctx context.Context, driverID int, completed bool) error {
	var query string
	if completed {
		query = `
			UPDATE drivers SET
				total_deliveries = total_deliveries + 1,
				completed_deliveries = completed_deliveries + 1,
				updated_at = NOW()
			WHERE id = $1
		`
	} else {
		query = `
			UPDATE drivers SET
				total_deliveries = total_deliveries + 1,
				cancelled_deliveries = cancelled_deliveries + 1,
				updated_at = NOW()
			WHERE id = $1
		`
	}

	_, err := r.db.ExecContext(ctx, query, driverID)
	return err
}
