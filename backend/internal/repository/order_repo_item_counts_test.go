package repository

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestGetOrderItemCountsScopesAndBatchesPage(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("create sql mock: %v", err)
	}
	defer db.Close()

	query := `(?s)SELECT oi\.order_id, COUNT\(\*\).*WHERE o\.tenant_id = \$1.*AND o\.restaurant_id = \$2.*AND oi\.tenant_id = \$1.*AND oi\.restaurant_id = \$2.*AND oi\.order_id = ANY\(\$3\).*GROUP BY oi\.order_id`
	if _, err := regexp.Compile(query); err != nil {
		t.Fatalf("invalid query matcher: %v", err)
	}
	mock.ExpectQuery(query).
		WithArgs(int64(7), int64(9), "{101,102,103}").
		WillReturnRows(sqlmock.NewRows([]string{"order_id", "count"}).
			AddRow(int64(101), 2).
			AddRow(int64(103), 4))

	repo := NewOrderRepository(db)
	counts, err := repo.GetOrderItemCounts(7, 9, []int64{101, 102, 103})
	if err != nil {
		t.Fatalf("GetOrderItemCounts returned error: %v", err)
	}
	if counts[101] != 2 || counts[102] != 0 || counts[103] != 4 {
		t.Fatalf("unexpected item counts: %#v", counts)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet SQL expectations: %v", err)
	}
}

func TestGetOrderItemCountsSkipsQueryForEmptyPage(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("create sql mock: %v", err)
	}
	defer db.Close()

	counts, err := NewOrderRepository(db).GetOrderItemCounts(7, 9, nil)
	if err != nil {
		t.Fatalf("GetOrderItemCounts returned error: %v", err)
	}
	if len(counts) != 0 {
		t.Fatalf("expected empty counts, got %#v", counts)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unexpected SQL query: %v", err)
	}
}
