package usecase

import (
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	"math"
	"pos-saas/internal/domain"
	"pos-saas/internal/repository"
	"strings"
	"testing"
	"time"
)

func TestOrderPricingPreservesServerTaxAndFeeWithCentRounding(t *testing.T) {
	product := domain.Product{Status: "active", Price: 9.99}
	item, err := pricedOrderItem(product, domain.CreateOrderItemRequest{ProductID: 7, Quantity: 3})
	if err != nil {
		t.Fatal(err)
	}
	order := &domain.Order{Items: []domain.OrderItem{item}}
	if err := applyOrderPricing(order); err != nil {
		t.Fatal(err)
	}
	if order.Subtotal != 29.97 || order.TaxAmount != 3 || order.DeliveryFee != 0 || order.TotalAmount != 32.97 {
		t.Fatalf("incorrect server pricing: %+v", order)
	}
}

func TestOrderPricingRejectsInvalidPricesQuantitiesAndNumericOverflow(t *testing.T) {
	for _, price := range []float64{-1, math.NaN(), math.Inf(1), 100_000_000} {
		if _, err := pricedOrderItem(domain.Product{Status: "active", Price: price}, domain.CreateOrderItemRequest{Quantity: 1}); err == nil {
			t.Fatal("invalid price accepted")
		}
	}
	if _, err := pricedOrderItem(domain.Product{Status: "active", Price: 1}, domain.CreateOrderItemRequest{Quantity: 0}); err == nil {
		t.Fatal("zero quantity accepted")
	}
	if _, err := pricedOrderItem(domain.Product{Status: "active", Price: 99_999_999}, domain.CreateOrderItemRequest{Quantity: 10000}); err == nil {
		t.Fatal("overflow accepted")
	}
}

func TestOrderQuoteUsesScopedDatabasePricesWithoutCreatingOrReservingOrders(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	mock.ExpectQuery(`FROM products WHERE tenant_id=\$1 AND restaurant_id=\$2 AND id=ANY`).WithArgs(int64(1), int64(2), sqlmock.AnyArg()).WillReturnRows(sqlmock.NewRows([]string{"id", "name", "price", "discount", "status", "available", "tracked", "stock"}).AddRow(7, "Synthetic", 25, nil, "active", true, true, 20))
	uc := NewOrderUseCase(nil, repository.NewProductRepository(db))
	quote, err := uc.QuoteOrder(context.Background(), 1, 2, &domain.OrderQuoteRequest{Items: []domain.CreateOrderItemRequest{{ProductID: 7, Quantity: 2}}})
	if err != nil {
		t.Fatal(err)
	}
	if quote.Subtotal != 50 || quote.TaxAmount != 5 || quote.DeliveryFee != 0 || quote.TotalAmount != 55 {
		t.Fatalf("unexpected quote: %+v", quote)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestOrderQuoteAggregatesDuplicateProductStockAndRejectsForeignProducts(t *testing.T) {
	for _, missing := range []bool{false, true} {
		t.Run(map[bool]string{false: "duplicate_stock", true: "foreign_scope"}[missing], func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatal(err)
			}
			defer db.Close()
			rows := sqlmock.NewRows([]string{"id", "name", "price", "discount", "status", "available", "tracked", "stock"})
			if !missing {
				rows.AddRow(7, "Synthetic", 25, nil, "active", true, true, 1)
			}
			mock.ExpectQuery("FROM products WHERE tenant_id").WithArgs(int64(1), int64(2), sqlmock.AnyArg()).WillReturnRows(rows)
			uc := NewOrderUseCase(nil, repository.NewProductRepository(db))
			_, err = uc.QuoteOrder(context.Background(), 1, 2, &domain.OrderQuoteRequest{Items: []domain.CreateOrderItemRequest{{ProductID: 7, Quantity: 1}, {ProductID: 7, Quantity: 1}}})
			if err == nil {
				t.Fatal("unsafe quote accepted")
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestOrderQuoteRejectsUnboundedOrInvalidItemsBeforeQuery(t *testing.T) {
	uc := NewOrderUseCase(nil, nil)
	for _, items := range [][]domain.CreateOrderItemRequest{nil, {{ProductID: 7, Quantity: 0}}, {{ProductID: -1, Quantity: 1}}, make([]domain.CreateOrderItemRequest, 101)} {
		if _, err := uc.QuoteOrder(context.Background(), 1, 2, &domain.OrderQuoteRequest{Items: items}); err == nil {
			t.Fatal("invalid quote accepted")
		}
	}
}

func TestOrderCreationRejectsChangedQuoteBeforeAnyDatabaseWrite(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	now := time.Now()
	mock.ExpectQuery("FROM products").WithArgs(101, 1, 2).WillReturnRows(sqlmock.NewRows(productColumnsForOrderTests()).AddRow(
		101, 1, 2, 1, nil, nil, "Synthetic", nil, nil, nil, 12.5, nil, nil, nil, nil, nil, nil, nil, nil, []byte(`[]`), false, false, false, false, true, nil, nil, []byte(`[]`), false, 100, 5, 0, 1, false, nil, "active", 1, nil, now, now))
	expected := 1.0
	uc := NewOrderUseCase(nil, repository.NewProductRepository(db))
	_, err = uc.CreateOrder(1, 2, &domain.CreateOrderRequest{CustomerName: "Synthetic", CustomerPhone: "+12025550123", PaymentMethod: "cash", ExpectedTotal: &expected, Items: []domain.CreateOrderItemRequest{{ProductID: 101, Quantity: 1}}})
	if err == nil || !strings.Contains(err.Error(), "total changed") {
		t.Fatalf("price change was not rejected: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}
