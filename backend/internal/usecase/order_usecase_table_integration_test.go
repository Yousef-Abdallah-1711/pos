package usecase

import (
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"

	"pos-saas/internal/domain"
	"pos-saas/internal/repository"
)

func tableColumnsForOrderTests() []string {
	return []string{
		"id", "tenant_id", "restaurant_id", "table_number", "name", "capacity", "location", "status",
		"description", "qr_code_config", "position_x", "position_y", "shape", "layout_id", "created_at", "updated_at",
	}
}

func productColumnsForOrderTests() []string {
	return []string{
		"id", "tenant_id", "restaurant_id", "category_id", "sku", "barcode",
		"name_en", "name_ar", "description_en", "description_ar",
		"price", "cost", "discount_price", "discount_percentage",
		"calories", "protein_g", "carbs_g", "fat_g", "fiber_g", "allergens",
		"is_vegetarian", "is_vegan", "is_spicy", "is_gluten_free",
		"is_available", "available_from", "available_until", "available_days",
		"track_inventory", "quantity_in_stock", "low_stock_threshold", "reorder_quantity",
		"display_order", "featured", "main_image_url", "status",
		"created_by", "updated_by", "created_at", "updated_at",
	}
}

func TestOrderUseCase_CreateOrderAutoOccupiesTable(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	orderRepo := repository.NewOrderRepository(db)
	productRepo := repository.NewProductRepository(db)
	tableRepo := repository.NewTableRepository(db)

	uc := NewOrderUseCase(orderRepo, productRepo)
	uc.SetTableRepository(tableRepo)

	now := time.Now()
	tableID := int64(8)

	mock.ExpectQuery("FROM restaurant_tables").
		WithArgs(tableID, int64(1), int64(2)).
		WillReturnRows(
			sqlmock.NewRows(tableColumnsForOrderTests()).
				AddRow(
					tableID, int64(1), int64(2), 8, "Garden 8", 4, "indoor", "available",
					"", []byte(`{}`), 0.0, 0.0, "square", sql.NullInt64{}, now, now,
				),
		)

	mock.ExpectQuery("FROM products").
		WithArgs(101, 1, 2).
		WillReturnRows(
			sqlmock.NewRows(productColumnsForOrderTests()).
				AddRow(
					101, 1, 2, 1, nil, nil,
					"Burger", nil, nil, nil,
					12.5, nil, nil, nil,
					nil, nil, nil, nil, nil, []byte(`[]`),
					false, false, false, false,
					true, nil, nil, []byte(`[]`),
					false, 100, 5, 0,
					1, false, nil, "active",
					1, nil, now, now,
				),
		)

	mock.ExpectQuery("INSERT INTO orders").
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).AddRow(int64(501), now, now))

	mock.ExpectQuery("INSERT INTO order_items").
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).AddRow(int64(1), now, now))

	mock.ExpectQuery("FROM orders o").
		WithArgs(int64(501), int64(1), int64(2)).
		WillReturnRows(
			sqlmock.NewRows([]string{
				"id", "tenant_id", "restaurant_id", "order_number", "customer_name", "customer_email",
				"customer_phone", "delivery_address", "delivery_city", "delivery_area",
				"delivery_zip_code", "delivery_latitude", "delivery_longitude",
				"delivery_instructions", "subtotal", "tax_amount", "discount_amount",
				"delivery_fee", "total_amount", "payment_method", "payment_status",
				"status", "estimated_delivery_time", "actual_delivery_time", "notes", "order_source",
				"table_id", "table_number", "table_name", "table_location", "capacity",
				"created_at", "updated_at",
			}).AddRow(
				int64(501), int64(1), int64(2), "ORD-2026-000001", "Table Guest", nil,
				"01000000000", nil, nil, nil,
				nil, nil, nil,
				nil, 25.0, 2.5, 0.0,
				0.0, 27.5, "cash", "pending",
				"pending", nil, nil, "", "dine_in",
				tableID, 8, "Garden 8", "indoor", 4,
				now, now,
			),
		)

	mock.ExpectQuery("FROM order_items").
		WithArgs(int64(1), int64(2), int64(501)).
		WillReturnRows(
			sqlmock.NewRows([]string{
				"id", "order_id", "product_id", "variant_id", "product_name", "variant_name",
				"quantity", "unit_price", "discount_amount", "total_price",
				"special_instructions", "addons", "created_at", "updated_at",
			}).AddRow(
				int64(1), int64(501), int64(101), nil, "Burger", nil,
				2, 12.5, 0.0, 25.0,
				"", []byte(`[]`), now, now,
			),
		)

	mock.ExpectExec("UPDATE orders SET public_token_hash").
		WithArgs(sqlmock.AnyArg(), int64(501), int64(1), int64(2)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectExec("UPDATE restaurant_tables").
		WithArgs("occupied", tableID, int64(1), int64(2)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	order, err := uc.CreateOrder(1, 2, &domain.CreateOrderRequest{
		CustomerName:  "Table Guest",
		CustomerPhone: "01000000000",
		PaymentMethod: "cash",
		OrderSource:   "dine_in",
		TableID:       &tableID,
		Items: []domain.CreateOrderItemRequest{
			{
				ProductID: 101,
				Quantity:  2,
			},
		},
	})
	if err != nil {
		t.Fatalf("expected no error creating order, got %v", err)
	}
	if order.TableID == nil || *order.TableID != tableID {
		t.Fatalf("expected order to include table_id=%d, got %+v", tableID, order.TableID)
	}
	if order.OrderSource != "dine_in" {
		t.Fatalf("expected order source dine_in, got %s", order.OrderSource)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sqlmock expectations: %v", err)
	}
}
