package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/lib/pq"

	"pos-saas/internal/domain"
)

// OrderRepository handles order data operations
type OrderRepository struct {
	db *sql.DB
}

// NewOrderRepository creates new order repository
func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

// CreateOrder inserts a new order into database
func (r *OrderRepository) CreateOrder(order *domain.Order) (*domain.Order, error) {
	return r.insertOrder(r.db, order)
}

func (r *OrderRepository) insertOrder(writer orderRowWriter, order *domain.Order) (*domain.Order, error) {
	query := `
		INSERT INTO orders (
			tenant_id, restaurant_id, order_number, customer_name, customer_email,
			customer_phone, delivery_address, delivery_city, delivery_area,
			delivery_zip_code, delivery_latitude, delivery_longitude,
			delivery_instructions, subtotal, tax_amount, discount_amount,
			delivery_fee, total_amount, payment_method, payment_status,
			status, estimated_delivery_time, notes, order_source, table_id
		) SELECT
			$1::bigint, $2::bigint, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15,
			$16, $17, $18, $19, $20, $21, $22, $23, $24, $25::bigint
		WHERE EXISTS (
			SELECT 1 FROM restaurants WHERE id = $2::bigint AND tenant_id = $1::bigint
		)
		AND ($25::bigint IS NULL OR EXISTS (
			SELECT 1 FROM restaurant_tables
			WHERE id = $25::bigint AND tenant_id = $1::bigint AND restaurant_id = $2::bigint
		))
		RETURNING id, created_at, updated_at
	`

	// Handle nullable fields
	var customerEmail, deliveryAddress, deliveryCity, deliveryArea, deliveryZipCode, deliveryInstructions sql.NullString
	var deliveryLatitude, deliveryLongitude sql.NullFloat64
	var paymentMethod sql.NullString
	var estimatedDeliveryTime sql.NullTime
	var tableID sql.NullInt64

	if order.CustomerEmail != "" {
		customerEmail = sql.NullString{String: order.CustomerEmail, Valid: true}
	}
	if order.DeliveryAddress != "" {
		deliveryAddress = sql.NullString{String: order.DeliveryAddress, Valid: true}
	}
	if order.DeliveryCity != "" {
		deliveryCity = sql.NullString{String: order.DeliveryCity, Valid: true}
	}
	if order.DeliveryArea != "" {
		deliveryArea = sql.NullString{String: order.DeliveryArea, Valid: true}
	}
	if order.DeliveryZipCode != "" {
		deliveryZipCode = sql.NullString{String: order.DeliveryZipCode, Valid: true}
	}
	if order.DeliveryInstructions != "" {
		deliveryInstructions = sql.NullString{String: order.DeliveryInstructions, Valid: true}
	}
	if order.DeliveryLatitude != 0 {
		deliveryLatitude = sql.NullFloat64{Float64: order.DeliveryLatitude, Valid: true}
	}
	if order.DeliveryLongitude != 0 {
		deliveryLongitude = sql.NullFloat64{Float64: order.DeliveryLongitude, Valid: true}
	}
	if order.PaymentMethod != "" {
		paymentMethod = sql.NullString{String: order.PaymentMethod, Valid: true}
	}
	if order.EstimatedDeliveryTime != nil {
		estimatedDeliveryTime = sql.NullTime{Time: *order.EstimatedDeliveryTime, Valid: true}
	}
	if order.TableID != nil && *order.TableID > 0 {
		tableID = sql.NullInt64{Int64: *order.TableID, Valid: true}
	}

	err := writer.QueryRow(query,
		order.TenantID,
		order.RestaurantID,
		order.OrderNumber,
		order.CustomerName,
		customerEmail,
		order.CustomerPhone,
		deliveryAddress,
		deliveryCity,
		deliveryArea,
		deliveryZipCode,
		deliveryLatitude,
		deliveryLongitude,
		deliveryInstructions,
		order.Subtotal,
		order.TaxAmount,
		order.DiscountAmount,
		order.DeliveryFee,
		order.TotalAmount,
		paymentMethod,
		order.PaymentStatus,
		order.Status,
		estimatedDeliveryTime,
		order.Notes,
		order.OrderSource,
		tableID,
	).Scan(&order.ID, &order.CreatedAt, &order.UpdatedAt)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			return nil, errors.New("order with this order number already exists")
		}
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	return order, nil
}

// GetOrderByID retrieves an order by ID with all items
func (r *OrderRepository) GetOrderByID(tenantID, restaurantID, orderID int64) (*domain.Order, error) {
	query := `
		SELECT
			o.id, o.tenant_id, o.restaurant_id, o.order_number, o.customer_name, o.customer_email,
			o.customer_phone, o.delivery_address, o.delivery_city, o.delivery_area,
			o.delivery_zip_code, o.delivery_latitude, o.delivery_longitude,
			o.delivery_instructions, o.subtotal, o.tax_amount, o.discount_amount,
			o.delivery_fee, o.total_amount, o.payment_method, o.payment_status,
			o.status, o.estimated_delivery_time, o.actual_delivery_time, o.notes, o.order_source,
			o.table_id, t.table_number, COALESCE(t.name, ''), COALESCE(t.location, ''), t.capacity,
			o.created_at, o.updated_at
		FROM orders o
		LEFT JOIN restaurant_tables t
			ON t.id = o.table_id
			AND t.tenant_id = o.tenant_id
			AND t.restaurant_id = o.restaurant_id
		WHERE o.id = $1 AND o.tenant_id = $2 AND o.restaurant_id = $3
	`

	order := &domain.Order{}
	var customerEmail, deliveryAddress, deliveryCity, deliveryArea, deliveryZipCode, deliveryInstructions sql.NullString
	var deliveryLatitude, deliveryLongitude sql.NullFloat64
	var paymentMethod, notes sql.NullString
	var estimatedDeliveryTime, actualDeliveryTime sql.NullTime
	var tableID sql.NullInt64
	var tableNumber sql.NullInt32
	var tableName sql.NullString
	var tableLocation sql.NullString
	var tableCapacity sql.NullInt32

	err := r.db.QueryRow(query, orderID, tenantID, restaurantID).Scan(
		&order.ID, &order.TenantID, &order.RestaurantID, &order.OrderNumber,
		&order.CustomerName, &customerEmail,
		&order.CustomerPhone, &deliveryAddress, &deliveryCity, &deliveryArea,
		&deliveryZipCode, &deliveryLatitude, &deliveryLongitude,
		&deliveryInstructions, &order.Subtotal, &order.TaxAmount, &order.DiscountAmount,
		&order.DeliveryFee, &order.TotalAmount, &paymentMethod, &order.PaymentStatus,
		&order.Status, &estimatedDeliveryTime, &actualDeliveryTime, &notes, &order.OrderSource,
		&tableID, &tableNumber, &tableName, &tableLocation, &tableCapacity,
		&order.CreatedAt, &order.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("order not found")
		}
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	// Handle nullable fields
	if customerEmail.Valid {
		order.CustomerEmail = customerEmail.String
	}
	if deliveryAddress.Valid {
		order.DeliveryAddress = deliveryAddress.String
	}
	if deliveryCity.Valid {
		order.DeliveryCity = deliveryCity.String
	}
	if deliveryArea.Valid {
		order.DeliveryArea = deliveryArea.String
	}
	if deliveryZipCode.Valid {
		order.DeliveryZipCode = deliveryZipCode.String
	}
	if deliveryInstructions.Valid {
		order.DeliveryInstructions = deliveryInstructions.String
	}
	if deliveryLatitude.Valid {
		order.DeliveryLatitude = deliveryLatitude.Float64
	}
	if deliveryLongitude.Valid {
		order.DeliveryLongitude = deliveryLongitude.Float64
	}
	if paymentMethod.Valid {
		order.PaymentMethod = paymentMethod.String
	}
	if estimatedDeliveryTime.Valid {
		order.EstimatedDeliveryTime = &estimatedDeliveryTime.Time
	}
	if actualDeliveryTime.Valid {
		order.ActualDeliveryTime = &actualDeliveryTime.Time
	}
	if notes.Valid {
		order.Notes = notes.String
	}
	if tableID.Valid {
		order.TableID = &tableID.Int64
	}
	if tableNumber.Valid {
		tableNumberValue := int(tableNumber.Int32)
		order.TableNumber = &tableNumberValue
	}
	if tableName.Valid {
		order.TableName = tableName.String
	}
	if tableLocation.Valid {
		order.TableLocation = tableLocation.String
	}
	if tableCapacity.Valid {
		tableCapacityValue := int(tableCapacity.Int32)
		order.TableCapacity = &tableCapacityValue
	}

	// Get order items
	items, err := r.GetOrderItems(tenantID, restaurantID, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get order items: %w", err)
	}
	order.Items = items

	return order, nil
}

// GetOrderByNumber retrieves an order by order number
func (r *OrderRepository) GetOrderByNumber(tenantID, restaurantID int64, orderNumber string) (*domain.Order, error) {
	query := `
		SELECT id FROM orders
		WHERE order_number = $1 AND tenant_id = $2 AND restaurant_id = $3
	`

	var orderID int64
	err := r.db.QueryRow(query, orderNumber, tenantID, restaurantID).Scan(&orderID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("order not found")
		}
		return nil, fmt.Errorf("failed to get order by number: %w", err)
	}

	return r.GetOrderByID(tenantID, restaurantID, orderID)
}

// ListOrders retrieves paginated orders for a restaurant
func (r *OrderRepository) ListOrders(tenantID, restaurantID int64, filters *domain.OrderListFilters) (*domain.OrderListResponse, error) {
	baseQuery := `
		FROM orders o
		LEFT JOIN restaurant_tables t
		  ON t.id = o.table_id
		 AND t.tenant_id = o.tenant_id
		 AND t.restaurant_id = o.restaurant_id
	`
	if filters != nil && filters.AssignedDriverID != nil {
		baseQuery += `
		JOIN driver_assignments da ON da.order_id = o.id
		`
	}
	baseQuery += `
		WHERE o.tenant_id = $1 AND o.restaurant_id = $2
	`
	query := `
		SELECT DISTINCT
			o.id, o.tenant_id, o.restaurant_id, o.order_number, o.customer_name, o.customer_email,
			o.customer_phone, o.delivery_address, o.subtotal, o.tax_amount, o.discount_amount, o.delivery_fee,
			o.total_amount, o.payment_method, o.payment_status, o.status, o.order_source, o.table_id,
			t.table_number, COALESCE(t.name, ''), COALESCE(t.location, ''), t.capacity,
			o.created_at, o.updated_at
	` + baseQuery
	countQuery := "SELECT COUNT(DISTINCT o.id) " + baseQuery

	args := []interface{}{tenantID, restaurantID}
	argCount := 2

	// Apply filters
	if filters != nil {
		if filters.Status != "" && domain.ValidStatus(filters.Status) {
			argCount++
			clause := fmt.Sprintf(" AND o.status = $%d", argCount)
			query += clause
			countQuery += clause
			args = append(args, filters.Status)
		}

		if filters.PaymentStatus != "" && domain.ValidPaymentStatus(filters.PaymentStatus) {
			argCount++
			clause := fmt.Sprintf(" AND o.payment_status = $%d", argCount)
			query += clause
			countQuery += clause
			args = append(args, filters.PaymentStatus)
		}

		if filters.CustomerName != "" {
			argCount++
			clause := fmt.Sprintf(" AND o.customer_name ILIKE $%d", argCount)
			query += clause
			countQuery += clause
			args = append(args, "%"+filters.CustomerName+"%")
		}

		if filters.CustomerEmail != "" {
			argCount++
			clause := fmt.Sprintf(" AND o.customer_email ILIKE $%d", argCount)
			query += clause
			countQuery += clause
			args = append(args, "%"+filters.CustomerEmail+"%")
		}

		if filters.OrderSource != "" {
			argCount++
			clause := fmt.Sprintf(" AND o.order_source = $%d", argCount)
			query += clause
			countQuery += clause
			args = append(args, filters.OrderSource)
		}

		if filters.TableID != nil && *filters.TableID > 0 {
			argCount++
			clause := fmt.Sprintf(" AND o.table_id = $%d", argCount)
			query += clause
			countQuery += clause
			args = append(args, *filters.TableID)
		}

		if !filters.StartDate.IsZero() {
			argCount++
			clause := fmt.Sprintf(" AND o.created_at >= $%d", argCount)
			query += clause
			countQuery += clause
			args = append(args, filters.StartDate)
		}

		if !filters.EndDate.IsZero() {
			argCount++
			clause := fmt.Sprintf(" AND o.created_at <= $%d", argCount)
			query += clause
			countQuery += clause
			args = append(args, filters.EndDate)
		}

		if filters.MinAmount > 0 {
			argCount++
			clause := fmt.Sprintf(" AND o.total_amount >= $%d", argCount)
			query += clause
			countQuery += clause
			args = append(args, filters.MinAmount)
		}

		if filters.MaxAmount > 0 {
			argCount++
			clause := fmt.Sprintf(" AND o.total_amount <= $%d", argCount)
			query += clause
			countQuery += clause
			args = append(args, filters.MaxAmount)
		}

		if filters.AssignedDriverID != nil {
			argCount++
			clause := fmt.Sprintf(" AND da.driver_id = $%d", argCount)
			query += clause
			countQuery += clause
			args = append(args, *filters.AssignedDriverID)
		}
	}

	// Count total matching orders
	var total int64
	err := r.db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("failed to count orders: %w", err)
	}

	// Add ordering and pagination
	query += " ORDER BY o.created_at DESC"

	if filters != nil && filters.Limit > 0 {
		offset := (filters.Page - 1) * filters.Limit
		argCount += 2
		query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argCount-1, argCount)
		args = append(args, filters.Limit, offset)
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list orders: %w", err)
	}
	defer rows.Close()

	orders := make([]domain.Order, 0)
	for rows.Next() {
		order := domain.Order{}
		var customerEmail, deliveryAddress, paymentMethod sql.NullString
		var tableID sql.NullInt64
		var tableNumber sql.NullInt32
		var tableName sql.NullString
		var tableLocation sql.NullString
		var tableCapacity sql.NullInt32

		err := rows.Scan(
			&order.ID, &order.TenantID, &order.RestaurantID, &order.OrderNumber,
			&order.CustomerName, &customerEmail,
			&order.CustomerPhone, &deliveryAddress, &order.Subtotal, &order.TaxAmount, &order.DiscountAmount,
			&order.DeliveryFee, &order.TotalAmount, &paymentMethod, &order.PaymentStatus, &order.Status,
			&order.OrderSource, &tableID, &tableNumber, &tableName, &tableLocation, &tableCapacity,
			&order.CreatedAt, &order.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan order: %w", err)
		}

		if customerEmail.Valid {
			order.CustomerEmail = customerEmail.String
		}
		if deliveryAddress.Valid {
			order.DeliveryAddress = deliveryAddress.String
		}
		if paymentMethod.Valid {
			order.PaymentMethod = paymentMethod.String
		}
		if tableID.Valid {
			order.TableID = &tableID.Int64
		}
		if tableNumber.Valid {
			tableNumberValue := int(tableNumber.Int32)
			order.TableNumber = &tableNumberValue
		}
		if tableName.Valid {
			order.TableName = tableName.String
		}
		if tableLocation.Valid {
			order.TableLocation = tableLocation.String
		}
		if tableCapacity.Valid {
			tableCapacityValue := int(tableCapacity.Int32)
			order.TableCapacity = &tableCapacityValue
		}

		orders = append(orders, order)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating orders: %w", err)
	}

	page := int64(1)
	limit := int64(10)
	if filters != nil {
		if filters.Page > 0 {
			page = filters.Page
		}
		if filters.Limit > 0 {
			limit = filters.Limit
		}
	}

	pages := (total + limit - 1) / limit

	return &domain.OrderListResponse{
		Orders: orders,
		Total:  total,
		Page:   page,
		Limit:  limit,
		Pages:  pages,
	}, nil
}

// GetOrderItemCounts returns item counts for a page of orders in one scoped query.
func (r *OrderRepository) GetOrderItemCounts(tenantID, restaurantID int64, orderIDs []int64) (map[int64]int, error) {
	counts := make(map[int64]int, len(orderIDs))
	if len(orderIDs) == 0 {
		return counts, nil
	}

	query := `
		SELECT oi.order_id, COUNT(*)
		FROM order_items oi
		JOIN orders o
		  ON o.id = oi.order_id
		 AND o.tenant_id = oi.tenant_id
		 AND o.restaurant_id = oi.restaurant_id
		WHERE o.tenant_id = $1
		  AND o.restaurant_id = $2
		  AND oi.tenant_id = $1
		  AND oi.restaurant_id = $2
		  AND oi.order_id = ANY($3)
		GROUP BY oi.order_id
	`

	rows, err := r.db.Query(query, tenantID, restaurantID, pq.Array(orderIDs))
	if err != nil {
		return nil, fmt.Errorf("failed to get order item counts: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var orderID int64
		var count int
		if err := rows.Scan(&orderID, &count); err != nil {
			return nil, fmt.Errorf("failed to scan order item count: %w", err)
		}
		counts[orderID] = count
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate order item counts: %w", err)
	}
	return counts, nil
}

// UpdateOrderStatus updates the status of an order
func (r *OrderRepository) UpdateOrderStatus(tenantID, restaurantID, orderID int64, newStatus string, changedBy *int64, reason string) error {
	if !domain.ValidStatus(newStatus) {
		return errors.New("invalid order status")
	}
	return r.transitionOrderStatus(tenantID, restaurantID, orderID, newStatus, changedBy, reason)
}

// UpdatePaymentStatus updates the payment status of an order
func (r *OrderRepository) UpdatePaymentStatus(tenantID, restaurantID, orderID int64, paymentStatus string) error {
	return r.UpdatePaymentStatusAs(tenantID, restaurantID, orderID, paymentStatus, nil)
}

func (r *OrderRepository) UpdatePaymentStatusAs(tenantID, restaurantID, orderID int64, paymentStatus string, postedBy *int) error {
	if !domain.ValidPaymentStatus(paymentStatus) {
		return errors.New("invalid payment status")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin payment status update: %w", err)
	}
	defer tx.Rollback()
	query := `
		UPDATE orders
		SET payment_status = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2 AND tenant_id = $3 AND restaurant_id = $4
	`

	result, err := tx.ExecContext(ctx, query, paymentStatus, orderID, tenantID, restaurantID)
	if err != nil {
		return fmt.Errorf("failed to update payment status: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return errors.New("order not found")
	}
	if strings.EqualFold(strings.TrimSpace(paymentStatus), "paid") {
		events := NewAccountingEventRepository(r.db)
		if err := events.EnqueuePaidOrderTx(ctx, tx, tenantID, restaurantID, orderID, postedBy); err != nil {
			return err
		}
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit payment status and accounting intent: %w", err)
	}

	return nil
}

// CancelOrder cancels an order
func (r *OrderRepository) CancelOrder(tenantID, restaurantID, orderID int64, reason string) error {
	return r.transitionOrderStatus(tenantID, restaurantID, orderID, "cancelled", nil, reason)
}

// transitionOrderStatus serializes status changes and commits the status and its
// audit row together so concurrent requests cannot validate against stale state.
func (r *OrderRepository) transitionOrderStatus(tenantID, restaurantID, orderID int64, newStatus string, changedBy *int64, reason string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("begin order status transition: %w", err)
	}
	defer tx.Rollback()

	var currentStatus string
	err = tx.QueryRow(
		`SELECT status FROM orders
		 WHERE id = $1 AND tenant_id = $2 AND restaurant_id = $3
		 FOR UPDATE`,
		orderID, tenantID, restaurantID,
	).Scan(&currentStatus)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("order not found")
		}
		return fmt.Errorf("get order status: %w", err)
	}
	if !domain.CanTransitionStatus(currentStatus, newStatus) {
		return fmt.Errorf("invalid status transition from %s to %s", currentStatus, newStatus)
	}

	if _, err := tx.Exec(
		`UPDATE orders SET status = $1, updated_at = CURRENT_TIMESTAMP
		 WHERE id = $2 AND tenant_id = $3 AND restaurant_id = $4`,
		newStatus, orderID, tenantID, restaurantID,
	); err != nil {
		return fmt.Errorf("update order status: %w", err)
	}

	var changedByValue, reasonValue interface{}
	if changedBy != nil {
		changedByValue = *changedBy
	}
	if reason != "" {
		reasonValue = reason
	}
	if _, err := tx.Exec(
		`INSERT INTO order_status_history (order_id, old_status, new_status, changed_by, change_reason)
		 VALUES ($1, $2, $3, $4, $5)`,
		orderID, currentStatus, newStatus, changedByValue, reasonValue,
	); err != nil {
		return fmt.Errorf("record order status change: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit order status transition: %w", err)
	}
	return nil
}

// CreateOrderItem inserts an order item
func (r *OrderRepository) CreateOrderItem(tenantID, restaurantID int64, item *domain.OrderItem) (*domain.OrderItem, error) {
	return r.insertOrderItem(r.db, tenantID, restaurantID, item)
}

func (r *OrderRepository) insertOrderItem(writer orderRowWriter, tenantID, restaurantID int64, item *domain.OrderItem) (*domain.OrderItem, error) {
	query := `
		INSERT INTO order_items (
			order_id, tenant_id, restaurant_id,
			product_id, variant_id, product_name, variant_name,
			quantity, unit_price, discount_amount, total_price,
			special_instructions, addons
		) SELECT
			$1::bigint, $2::bigint, $3::bigint, $4::bigint, $5::bigint, $6, $7,
			$8, $9, $10, $11, $12, $13
		WHERE EXISTS (
			SELECT 1 FROM orders WHERE id = $1::bigint AND tenant_id = $2::bigint AND restaurant_id = $3::bigint
		)
		AND EXISTS (
			SELECT 1 FROM products WHERE id = $4::bigint AND tenant_id = $2::bigint AND restaurant_id = $3::bigint
		)
		RETURNING id, created_at, updated_at
	`

	// Marshal addons to JSON
	var addonsJSON interface{} = nil
	if len(item.AddOns) > 0 {
		bytes, err := json.Marshal(item.AddOns)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal addons: %w", err)
		}
		addonsJSON = string(bytes)
	}

	// Handle nullable fields
	var variantID sql.NullInt64
	if item.VariantID != nil {
		variantID = sql.NullInt64{Int64: *item.VariantID, Valid: true}
	}

	var variantNameStr sql.NullString
	if item.VariantName != "" {
		variantNameStr = sql.NullString{String: item.VariantName, Valid: true}
	}

	var specialInstructions sql.NullString
	if item.SpecialInstructions != "" {
		specialInstructions = sql.NullString{String: item.SpecialInstructions, Valid: true}
	}

	err := writer.QueryRow(query,
		item.OrderID,
		tenantID,
		restaurantID,
		item.ProductID,
		variantID,
		item.ProductName,
		variantNameStr,
		item.Quantity,
		item.UnitPrice,
		item.DiscountAmount,
		item.TotalPrice,
		specialInstructions,
		addonsJSON,
	).Scan(&item.ID, &item.CreatedAt, &item.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create order item: %w", err)
	}

	return item, nil
}

// GetOrderItems retrieves all items for an order
func (r *OrderRepository) GetOrderItems(tenantID, restaurantID, orderID int64) ([]domain.OrderItem, error) {
	query := `
		SELECT
			id, order_id, product_id, variant_id, product_name, variant_name,
			quantity - cancelled_quantity AS quantity, unit_price, discount_amount, total_price,
			special_instructions, addons, created_at, updated_at
		FROM order_items
		WHERE tenant_id = $1 AND restaurant_id = $2 AND order_id = $3
		ORDER BY created_at ASC
	`

	rows, err := r.db.Query(query, tenantID, restaurantID, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get order items: %w", err)
	}
	defer rows.Close()

	items := make([]domain.OrderItem, 0)
	for rows.Next() {
		item := domain.OrderItem{}
		var variantID sql.NullInt64
		var variantName, specialInstructions sql.NullString
		var addonsJSON []byte

		err := rows.Scan(
			&item.ID, &item.OrderID, &item.ProductID, &variantID,
			&item.ProductName, &variantName,
			&item.Quantity, &item.UnitPrice, &item.DiscountAmount, &item.TotalPrice,
			&specialInstructions, &addonsJSON, &item.CreatedAt, &item.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan order item: %w", err)
		}

		if variantID.Valid {
			item.VariantID = &variantID.Int64
		}

		if variantName.Valid {
			item.VariantName = variantName.String
		}

		if specialInstructions.Valid {
			item.SpecialInstructions = specialInstructions.String
		}

		// Store addons as raw JSON
		if len(addonsJSON) > 0 {
			item.AddOns = json.RawMessage(addonsJSON)
		}

		items = append(items, item)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating order items: %w", err)
	}

	return items, nil
}

// GetOrderStatusHistory retrieves the status change history for an order
func (r *OrderRepository) GetOrderStatusHistory(tenantID, restaurantID, orderID int64) ([]domain.OrderStatusHistory, error) {
	query := `
		SELECT
			h.id, h.order_id, h.old_status, h.new_status, h.changed_by,
			h.change_reason, h.created_at
		FROM order_status_history h
		JOIN orders o ON o.id = h.order_id
		WHERE o.tenant_id = $1 AND o.restaurant_id = $2 AND h.order_id = $3
		ORDER BY h.created_at DESC
	`

	rows, err := r.db.Query(query, tenantID, restaurantID, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get order status history: %w", err)
	}
	defer rows.Close()

	history := make([]domain.OrderStatusHistory, 0)
	for rows.Next() {
		h := domain.OrderStatusHistory{}
		var oldStatus, changeReason sql.NullString
		var changedBy sql.NullInt64

		err := rows.Scan(
			&h.ID, &h.OrderID, &oldStatus, &h.NewStatus, &changedBy,
			&changeReason, &h.CreatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan status history: %w", err)
		}

		if oldStatus.Valid {
			h.OldStatus = &oldStatus.String
		}

		if changedBy.Valid {
			h.ChangedBy = &changedBy.Int64
		}

		if changeReason.Valid {
			h.ChangeReason = changeReason.String
		}

		history = append(history, h)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating status history: %w", err)
	}

	return history, nil
}

// GetOrderStats retrieves statistics for orders
func (r *OrderRepository) GetOrderStats(tenantID, restaurantID int64, startDate, endDate time.Time) (*domain.OrderStats, error) {
	// Total orders and revenue
	query := `
		SELECT
			COUNT(*) as total_orders,
			COALESCE(SUM(total_amount), 0) as total_revenue,
			COALESCE(AVG(total_amount), 0) as average_order_value,
			COUNT(CASE WHEN status = 'delivered' THEN 1 END) as delivered_count,
			COUNT(CASE WHEN status = 'pending' THEN 1 END) as pending_count,
			COUNT(CASE WHEN created_at >= CURRENT_DATE THEN 1 END) as delivered_today
		FROM orders
		WHERE tenant_id = $1 AND restaurant_id = $2
		AND created_at >= $3 AND created_at <= $4
	`

	stats := &domain.OrderStats{
		OrdersByStatus:        make(map[string]int64),
		OrdersByPaymentMethod: make(map[string]int64),
	}

	var deliveredToday int64
	err := r.db.QueryRow(query, tenantID, restaurantID, startDate, endDate).Scan(
		&stats.TotalOrders,
		&stats.TotalRevenue,
		&stats.AverageOrderValue,
		&stats.DeliveredCount,
		&stats.PendingCount,
		&deliveredToday,
	)

	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to get order stats: %w", err)
	}

	stats.DeliveredToday = deliveredToday

	// Orders by status
	statusQuery := `
		SELECT status, COUNT(*) as count
		FROM orders
		WHERE tenant_id = $1 AND restaurant_id = $2
		AND created_at >= $3 AND created_at <= $4
		GROUP BY status
	`

	rows, err := r.db.Query(statusQuery, tenantID, restaurantID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get status breakdown: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var status string
		var count int64
		err := rows.Scan(&status, &count)
		if err != nil {
			return nil, fmt.Errorf("failed to scan status breakdown: %w", err)
		}
		stats.OrdersByStatus[status] = count
	}

	// Orders by payment method
	paymentQuery := `
		SELECT payment_method, COUNT(*) as count
		FROM orders
		WHERE tenant_id = $1 AND restaurant_id = $2
		AND created_at >= $3 AND created_at <= $4
		AND payment_method IS NOT NULL
		GROUP BY payment_method
	`

	rows, err = r.db.Query(paymentQuery, tenantID, restaurantID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get payment method breakdown: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var method string
		var count int64
		err := rows.Scan(&method, &count)
		if err != nil {
			return nil, fmt.Errorf("failed to scan payment method breakdown: %w", err)
		}
		stats.OrdersByPaymentMethod[method] = count
	}

	return stats, nil
}

// DeleteOrder deletes an order (soft delete by cancelling it)
func (r *OrderRepository) DeleteOrder(tenantID, restaurantID, orderID int64) error {
	// We prefer to cancel orders rather than delete them for audit trail
	return r.CancelOrder(tenantID, restaurantID, orderID, "Deleted by admin")
}

// SetPublicTokenHash stores the hash of the public tracking token for an order.
func (r *OrderRepository) SetPublicTokenHash(tenantID, restaurantID, orderID int64, hash string) error {
	res, err := r.db.Exec(
		`UPDATE orders SET public_token_hash = $1 WHERE id = $2 AND tenant_id = $3 AND restaurant_id = $4`,
		hash, orderID, tenantID, restaurantID,
	)
	if err != nil {
		return fmt.Errorf("failed to store public token: %w", err)
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return errors.New("order not found")
	}
	return nil
}

// GetPublicAccessInfo returns the stored public token hash (empty for legacy orders) and the
// customer phone of an order, scoped to tenant and restaurant.
func (r *OrderRepository) GetPublicAccessInfo(tenantID, restaurantID, orderID int64) (hash string, phone string, err error) {
	var h sql.NullString
	err = r.db.QueryRow(
		`SELECT public_token_hash, customer_phone FROM orders WHERE id = $1 AND tenant_id = $2 AND restaurant_id = $3`,
		orderID, tenantID, restaurantID,
	).Scan(&h, &phone)
	if err == sql.ErrNoRows {
		return "", "", errors.New("order not found")
	}
	if err != nil {
		return "", "", fmt.Errorf("failed to read order access info: %w", err)
	}
	return h.String, phone, nil
}
