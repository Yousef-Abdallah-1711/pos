package repository

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"pos-saas/internal/domain"
	"strings"
)

type tableScanner interface {
	Scan(dest ...interface{}) error
}

// TableRepository handles table-management persistence.
type TableRepository struct {
	db *sql.DB
}

func NewTableRepository(db *sql.DB) *TableRepository {
	return &TableRepository{db: db}
}

func (r *TableRepository) CreateTable(table *domain.RestaurantTable) (*domain.RestaurantTable, error) {
	if r.db == nil {
		return nil, errors.New("database is not initialized")
	}
	if table == nil {
		return nil, errors.New("table payload is required")
	}

	query := `
		INSERT INTO restaurant_tables (
			tenant_id, restaurant_id, table_number, name, capacity, location, status,
			description, qr_code_config, position_x, position_y, shape, layout_id
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7,
			$8, $9, $10, $11, $12, $13
		)
		RETURNING id, created_at, updated_at
	`

	var name sql.NullString
	if strings.TrimSpace(table.Name) != "" {
		name = sql.NullString{String: strings.TrimSpace(table.Name), Valid: true}
	}

	var description sql.NullString
	if strings.TrimSpace(table.Description) != "" {
		description = sql.NullString{String: strings.TrimSpace(table.Description), Valid: true}
	}

	var layoutID sql.NullInt64
	if table.LayoutID != nil && *table.LayoutID > 0 {
		layoutID = sql.NullInt64{Int64: *table.LayoutID, Valid: true}
	}

	qrConfig := table.QRCodeConfig
	if len(qrConfig) == 0 {
		qrConfig = json.RawMessage("{}")
	}

	err := r.db.QueryRow(
		query,
		table.TenantID,
		table.RestaurantID,
		table.TableNumber,
		name,
		table.Capacity,
		table.Location,
		table.Status,
		description,
		qrConfig,
		table.PositionX,
		table.PositionY,
		table.Shape,
		layoutID,
	).Scan(&table.ID, &table.CreatedAt, &table.UpdatedAt)
	if err != nil {
		if strings.Contains(err.Error(), "uq_restaurant_tables_number") || strings.Contains(err.Error(), "duplicate key") {
			return nil, errors.New("table number already exists for this restaurant")
		}
		return nil, fmt.Errorf("failed to create table: %w", err)
	}

	table.QRCodeConfig = qrConfig
	return table, nil
}

func (r *TableRepository) GetTableByID(tenantID, restaurantID, tableID int64) (*domain.RestaurantTable, error) {
	if r.db == nil {
		return nil, errors.New("database is not initialized")
	}
	query := `
		SELECT
			id, tenant_id, restaurant_id, table_number, name, capacity, location, status,
			description, qr_code_config, position_x, position_y, shape, layout_id, created_at, updated_at
		FROM restaurant_tables
		WHERE id = $1 AND tenant_id = $2 AND restaurant_id = $3
	`
	table, err := scanRestaurantTable(r.db.QueryRow(query, tableID, tenantID, restaurantID))
	if err != nil {
		return nil, err
	}
	return table, nil
}

func (r *TableRepository) GetTableByNumber(tenantID, restaurantID int64, tableNumber int) (*domain.RestaurantTable, error) {
	if r.db == nil {
		return nil, errors.New("database is not initialized")
	}
	query := `
		SELECT
			id, tenant_id, restaurant_id, table_number, name, capacity, location, status,
			description, qr_code_config, position_x, position_y, shape, layout_id, created_at, updated_at
		FROM restaurant_tables
		WHERE tenant_id = $1 AND restaurant_id = $2 AND table_number = $3
	`
	table, err := scanRestaurantTable(r.db.QueryRow(query, tenantID, restaurantID, tableNumber))
	if err != nil {
		return nil, err
	}
	return table, nil
}

func (r *TableRepository) ListTables(tenantID, restaurantID int64, status string) ([]domain.RestaurantTable, error) {
	if r.db == nil {
		return nil, errors.New("database is not initialized")
	}

	query := `
		SELECT
			id, tenant_id, restaurant_id, table_number, name, capacity, location, status,
			description, qr_code_config, position_x, position_y, shape, layout_id, created_at, updated_at
		FROM restaurant_tables
		WHERE tenant_id = $1 AND restaurant_id = $2
	`
	args := []interface{}{tenantID, restaurantID}
	if strings.TrimSpace(status) != "" {
		query += " AND status = $3"
		args = append(args, strings.TrimSpace(status))
	}
	query += " ORDER BY table_number ASC, created_at ASC"

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list tables: %w", err)
	}
	defer rows.Close()

	tables := make([]domain.RestaurantTable, 0)
	for rows.Next() {
		table, err := scanRestaurantTable(rows)
		if err != nil {
			return nil, err
		}
		tables = append(tables, *table)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate tables: %w", err)
	}
	return tables, nil
}

func (r *TableRepository) UpdateTable(tenantID, restaurantID, tableID int64, req *domain.UpdateTableRequest) error {
	if r.db == nil {
		return errors.New("database is not initialized")
	}
	if req == nil {
		return errors.New("update payload is required")
	}

	current, err := r.GetTableByID(tenantID, restaurantID, tableID)
	if err != nil {
		return err
	}

	if req.TableNumber != nil {
		current.TableNumber = *req.TableNumber
	}
	if req.Name != nil {
		current.Name = strings.TrimSpace(*req.Name)
	}
	if req.Capacity != nil {
		current.Capacity = *req.Capacity
	}
	if req.Location != nil {
		current.Location = domain.TableLocation(strings.TrimSpace(*req.Location))
	}
	if req.Status != nil {
		current.Status = domain.TableStatus(strings.TrimSpace(*req.Status))
	}
	if req.Description != nil {
		current.Description = strings.TrimSpace(*req.Description)
	}
	if req.PositionX != nil {
		current.PositionX = *req.PositionX
	}
	if req.PositionY != nil {
		current.PositionY = *req.PositionY
	}
	if req.Shape != nil {
		current.Shape = domain.TableShape(strings.TrimSpace(*req.Shape))
	}
	if req.LayoutID != nil {
		current.LayoutID = req.LayoutID
	}

	var name sql.NullString
	if current.Name != "" {
		name = sql.NullString{String: current.Name, Valid: true}
	}
	var description sql.NullString
	if current.Description != "" {
		description = sql.NullString{String: current.Description, Valid: true}
	}
	var layoutID sql.NullInt64
	if current.LayoutID != nil && *current.LayoutID > 0 {
		layoutID = sql.NullInt64{Int64: *current.LayoutID, Valid: true}
	}

	query := `
		UPDATE restaurant_tables
		SET
			table_number = $1,
			name = $2,
			capacity = $3,
			location = $4,
			status = $5,
			description = $6,
			position_x = $7,
			position_y = $8,
			shape = $9,
			layout_id = $10,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $11 AND tenant_id = $12 AND restaurant_id = $13
	`

	result, err := r.db.Exec(
		query,
		current.TableNumber,
		name,
		current.Capacity,
		current.Location,
		current.Status,
		description,
		current.PositionX,
		current.PositionY,
		current.Shape,
		layoutID,
		tableID,
		tenantID,
		restaurantID,
	)
	if err != nil {
		if strings.Contains(err.Error(), "uq_restaurant_tables_number") || strings.Contains(err.Error(), "duplicate key") {
			return errors.New("table number already exists for this restaurant")
		}
		return fmt.Errorf("failed to update table: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to inspect updated rows: %w", err)
	}
	if rowsAffected == 0 {
		return errors.New("table not found")
	}
	return nil
}

func (r *TableRepository) DeleteTable(tenantID, restaurantID, tableID int64) error {
	if r.db == nil {
		return errors.New("database is not initialized")
	}
	result, err := r.db.Exec(
		`DELETE FROM restaurant_tables WHERE id = $1 AND tenant_id = $2 AND restaurant_id = $3`,
		tableID, tenantID, restaurantID,
	)
	if err != nil {
		return fmt.Errorf("failed to delete table: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to inspect deleted rows: %w", err)
	}
	if rowsAffected == 0 {
		return errors.New("table not found")
	}
	return nil
}

func (r *TableRepository) UpdateTableStatus(tenantID, restaurantID, tableID int64, status string) error {
	if r.db == nil {
		return errors.New("database is not initialized")
	}
	result, err := r.db.Exec(
		`UPDATE restaurant_tables
		 SET status = $1, updated_at = CURRENT_TIMESTAMP
		 WHERE id = $2 AND tenant_id = $3 AND restaurant_id = $4`,
		status, tableID, tenantID, restaurantID,
	)
	if err != nil {
		return fmt.Errorf("failed to update table status: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to inspect updated rows: %w", err)
	}
	if rowsAffected == 0 {
		return errors.New("table not found")
	}
	return nil
}

func (r *TableRepository) GetTableStats(tenantID, restaurantID int64) (*domain.TableStats, error) {
	if r.db == nil {
		return nil, errors.New("database is not initialized")
	}
	query := `
		SELECT
			COUNT(*) AS total_tables,
			COUNT(*) FILTER (WHERE status = 'available') AS available,
			COUNT(*) FILTER (WHERE status = 'occupied') AS occupied,
			COUNT(*) FILTER (WHERE status = 'reserved') AS reserved,
			COUNT(*) FILTER (WHERE status = 'maintenance') AS maintenance,
			COALESCE(SUM(capacity), 0) AS total_capacity
		FROM restaurant_tables
		WHERE tenant_id = $1 AND restaurant_id = $2
	`
	stats := &domain.TableStats{}
	if err := r.db.QueryRow(query, tenantID, restaurantID).Scan(
		&stats.TotalTables,
		&stats.Available,
		&stats.Occupied,
		&stats.Reserved,
		&stats.Maintenance,
		&stats.TotalCapacity,
	); err != nil {
		return nil, fmt.Errorf("failed to fetch table stats: %w", err)
	}
	return stats, nil
}

func (r *TableRepository) UpdateTablePositions(tenantID, restaurantID int64, positions []domain.TablePosition) error {
	if r.db == nil {
		return errors.New("database is not initialized")
	}
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to start positions transaction: %w", err)
	}
	defer tx.Rollback()

	for _, pos := range positions {
		result, err := tx.Exec(
			`UPDATE restaurant_tables
			 SET position_x = $1, position_y = $2, updated_at = CURRENT_TIMESTAMP
			 WHERE id = $3 AND tenant_id = $4 AND restaurant_id = $5`,
			pos.PositionX, pos.PositionY, pos.TableID, tenantID, restaurantID,
		)
		if err != nil {
			return fmt.Errorf("failed to update table %d position: %w", pos.TableID, err)
		}
		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return fmt.Errorf("failed to inspect updated rows for table %d: %w", pos.TableID, err)
		}
		if rowsAffected == 0 {
			return fmt.Errorf("table %d not found", pos.TableID)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit position updates: %w", err)
	}
	return nil
}

func (r *TableRepository) CreateLayout(layout *domain.TableLayout) (*domain.TableLayout, error) {
	if r.db == nil {
		return nil, errors.New("database is not initialized")
	}
	if layout == nil {
		return nil, errors.New("layout payload is required")
	}

	var backgroundImage sql.NullString
	if strings.TrimSpace(layout.BackgroundImage) != "" {
		backgroundImage = sql.NullString{String: strings.TrimSpace(layout.BackgroundImage), Valid: true}
	}

	query := `
		INSERT INTO table_layouts (
			tenant_id, restaurant_id, name, floor_number, canvas_width, canvas_height, background_image, is_default
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at, updated_at
	`
	if err := r.db.QueryRow(
		query,
		layout.TenantID,
		layout.RestaurantID,
		layout.Name,
		layout.FloorNumber,
		layout.CanvasWidth,
		layout.CanvasHeight,
		backgroundImage,
		layout.IsDefault,
	).Scan(&layout.ID, &layout.CreatedAt, &layout.UpdatedAt); err != nil {
		return nil, fmt.Errorf("failed to create layout: %w", err)
	}
	return layout, nil
}

func (r *TableRepository) GetLayouts(tenantID, restaurantID int64) ([]domain.TableLayout, error) {
	if r.db == nil {
		return nil, errors.New("database is not initialized")
	}
	query := `
		SELECT
			id, tenant_id, restaurant_id, name, floor_number, canvas_width, canvas_height,
			background_image, is_default, created_at, updated_at
		FROM table_layouts
		WHERE tenant_id = $1 AND restaurant_id = $2
		ORDER BY floor_number ASC, name ASC
	`
	rows, err := r.db.Query(query, tenantID, restaurantID)
	if err != nil {
		return nil, fmt.Errorf("failed to list layouts: %w", err)
	}
	defer rows.Close()

	layouts := make([]domain.TableLayout, 0)
	for rows.Next() {
		layout := domain.TableLayout{}
		var backgroundImage sql.NullString
		if err := rows.Scan(
			&layout.ID,
			&layout.TenantID,
			&layout.RestaurantID,
			&layout.Name,
			&layout.FloorNumber,
			&layout.CanvasWidth,
			&layout.CanvasHeight,
			&backgroundImage,
			&layout.IsDefault,
			&layout.CreatedAt,
			&layout.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan layout: %w", err)
		}
		if backgroundImage.Valid {
			layout.BackgroundImage = backgroundImage.String
		}
		layouts = append(layouts, layout)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate layouts: %w", err)
	}
	return layouts, nil
}

func (r *TableRepository) GetTablePublicInfo(tableID int64) (*domain.TablePublicInfo, error) {
	if r.db == nil {
		return nil, errors.New("database is not initialized")
	}
	query := `
		SELECT
			t.id,
			t.table_number,
			COALESCE(t.name, ''),
			t.capacity,
			t.location,
			t.status,
			t.tenant_id,
			t.restaurant_id,
			COALESCE(r.slug, '')
		FROM restaurant_tables t
		JOIN restaurants r ON r.id = t.restaurant_id
		WHERE t.id = $1
	`

	info := &domain.TablePublicInfo{}
	if err := r.db.QueryRow(query, tableID).Scan(
		&info.ID,
		&info.TableNumber,
		&info.Name,
		&info.Capacity,
		&info.Location,
		&info.Status,
		&info.TenantID,
		&info.RestaurantID,
		&info.RestaurantSlug,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("table not found")
		}
		return nil, fmt.Errorf("failed to fetch public table info: %w", err)
	}

	return info, nil
}

func (r *TableRepository) HasActiveOrders(tenantID, restaurantID, tableID int64) (bool, error) {
	if r.db == nil {
		return false, errors.New("database is not initialized")
	}
	query := `
		SELECT COUNT(1)
		FROM orders
		WHERE tenant_id = $1
		  AND restaurant_id = $2
		  AND table_id = $3
		  AND status NOT IN ('delivered', 'cancelled')
	`
	var count int
	if err := r.db.QueryRow(query, tenantID, restaurantID, tableID).Scan(&count); err != nil {
		return false, fmt.Errorf("failed to check active table orders: %w", err)
	}
	return count > 0, nil
}

func scanRestaurantTable(scanner tableScanner) (*domain.RestaurantTable, error) {
	table := &domain.RestaurantTable{}
	var name sql.NullString
	var description sql.NullString
	var qrConfig []byte
	var layoutID sql.NullInt64
	var location string
	var status string
	var shape string

	err := scanner.Scan(
		&table.ID,
		&table.TenantID,
		&table.RestaurantID,
		&table.TableNumber,
		&name,
		&table.Capacity,
		&location,
		&status,
		&description,
		&qrConfig,
		&table.PositionX,
		&table.PositionY,
		&shape,
		&layoutID,
		&table.CreatedAt,
		&table.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("table not found")
		}
		return nil, fmt.Errorf("failed to scan table: %w", err)
	}

	if name.Valid {
		table.Name = name.String
	}
	if description.Valid {
		table.Description = description.String
	}
	if len(qrConfig) == 0 {
		qrConfig = []byte("{}")
	}
	table.QRCodeConfig = json.RawMessage(qrConfig)
	table.Location = domain.TableLocation(location)
	table.Status = domain.TableStatus(status)
	table.Shape = domain.TableShape(shape)
	if layoutID.Valid {
		table.LayoutID = &layoutID.Int64
	}

	return table, nil
}
